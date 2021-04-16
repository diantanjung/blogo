package psql_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"

	userPsqlRepo "github.com/diantanjung/blogo/user-service/user/repository/psql"
	"github.com/diantanjung/blogo/user-service/domain"
	"github.com/diantanjung/blogo/user-service/user/repository"
)

func TestFetch(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mockUsers := []domain.User{
		domain.User{
			ID: 1, Username: "user1", Name: "user 1", Email:"user1@gmail.com",
			UpdatedAt: time.Now(), CreatedAt: time.Now(),
		},
		domain.User{
			ID: 1, Username: "user2", Name: "user 2", Email:"user2@gmail.com",
			UpdatedAt: time.Now(), CreatedAt: time.Now(),
		},
	}

	rows := sqlmock.NewRows([]string{"id", "username", "name", "email", "updated_at", "created_at"}).
		AddRow(mockUsers[0].ID, mockUsers[0].Username, mockUsers[0].Name,
			mockUsers[0].Email, mockUsers[0].UpdatedAt, mockUsers[0].CreatedAt).
			AddRow(mockUsers[1].ID, mockUsers[1].Username, mockUsers[1].Name,
				mockUsers[1].Email, mockUsers[1].UpdatedAt, mockUsers[1].CreatedAt)

	query := "SELECT id, username, name, email, created_at, updated_at  FROM users WHERE created_at > \\? ORDER BY created_at LIMIT \\?"

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := userPsqlRepo.NewPsqlUserRepository(db)
	cursor := repository.EncodeCursor(mockUsers[1].CreatedAt)
	num := int64(2)
	list, nextCursor, err := a.Fetch(context.TODO(), cursor, num)
	assert.NotEmpty(t, nextCursor)
	assert.NoError(t, err)
	assert.Len(t, list, 2)
}

func TestGetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows := sqlmock.NewRows([]string{"id", "username", "name", "email", "created_at", "updated_at"}).
		AddRow(1, "usrname1", "Name 1", "username1@gmail.com", time.Now(), time.Now())

	query := "SELECT id, username, name, email, created_at, updated_at FROM users WHERE ID = \\?"

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := userPsqlRepo.NewPsqlUserRepository(db)

	num := int64(5)
	anUser, err := a.GetByID(context.TODO(), num)
	assert.NoError(t, err)
	assert.NotNil(t, anUser)
}

func TestStore(t *testing.T) {
	now := time.Now()
	u := &domain.User{
		Username:   "username1",
		Name:   	"Nama1",
		Email:   	"email1@gmail.com",
		Password:	"asdf123",
		CreatedAt: now,
		UpdatedAt: now,
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	query := "INSERT INTO users \\(username,name,email,password,created_at,updated_at\\) VALUE \\(username=\\?, name=\\?, email=\\?, password=\\?, created_at=\\?, updated_at=\\?\\)"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(u.Username, u.Name, u.Email, u.Password, u.CreatedAt, u.UpdatedAt).WillReturnResult(sqlmock.NewResult(12, 1))

	a := userPsqlRepo.NewPsqlUserRepository(db)

	err = a.Store(context.TODO(), u)
	assert.NoError(t, err)
	assert.Equal(t, int64(12), u.ID)
}

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	query := "DELETE FROM users WHERE id = \\?"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(12).WillReturnResult(sqlmock.NewResult(12, 1))

	a := userPsqlRepo.NewPsqlUserRepository(db)

	num := int64(12)
	err = a.Delete(context.TODO(), num)
	assert.NoError(t, err)
}

func TestUpdate(t *testing.T) {
	now := time.Now()
	u := &domain.User{
		Username:   "username1",
		Name:   	"Nama1",
		Email:   	"email1@gmail.com",
		Password:	"asdf123",
		CreatedAt: now,
		UpdatedAt: now,
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	query := "UPDATE users SET username=\\?, password=\\?, name=\\?, email=\\?, updated_at=\\? WHERE id=\\?"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(u.Username, u.Password, u.Name, u.Email, u.UpdatedAt, u.ID).WillReturnResult(sqlmock.NewResult(12, 1))

	a := userPsqlRepo.NewPsqlUserRepository(db)

	err = a.Update(context.TODO(), u)
	assert.NoError(t, err)
}