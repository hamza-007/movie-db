package pg

import (
	errors "emperror.dev/errors"
	pgxscan "github.com/georgysavva/scany/pgxscan"
	pgconn "github.com/jackc/pgconn"
	pgx "github.com/jackc/pgx/v4"
)

/*============================================================================*/
/*=====*                             Error                              *=====*/
/*============================================================================*/

// NotFound ...
func IsNotFound(err error) bool {
	return pgxscan.NotFound(err) || errors.Is(err, pgx.ErrNoRows)
}

// HavePGErr ...
func HavePGErr(err error) *pgconn.PgError {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr
	}
	return nil
}

// IsErrCode ...
func IsErrCode(err error, code string) bool {
	pgErr := HavePGErr(err)
	return pgErr != nil && pgErr.Code == code
}

// IsErrConstraint ...
func IsErrConstraint(err error, code, constraint string) bool {
	pgErr := HavePGErr(err)
	return pgErr != nil && pgErr.Code == code && pgErr.ConstraintName == constraint
}
