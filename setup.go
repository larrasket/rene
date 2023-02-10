package main

import (
	"database/sql"
	"io"
	"os"

	_ "embed"
	"golang.org/x/exp/slog"
)

// setup logger
var logger *slog.Logger

func SetLogger() error {
	file, err := os.OpenFile(".rene.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0666)
	if err != nil {
		logger = slog.New(slog.HandlerOptions{AddSource: true}.
			NewTextHandler(os.Stdout))
	} else {
		mw := io.MultiWriter(os.Stdout, file)
		logger = slog.New(slog.HandlerOptions{AddSource: true}.
			NewTextHandler(mw))
	}
	slog.SetDefault(logger)
	return err
}

// setup db

//go:embed sql/db.sql
var sqlIntialize string
var Db *sql.DB

func SetDb() error {
	var err error
	Db, err = sql.Open("sqlite3", "bin.db")
	if err != nil {
		logger.Error("Couldn't open database", err)
		return err
	}
	_, err = Db.Exec(sqlIntialize)
	return err
}
