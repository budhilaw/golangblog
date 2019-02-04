package repository_test

import (
	"context"
	"database/sql/driver"
	"log"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/budhilaw/blog/models"
	userRepo "github.com/budhilaw/blog/user/repository"
	"github.com/stretchr/testify/assert"
)

type AnyTime struct{}
type Anything struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func (a Anything) Match(v driver.Value) bool {
	if s, ok := v.(string); ok {
		return len(s) > 0
	}
	return false
}

func TestFetch(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	now := time.Now()
	mockUser := []models.User{
		models.User{
			ID:        1,
			Name:      "Kai",
			Email:     "kai@gmail.com",
			Password:  "1550",
			CreadedAt: now,
			UpdatedAt: now,
		},
		models.User{
			ID:        2,
			Name:      "Eric",
			Email:     "eric@gmail.com",
			Password:  "12345",
			CreadedAt: now,
			UpdatedAt: now,
		},
	}

	rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "created_at", "updated_at"}).
		AddRow(mockUser[0].ID, mockUser[0].Name, mockUser[0].Email,
			mockUser[0].Password, mockUser[0].CreadedAt, mockUser[0].UpdatedAt).
		AddRow(mockUser[1].ID, mockUser[1].Name, mockUser[1].Email,
			mockUser[1].Password, mockUser[1].CreadedAt, mockUser[1].UpdatedAt)

	query := "SELECT * FROM users ORDER BY created_at LIMIT 2"

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := userRepo.New(db)
	num := int64(2)
	list, err := a.Fetch(context.TODO(), num)
	assert.NoError(t, err)
	assert.Len(t, list, 2)

}

func TestStore(t *testing.T) {
	now := time.Now()
	ur := &models.User{
		Name:      "Kai",
		Email:     "kai@gmail.com",
		Password:  "1550",
		CreadedAt: now,
		UpdatedAt: now,
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	query := "INSERT INTO users SET name=\\?, email=\\?, password=\\?, updated_at=\\?, created_at=\\?"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(ur.Name, ur.Email, Anything{}, AnyTime{}, AnyTime{}).WillReturnResult(sqlmock.NewResult(12, 1))

	a := userRepo.New(db)

	err = a.Store(context.TODO(), ur)
	assert.NoError(t, err)
	assert.Equal(t, int64(12), ur.ID)
}

func getPwd(pwd string) []byte {
	return []byte(pwd)
}

func hashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Fatal(err)
	}

	return string(hash)
}

func comparePwd(hashedPwd string, plainPwd []byte) (bool, error) {
	byteHash := []byte(hashedPwd)

	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		return false, err
	}
	return true, nil
}
