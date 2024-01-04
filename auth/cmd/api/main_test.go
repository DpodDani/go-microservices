package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DpodDani/auth/cmd/data"
	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/require"
)

const TIMEOUT = 5

func createServer(db *sql.DB) *httptest.Server {
	config := Config{
		DB:     db,
		Models: data.New(db),
	}
	server := httptest.NewServer(config.routes())
	return server
}

func createDB() *sql.DB {
	// TODO: Connect to DB and create temp database
	return connectToDB("host=localhost port=5432 user=postgres password=password dbname=users sslmode=disable timezone=UTC connect_timeout=5")
}

func mustCreateTestUser(t *testing.T, db *sql.DB) *data.User {
	config := Config{
		DB:     db,
		Models: data.New(db),
	}

	now := time.Now()
	expectedPassword := "password"

	userID, err := config.Models.User.Insert(data.User{
		Email:     fmt.Sprintf("test%d@email.com", now.Unix()),
		FirstName: "test",
		LastName:  "user",
		Password:  expectedPassword,
		Active:    1,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		require.NoError(t, err, "Failed to create test user")
	}

	user, err := config.Models.User.GetOne(userID)
	if err != nil {
		require.NoError(t, err, "Failed to fetch test user")
	}

	user.Password = expectedPassword

	return user
}

func TestAuthenticate_when_user_does_not_exist(t *testing.T) {
	db := createDB()
	server := createServer(db)
	defer server.Close()

	e := httpexpect.Default(t, server.URL)

	e.POST("/authenticate").
		WithFormField("email", "fake_user@email.com").
		WithFormField("password", "fake_password").
		Expect().
		Status(http.StatusBadRequest)
}

func TestAuthenticate_password_matches(t *testing.T) {
	db := createDB()
	server := createServer(db)
	defer server.Close()

	user := mustCreateTestUser(t, db)
	e := httpexpect.Default(t, server.URL)

	e.POST("/authenticate").
		WithJSON(map[string]interface{}{
			"email":    user.Email,
			"password": user.Password,
		}).
		Expect().
		Status(http.StatusAccepted)
}
