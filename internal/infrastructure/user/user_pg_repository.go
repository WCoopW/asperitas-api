package user

import (
	"context"
	"database/sql"
	"errors"

	domain "reddit/internal/domain/user"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

const pgUniqueViolation = "23505"

type UserPGRepository struct {
	db *sqlx.DB
}

func NewUserPGRepository(db *sqlx.DB) domain.UserRepository {
	return &UserPGRepository{db: db}
}

func (r *UserPGRepository) GetUserByUsername(username string) (domain.User, error) {
	return r.getUser(
		`SELECT id, username, password_hash FROM users WHERE username = $1`,
		username,
	)
}

func (r *UserPGRepository) GetUserByID(id string) (domain.User, error) {
	return r.getUser(
		`SELECT id, username, password_hash FROM users WHERE id = $1`,
		id,
	)
}

func (r *UserPGRepository) getUser(query string, arg string) (domain.User, error) {
	var u domain.User
	err := r.db.QueryRowContext(context.Background(), query, arg).Scan(
		&u.ID, &u.Username, &u.PasswordHash,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, domain.ErrNotFound
		}
		return domain.User{}, err
	}
	return u, nil
}

func (r *UserPGRepository) CreateUser(u domain.User) (domain.User, error) {
	err := r.db.QueryRowContext(
		context.Background(),
		`INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id`,
		u.Username, u.PasswordHash,
	).Scan(&u.ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
			return domain.User{}, domain.ErrUsernameTaken
		}
		return domain.User{}, err
	}
	return u, nil
}
