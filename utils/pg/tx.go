package pg

import (
	"context"
	"errors"

	logger "movies/utils/logger"

	pgconn "github.com/jackc/pgconn"
	pgx "github.com/jackc/pgx/v4"
	lo "github.com/samber/lo"
)

/*============================================================================*/
/*=====*                          Transaction                           *=====*/
/*============================================================================*/

// Wrapper of pgx Tx
type Tx struct {
	pgxTx   pgx.Tx
	scopped bool
	empty   bool

	commitHooks,
	rollbackHooks *[]func() error
}

// EmptyTx: Create an empty transaction
func EmptyTx() Tx {
	return Tx{
		empty:         true,
		commitHooks:   &[]func() error{},
		rollbackHooks: &[]func() error{},
	}
}

// IsEmpty: Is a transaction empty ?
func (t Tx) IsEmpty() bool { return t.empty }

// NewTx: Create a new transaction
func NewTx(ctx context.Context) (Tx, error) {
	t, err := pool().Begin(ctx)
	if err != nil {
		return EmptyTx(), err
	}
	return Tx{
		pgxTx:         t,
		commitHooks:   &[]func() error{},
		rollbackHooks: &[]func() error{},
	}, nil
}

// EnsureTx: Ensure to be inside a transaction
func EnsureTx(ctx context.Context, t Tx) (Tx, error) {
	if t.empty {
		return NewTx(ctx)
	}
	return Tx{
		pgxTx:         t.pgxTx,
		commitHooks:   t.commitHooks,
		rollbackHooks: t.rollbackHooks,
		scopped:       true,
	}, nil
}

// Rollback: Rollback a transaction
func (t *Tx) Rollback(ctx context.Context) error {
	if t.scopped || t.empty {
		return nil
	}
	err := t.pgxTx.Rollback(ctx)
	if errors.Is(err, pgx.ErrTxClosed) {
		return nil
	} else if err != nil {
		logger.Error(ctx, "Rollback failed: %v", err)
	} else {
		// Execute rollback functions
		lo.ForEach(*t.rollbackHooks, func(fct func() error, _ int) {
			if err := fct(); err != nil {
				logger.Error(ctx, err.Error())
			}
		})
	}
	return err
}

// RollbackDefer: Rollback a transaction and log error
func (t *Tx) RollbackDefer(ctx context.Context) {
	if err := t.Rollback(ctx); err != nil {
		logger.Error(ctx, err.Error())
	}
}

// Commit: Commit a transaction
func (t *Tx) Commit(ctx context.Context) error {
	if t.scopped || t.empty {
		return nil
	}

	err := t.pgxTx.Commit(ctx)
	if err == nil {
		// Execute commit functions
		lo.ForEach(*t.commitHooks, func(fct func() error, _ int) {
			if err := fct(); err != nil {
				logger.Error(ctx, err.Error())
			}
		})
	}
	return err
}

// OnCommit: Execute function on success commit
func (t Tx) OnCommit(fct func() error) {
	*t.commitHooks = append(*t.commitHooks, fct)
}

// OnRollback: Execute function on rollback
func (t Tx) OnRollback(fct func() error) {
	*t.rollbackHooks = append(*t.rollbackHooks, fct)
}

// Commit: Commit a transaction
func (t *Tx) Exec(ctx context.Context, sql string) (pgconn.CommandTag, error) {
	return t.pgxTx.Exec(ctx, sql)
}
