package sql

import (
	"context"
	"embed"
	"fmt"
	"log"
	"path"
	"strings"
	"time"

	pg "movies/utils/pg"

	pgerrcode "github.com/jackc/pgerrcode"
)

//go:embed schema/*
var sqlsFs embed.FS

var sqlDir = "schema"

// Reset the database.
func Reset(ctx context.Context) error {
	startAt := time.Now()
	log.SetPrefix("[sql] ")

	// Drop
	log.Print("Drop")
	if err := Drop(ctx); err != nil {
		return err
	}

	// Init
	log.Print("Init")
	if err := Init(ctx); err != nil {
		return err
	}

	log.Printf("End %v", time.Since(startAt))
	log.SetPrefix("")

	return nil
}

// Init the database.
func Init(ctx context.Context) error {
	// Get SQL files
	files, err := sqlsFs.ReadDir(sqlDir)
	if err != nil {
		panic(err)
	}

	// Run SQL files with a suffix
	runSQL := func(suffix string) error {
		for _, file := range files {
			if strings.Contains(file.Name(), suffix) {
				if err = RunFile(ctx, file.Name()); err != nil {
					return err
				}
			}
		}
		return nil
	}

	if err := runSQL("extensions"); err != nil {
		return err
	} else if err := runSQL("misc"); err != nil {
		return err
	} else if err := runSQL("table"); err != nil {
		return err
	} else if err := runSQL("alter"); err != nil {
		return err
	}

	// Goose Up
	if err := GooseUp(); err != nil {
		return err
	}

	return nil
}

// Drop the database.
func Drop(ctx context.Context) error {
	var username string
	var schema string = "private"

	// Get current username
	if err := pg.Client(pg.EmptyTx()).QueryRow(ctx, `SELECT current_user`).Scan(&username); err != nil {
		return err
	}

	// Drop schema if not exists
	_, err := pg.Client(pg.EmptyTx()).Exec(ctx, fmt.Sprintf(`DROP SCHEMA "%s" CASCADE`, schema))
	if err != nil && !pg.IsErrCode(err, pgerrcode.InvalidSchemaName) {
		return err
	}

	// Create schema
	if _, err = pg.Client(pg.EmptyTx()).Exec(ctx, fmt.Sprintf(`CREATE SCHEMA "%s" AUTHORIZATION "%s"`, schema, username)); err != nil {
		return err
	}

	// Set default schema
	if _, err = pg.Client(pg.EmptyTx()).Exec(ctx, fmt.Sprintf(`SET search_path TO "%s","public"`, schema)); err != nil {
		return err
	}

	// Add user to schema
	if _, err = pg.Client(pg.EmptyTx()).Exec(ctx, fmt.Sprintf(`ALTER USER "%s" SET search_path = "%s","public"`, username, schema)); err != nil {
		return err
	}

	return nil
}

// Run a SQL file.
func RunFile(ctx context.Context, name string) error {
	path := path.Join(sqlDir, name)
	if sql, err := sqlsFs.ReadFile(path); err != nil {
		return err
	} else if _, err = pg.Client(pg.EmptyTx()).Exec(ctx, string(sql)); err != nil {
		return err
	}
	return nil
}
