package psql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/diantanjung/blogo/user-service/domain"
	"github.com/diantanjung/blogo/user-service/user/repository"
)

type psqlUserRepository struct {
	Conn *sql.DB
}

// NewMysqlUserRepository will create an object that represent the user.Repository interface
func NewPsqlUserRepository(Conn *sql.DB) domain.UserRepository {
	return &psqlUserRepository{Conn}
}

func (m *psqlUserRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.User, err error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			logrus.Error(errRow)
		}
	}()

	result = make([]domain.User, 0)
	for rows.Next() {
		user := domain.User{}
		err = rows.Scan(
			&user.ID,
			&user.Username,
			&user.Name,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, user)
	}

	return result, nil
}

func (m *psqlUserRepository) Fetch(ctx context.Context, cursor string, num int64) (res []domain.User, nextCursor string, err error) {
	query := `SELECT id, username, name, email, created_at, updated_at
  						FROM users WHERE created_at > ? ORDER BY created_at LIMIT ? `

	decodedCursor, err := repository.DecodeCursor(cursor)
	if err != nil && cursor != "" {
		return nil, "", domain.ErrBadParamInput
	}

	res, err = m.fetch(ctx, query, decodedCursor, num)
	if err != nil {
		return nil, "", err
	}

	if len(res) == int(num) {
		nextCursor = repository.EncodeCursor(res[len(res)-1].CreatedAt)
	}

	return
}
func (m *psqlUserRepository) GetByID(ctx context.Context, id int64) (res domain.User, err error) {
	query := `SELECT id, username, name, email, created_at, updated_at
  						FROM users WHERE ID = ?`

	list, err := m.fetch(ctx, query, id)
	if err != nil {
		return domain.User{}, err
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, domain.ErrNotFound
	}

	return
}

func (m *psqlUserRepository) Update(ctx context.Context, u *domain.User) (err error) {

	query := `UPDATE users SET username=?, password=?, name=?, email=?, updated_at=? WHERE id=?`

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, u.Username, u.Password, u.Name, u.Email, u.UpdatedAt, u.ID)
	if err != nil {
		return
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return
	}
	if affect != 1 {
		err = fmt.Errorf("Weird  Behavior. Total Affected: %d", affect)
		return
	}

	return
}

func (m *psqlUserRepository) Store(ctx context.Context, u *domain.User) (err error) {
	query := `INSERT INTO users (username,name,email,password,created_at,updated_at) VALUE (username=?, name=?, email=?, password=?, created_at=?, updated_at=?)`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, u.Username, u.Name, u.Email, u.Password, u.CreatedAt, u.UpdatedAt)
	if err != nil {
		return
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		return
	}
	u.ID = lastID
	return
}

func (m *psqlUserRepository) Delete(ctx context.Context, id int64) (err error) {
	query := "DELETE FROM users WHERE id = ?"

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return
	}

	rowsAfected, err := res.RowsAffected()
	if err != nil {
		return
	}

	if rowsAfected != 1 {
		err = fmt.Errorf("Weird  Behavior. Total Affected: %d", rowsAfected)
		return
	}

	return
}
