package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Extaleus/password-manager-app/internal/storage"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New" // for pure debug

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// TODO: need to read about DB migrations
	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS url(
		id INTEGER PRIMARY KEY,
		alias TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL);
	CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveUrl(urlToSave string, alias string) (int64, error) { // int64 for (LastInsertId()) get index of new url in db // not work in some db
	const op = "storage.sqlite.SaveUrl" // for pure debug

	stmt, err := s.db.Prepare("INSERT INTO url(url, alias) VALUES(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(urlToSave, alias)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetUrl(alias string) (string, error) {
	const op = "storage.sqlite.GetUrl" // for pure debug

	stmt, err := s.db.Prepare("SELECT url from url WHERE alias = ?")
	if err != nil {
		return "", fmt.Errorf("%s: prepare: %w", op, err)
	}

	var resUrl string
	err = stmt.QueryRow(alias).Scan(&resUrl)
	if errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf(storage.ErrURLNotFound.Error())
	}
	if err != nil {
		return "", fmt.Errorf("%s: execute: %w", op, err)
	}

	return resUrl, nil
}

func (s *Storage) DeleteUrl(alias string) error {
	const op = "storage.sqlite.DeleteUrl" // for pure debug

	stmt, err := s.db.Prepare("DELETE FROM url WHERE alias = ?")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(alias)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_ = res

	return nil
}
