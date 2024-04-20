package auth

import (
	"context"
	"database/sql"
	"errors"

	"github.com/KretovDmitry/gophermart-loyalty-service/internal/models/errs"
	"github.com/KretovDmitry/gophermart-loyalty-service/internal/models/user"
	"github.com/KretovDmitry/gophermart-loyalty-service/pkg/logger"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type Repository interface {
	GetUserByID(ctx context.Context, userID int) (*user.User, error)
	GetUserByLogin(ctx context.Context, login string) (*user.User, error)
	CreateUser(ctx context.Context, login, password string) (id int, err error)
}

type repository struct {
	db     *sql.DB
	logger logger.Logger
}

func NewRepository(db *sql.DB, logger logger.Logger) (*repository, error) {
	if db == nil {
		return nil, errors.New("nil dependency: database")
	}

	return &repository{
		db:     db,
		logger: logger,
	}, nil
}

var _ Repository = (*repository)(nil)

func (r *repository) GetUserByID(ctx context.Context, userID int) (*user.User, error) {
	const query = "SELECT * FROM users WHERE id = $1"

	u := new(user.User)

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&u.ID,
		&u.Login,
		&u.Password,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		return nil, err
	}

	return u, nil
}

func (r *repository) GetUserByLogin(ctx context.Context, login string) (*user.User, error) {
	const query = "SELECT * FROM users WHERE login = $1"

	u := new(user.User)

	err := r.db.QueryRowContext(ctx, query, login).Scan(
		&u.ID,
		&u.Login,
		&u.Password,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		return nil, err
	}

	return u, nil
}

func (r *repository) CreateUser(ctx context.Context, login, password string) (int, error) {
	const query = "INSERT INTO users (login, password) VALUES ($1, $2) RETURNING id"

	var id int

	err := r.db.QueryRowContext(ctx, query, login, password).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			r.logger.Errorf("%v", pgErr.Code)
			if pgErr.Code == pgerrcode.UniqueViolation {
				return -1, errs.ErrConflict
			}
		}
		return -1, err
	}

	return id, nil
}
