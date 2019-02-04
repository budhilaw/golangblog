package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/budhilaw/blog/models"
	"github.com/budhilaw/blog/user"
	"golang.org/x/crypto/bcrypt"
)

// MysqlUserRepo : MySQL DB on User Repository
type MysqlUserRepo struct {
	Conn *sql.DB
}

// New : Will create object that represent of the user.Repository interface
func New(Conn *sql.DB) user.Repository {
	return &MysqlUserRepo{Conn}
}

// Fetch : Fetch user data
func (m *MysqlUserRepo) Fetch(ctx context.Context, num int64) ([]*models.User, error) {
	query := "SELECT * FROM users ORDER BY created_at LIMIT ?"
	rows, err := m.Conn.QueryContext(ctx, query, num)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*models.User, 0)
	for rows.Next() {
		t := new(models.User)
		err = rows.Scan(
			&t.ID,
			&t.Name,
			&t.Email,
			&t.Password,
			&t.CreadedAt,
			&t.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		result = append(result, t)
	}
	return result, nil
}

// Store : Add a new user
func (m *MysqlUserRepo) Store(ctx context.Context, a *models.User) error {
	query := "INSERT INTO users SET name=?, email=?, password=?, updated_at=?, created_at=?"
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer stmt.Close()

	pwd := getPwd(a.Password)
	hash := hashAndSalt(pwd)

	a.CreadedAt = time.Now()
	a.UpdatedAt = time.Now()

	log.Println("Created at: ", a.CreadedAt)
	res, err := stmt.ExecContext(ctx, a.Name, a.Email, hash, a.UpdatedAt, a.CreadedAt)
	if err != nil {
		log.Fatal(err)
		return err
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
		return err
	}
	a.ID = lastID
	return nil
}

// Login : Login function
func (m *MysqlUserRepo) Login(ctx context.Context, email string, pass string, a *models.User) error {
	query := "SELECT * FROM users WHERE email=?"
	err := m.Conn.QueryRowContext(ctx, query, email).Scan(
		&a.ID,
		&a.Name,
		&a.Email,
		&a.Password,
		&a.UpdatedAt,
		&a.CreadedAt,
	)

	if err == sql.ErrNoRows {
		err = errors.New("No user found")
		return err
	}

	pwd := getPwd(pass)
	_, err = comparePwd(a.Password, pwd)
	if err != nil {
		err = errors.New("Password is wrong")
		return err
	}

	return nil
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
