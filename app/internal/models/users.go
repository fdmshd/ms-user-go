package models

import (
	"database/sql"
	"errors"
	"fmt"
)

var ErrNoRecord = errors.New("models: no matching record found")

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username" validate:"required,min=3,max=32"`
	Email    string `json:"email,omitempty" validate:"required,email,max=64"`
	IsAdmin  bool   `json:"is_admin"`
	Token    string `json:"token,omitempty"`
	Password string `json:"password,omitempty" validate:"required,min=6,max=32"`
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
	stmt := `SELECT id, username,email,is_admin, password FROM users
    WHERE username = ?`

	row := m.DB.QueryRow(stmt, name)
	u := &User{}

	err := row.Scan(&u.Id, &u.Username, &u.Email, &u.IsAdmin, &u.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return u, nil
}

func (m *UserModel) List() ([]*User, error) {
	stmt := `SELECT id,username,email FROM users`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	users := []*User{}
	n := 0
	for rows.Next() {
		users = append(users, &User{})
		err := rows.Scan(&users[n].Id, &users[n].Username, &users[n].Email)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrNoRecord
			} else {
				return nil, err
			}
		}
		n++
	}
	return users, nil
}

func (m *UserModel) Update(u User) error {
	stmt := `UPDATE users SET username = ? , email = ? WHERE id = ?`
	res, err := m.DB.Exec(stmt, u.Username, u.Email, u.Id)
	if err != nil {
		return err
	}
	affectedNum, _ := res.RowsAffected()
	if affectedNum == 0 {
		return fmt.Errorf("nothing was updated")
	}
	return nil
}

func (m *UserModel) Delete(id int) error {
	stmt := `DELETE FROM users WHERE id = ?`
	res, err := m.DB.Exec(stmt, id)
	if err != nil {
		return err
	}
	affectedNum, _ := res.RowsAffected()
	if affectedNum == 0 {
		return fmt.Errorf("no user found with id %d", id)
	}
	return nil
}
