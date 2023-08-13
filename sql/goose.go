package sql

import (
	"database/sql"
	"embed"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
	config "movies/utils/config"

	goose "github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

var folderMigrations = "migrations"

func gooseInit() (*sql.DB, error) {
	pgconfig := config.PostgreSQL()
	config := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		pgconfig.User, pgconfig.Password, pgconfig.Host, pgconfig.Port, pgconfig.Database, "disable",
	)
	goose.SetBaseFS(embedMigrations)
	return goose.OpenDBWithDriver("pgx", config)
}

// Migrate the DB to the most recent version available.
func GooseUp() error {
	db, err := gooseInit()
	if err != nil {
		return err
	}
	return goose.Up(db, folderMigrations)
}

// Migrate the DB up by 1.
func GooseUpByOne() error {
	db, err := gooseInit()
	if err != nil {
		return err
	}
	return goose.UpByOne(db, folderMigrations)
}

// Migrate the DB to a specific VERSION.
func GooseUpTo(version int64) error {
	db, err := gooseInit()
	if err != nil {
		return err
	}
	return goose.UpTo(db, folderMigrations, version)
}

// Creates new migration file with the current timestamp.
func GooseCreate(name string) error {
	db, err := gooseInit()
	if err != nil {
		return err
	}
	return goose.Create(db, "sql/"+folderMigrations, name, "sql")
}

// Print the current version of the database.
func GooseVersion() error {
	db, err := gooseInit()
	if err != nil {
		return err
	}
	return goose.Version(db, folderMigrations)
}

// Apply sequential ordering to migrations.
func GooseFix() error {
	return goose.Fix("sql/" + folderMigrations)
}

// Dump the migration status for the current DB.
func GooseStatus() error {
	db, err := gooseInit()
	if err != nil {
		return err
	}
	return goose.Status(db, folderMigrations)
}

// Roll back all migrations.
func GooseReset() error {
	db, err := gooseInit()
	if err != nil {
		return err
	}
	return goose.Reset(db, folderMigrations)
}

// Re-run the latest migration.
func GooseRedo() error {
	db, err := gooseInit()
	if err != nil {
		return err
	}
	return goose.Redo(db, folderMigrations)
}

// Roll back the version by 1.
func GooseDown() error {
	db, err := gooseInit()
	if err != nil {
		return err
	}
	return goose.Down(db, folderMigrations)
}

// Roll back to a specific VERSION.
func GooseDownTo(version int64) error {
	db, err := gooseInit()
	if err != nil {
		return err
	}
	return goose.DownTo(db, folderMigrations, version)
}
