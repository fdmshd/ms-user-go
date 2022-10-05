package models

import (
	"database/sql"
	"errors"
)

var ErrNoRecord = errors.New("models: no matching record found")

type User struct {
	Id       int    `json: "id"`
	Username string `json: "username"`
	Email    string `json: "email"`
	Token    string `json:token, omitempty`
	Password string `json:"password,omitempty"`
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(u User) (int, error) {

	stmt := `INSERT INTO users (username, email, password)
	VALUES(?,?,?)`

	res, err := m.DB.Exec(stmt, u.Username, u.Email, u.Password)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (m *UserModel) Get(id int) (*User, error) {

	stmt := `SELECT id, username, email FROM users
    WHERE id = ?`

	row := m.DB.QueryRow(stmt, id)
	u := &User{}

	err := row.Scan(&u.Id, &u.Username, &u.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return u, nil
}

func (m *UserModel) GetByName(name string) (*User, error) {
	stmt := `SELECT id, username,email, password FROM users
    WHERE username = ?`

	row := m.DB.QueryRow(stmt, name)
	u := &User{}

	err := row.Scan(&u.Id, &u.Username, &u.Email, &u.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return u, nil
}

func (m *UserModel) Update(id int) (*User, error) {
	return nil, nil
}

func (m *UserModel) Delete(id int) error {
	return nil
}
