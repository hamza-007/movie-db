// package pg

// import (
// 	"context"
// 	"fmt"
// 	"os"

// 	dotenv "github.com/joho/godotenv"
// 	postgres "gorm.io/driver/postgres"
// 	gorm "gorm.io/gorm"

// // PG Driver for goqu.
// _ "github.com/doug-martin/goqu/v9/dialect/postgres"
// pgconn "github.com/jackc/pgconn"
// pgtype "github.com/jackc/pgtype"
// shopspring "github.com/jackc/pgtype/ext/shopspring-numeric"
// pgx "github.com/jackc/pgx/v4"
// pgxpool "github.com/jackc/pgx/v4/pgxpool"
// lo "github.com/samber/lo"
package pg

import (
	"context"
	"fmt"
	"hash/fnv"
	"log"
	"sync"

	config "movies/utils/config"

	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	pgconn "github.com/jackc/pgconn"
	pgtype "github.com/jackc/pgtype"
	shopspring "github.com/jackc/pgtype/ext/shopspring-numeric"
	pgx "github.com/jackc/pgx/v4"
	pgxpool "github.com/jackc/pgx/v4/pgxpool"
	lo "github.com/samber/lo"
	viper "github.com/spf13/viper"
)

var (
	pgxp     *pgxpool.Pool
	pgxpOnce sync.Once
)

type qLogger struct{}

func (l *qLogger) Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]any) {
	if level == pgx.LogLevelInfo && msg == "Query" {
		log.Printf("SQL: %s || ARGS: %v\n", data["sql"], data["args"])
	}
}

// Connect Init conn.
func connect() {
	c := config.PostgreSQL()

	config := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.Database, "disable",
	)

	poolConfig, err := pgxpool.ParseConfig(config)
	if err != nil {
		log.Fatalf("Error pool config: %v", err)
	}

	if viper.GetBool("log-db") {
		poolConfig.ConnConfig.Logger = &qLogger{}
	}

	poolConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		conn.ConnInfo().RegisterDataType(pgtype.DataType{
			Value: &shopspring.Numeric{},
			Name:  "numeric",
			OID:   pgtype.NumericOID,
		})
		return nil
	}

	pool, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatalf("Unable to connect PGX: %v", err)
	}
	pgxp = pool
}

func pool() *pgxpool.Pool {
	pgxpOnce.Do(func() {
		connect()
	})
	return pgxp
}

// Close Close all conn.
func Close() {
	if pool() != nil {
		pool().Close()
	}
}

/*============================================================================*/
/*=====*                            Querier                             *=====*/
/*============================================================================*/

// Querier is something that pgxscan can query and get the pgx.Rows from.
// It can be: *pgxpool.Pool, *pgx.Conn or pgx.Tx.
type Querier interface {
	Query(context.Context, string, ...any) (pgx.Rows, error)
	QueryRow(context.Context, string, ...any) pgx.Row
	QueryFunc(context.Context, string, []any, []any, func(pgx.QueryFuncRow) error) (pgconn.CommandTag, error)
	Begin(context.Context) (pgx.Tx, error)
	SendBatch(context.Context, *pgx.Batch) pgx.BatchResults
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
}

// Client : Get transaction or new Querier.
func Client(tx Tx) Querier {
	return lo.Ternary[Querier](!tx.empty, tx.pgxTx, pool())
}

func Ping(ctx context.Context) error {
	return pool().Ping(ctx)
}

func Lock(ctx context.Context, tx Tx, key1, key2 string) error {
	// Ensure to be inside a transaction
	tx, err := EnsureTx(ctx, tx)
	defer tx.RollbackDefer(ctx)
	if err != nil {
		return err
	}

	if config.IsTest() {
		sql := "set deadlock_timeout='3s';"
		if _, err := Client(tx).Exec(ctx, sql); err != nil {
			return err
		}
	}

	k1 := fnv.New32a()
	k1.Write([]byte(key1))
	k2 := fnv.New32a()
	k2.Write([]byte(key2))
	sql := fmt.Sprintf(
		"SELECT pg_advisory_xact_lock(%d,%d);",
		int32(k1.Sum32()),
		int32(k2.Sum32()),
	)

	if _, err := Client(tx).Exec(ctx, sql); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
