package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const dbTimeout = time.Second * 3

type PostgresRepository struct {
	Conn *sql.DB
}

func NewPostgresRepository(pool *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		Conn: pool,
	}
}

// func New(dbPool *sql.DB) Models {
// 	db = dbPool

// 	return Models{
// 		User: User{},
// 	}
// }

// type Models struct {
// 	User User
// }

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	Password  string    `json:"-"` // field is always omitted from JSON
	Active    int       `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var tableColumns [8]string = [8]string{
	"id",
	"email",
	"first_name",
	"last_name",
	"password",
	"user_active",
	"created_at",
	"updated_at",
}

func getTableColString() string {
	return strings.Join(tableColumns[:], ", ")
}

func (u *PostgresRepository) GetAll() ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	cols := getTableColString()
	query := fmt.Sprintf("select %s from users order by last_name", cols)

	rows, err := u.Conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.Password,
			&user.Active,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			log.Println("Error scanning", err)
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}

func (u *PostgresRepository) GetByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	cols := getTableColString()
	query := fmt.Sprintf("select %s from users where email = $1", cols)

	var user User

	row := u.Conn.QueryRowContext(ctx, query, email)
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *PostgresRepository) GetOne(id int) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	cols := getTableColString()
	query := fmt.Sprintf("select %s from users where id = $1", cols)

	var user User
	row := u.Conn.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// update user in DB, using information stored in User struct instance
func (u *PostgresRepository) Update(user User) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `update users set
		email = $1,
		first_name = $2,
		last_name = $3,
		user_active = $4,
		updated_at = $5
		where id = $6
	`

	_, err := u.Conn.ExecContext(ctx, stmt,
		user.Email,
		user.FirstName,
		user.LastName,
		user.Active,
		time.Now(),
		user.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

// Delete user using ID specified to func
func (u *PostgresRepository) DeleteByID(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := "delete from users where id = $1"

	_, err := u.Conn.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	return nil
}

func (u *PostgresRepository) Insert(user User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return 0, err
	}

	var newID int
	stmt := `insert into users (email, first_name, last_name, password, user_active, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6, $7) returning id`

	err = u.Conn.QueryRowContext(ctx, stmt,
		user.Email,
		user.FirstName,
		user.LastName,
		hashedPassword,
		user.Active,
		time.Now(),
		time.Now(),
	).Scan(&newID)
	if err != nil {
		return 0, err
	}

	return newID, nil
}

func (u *PostgresRepository) ResetPassword(password string, user User) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := "update users set password = $1 where id = $2"
	_, err = u.Conn.ExecContext(ctx, stmt, hashedPassword, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (u *PostgresRepository) PasswordMatches(plainText string, user User) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(plainText))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			// passwords don't match
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}
