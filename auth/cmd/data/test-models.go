package data

import (
	"database/sql"
)

type PostgresTestRepository struct {
	Conn *sql.DB
}

func NewPostgresTestRepository(db *sql.DB) *PostgresTestRepository {
	return &PostgresTestRepository{Conn: db}
}

func (u *PostgresTestRepository) GetAll() ([]*User, error) {
	users := []*User{}
	return users, nil
}

func (u *PostgresTestRepository) GetByEmail(email string) (*User, error) {
	user := User{}
	return &user, nil
}

func (u *PostgresTestRepository) GetOne(id int) (*User, error) {
	user := User{}
	return &user, nil
}

// update user in DB, using information stored in User struct instance
func (u *PostgresTestRepository) Update(user User) error {
	return nil
}

// Delete user using ID specified to func
func (u *PostgresTestRepository) DeleteByID(id int) error {
	return nil
}

func (u *PostgresTestRepository) Insert(user User) (int, error) {
	return 2, nil
}

func (u *PostgresTestRepository) ResetPassword(password string, user User) error {
	return nil
}

func (u *PostgresTestRepository) PasswordMatches(plainText string, user User) (bool, error) {
	return true, nil
}