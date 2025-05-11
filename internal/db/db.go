package database

import (
	app "calc_parallel/internal/app"
	"database/sql"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db  *sql.DB
	mux *sync.Mutex
}

type User_Expression struct {
	Expression_ID
	UserID string
}

type Expression_ID struct {
	ID uint32 `json:"id"`
	app.Expression
}

func (st *Storage) DB_Close() error {
	st.mux.Lock()
	err := st.db.Close()
	st.mux.Unlock()
	return err
}

func DB_Open(path string) (*Storage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	st := &Storage{db, &sync.Mutex{}}
	return st, st.Create_Tab()
}

func (s *Storage) Create_Tab() error {
	s.mux.Lock()
	defer s.mux.Unlock()
	const (
		usersTable = `
	CREATE TABLE IF NOT EXISTS users(
		id TEXT PRIMARY KEY, 
		login TEXT,
		password TEXT
	);`

		expressionsTable = `
	CREATE TABLE IF NOT EXISTS expressions(
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		data TEXT NOT NULL,
		status TEXT NOT NULL,
		result FLOAT,
		user_id TEXT NOT NULL
	);`
	)

	if _, err := s.db.Exec(usersTable); err != nil {
		return err
	}

	if _, err := s.db.Exec(expressionsTable); err != nil {
		return err
	}

	return nil
}

func (s *Storage) Clear() error {
	s.mux.Lock()
	defer s.mux.Unlock()
	var (
		q1 = `DELETE FROM users;`
		q2 = `DELETE FROM expressions;`
	)
	if _, err := s.db.Exec(q1); err != nil {
		return err
	}
	if _, err := s.db.Exec(q2); err != nil {
		return err
	}
	return nil
}

func (s *Storage) InsertUser(user *app.User) error {
	s.mux.Lock()
	defer s.mux.Unlock()
	var q = `INSERT INTO users (id, login, password) values ($1, $2, $3)`
	_, err := s.db.Exec(q, user.ID, user.Login, user.Password)
	return err
}

func (s *Storage) InsertExpression(exp Expression_ID, forUser *app.User) error {
	s.mux.Lock()
	defer s.mux.Unlock()
	var q = `INSERT INTO expressions (id, data, status, result, user_id) values ($1, $2, $3, $4, $5)`
	_, err := s.db.Exec(q, exp.ID, exp.Expr, exp.Status, exp.Result, forUser.ID)
	return err
}

func (s *Storage) SelectAllUsers() ([]*app.User, error) {
	s.mux.Lock()
	defer s.mux.Unlock()
	var users []*app.User
	var q = `SELECT id, login, password FROM users`
	rows, err := s.db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		u := &app.User{}
		err := rows.Scan(&u.ID, &u.Login, &u.Password)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (s *Storage) SelectExpressionsForUser(user *app.User) ([]*Expression_ID, error) {
	s.mux.Lock()
	defer s.mux.Unlock()
	var expressions []*Expression_ID
	var q = `SELECT id, data, status, result, user_id FROM expressions`

	rows, err := s.db.Query(q)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		e := &User_Expression{}
		err := rows.Scan(&e.ID, &e.Expr, &e.Status, &e.Result, &e.UserID)
		if err != nil {
			return nil, err
		}
		if e.UserID == user.ID {
			expressions = append(expressions, &e.Expression_ID)
		}
	}

	return expressions, rows.Close()
}

func (s *Storage) SelectExpressions() ([]*User_Expression, error) {
	s.mux.Lock()
	defer s.mux.Unlock()
	var expressions []*User_Expression
	var q = `SELECT id, data, status, result, user_id FROM expressions`

	rows, err := s.db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		e := &User_Expression{}
		err := rows.Scan(&e.ID, &e.Expr, &e.Status, &e.Result, &e.UserID)
		if err != nil {
			return nil, err
		}
		expressions = append(expressions, e)
	}

	return expressions, nil
}
