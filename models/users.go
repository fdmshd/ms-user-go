package models

import "database/sql"

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

	statement := `INSERT INTO users (username, email, password)
	VALUES(?,?,?)`
	res, err := m.DB.Exec(statement, u.Username, u.Email, u.Password)
	if err != nil {
		return 0, nil
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (m *UserModel) Get(id int) (*User, error) {
	return nil, nil
}

func (m *UserModel) GetByName()

func (m *UserModel) Update(id int) (*User, error) {
	return nil, nil
}

func (m *UserModel) Delete(id int) error {
	return nil
}
