package user

import (
	"context"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	pgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const userTableName = "users"

var (
	psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	ErrUserExists = errors.New("user already exists")
)

type repository struct {
	db *pgx.Conn
}

func newRepository(db *pgx.Conn) *repository {
	return &repository{db: db}
}

func (r *repository) saveUser(ctx context.Context, user User) error {
	query, args, err := psql.
		Insert(userTableName).
		Columns("email", "username", "password_hash").
		Values(user.Email, user.Username, user.PasswordHash).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		if isDuplicateKeyError(err) {
			return ErrUserExists
		}
		return err
	}

	return nil
}

func (r *repository) getUserByEmail(ctx context.Context, email string) (User, error) {
	query, args, err := psql.
		Select("id", "email", "username", "password_hash").
		From(userTableName).
		Where(squirrel.Eq{"email": email}).
		ToSql()
	if err != nil {
		return User{}, err
	}

	var user User
	err = pgxscan.Get(ctx, r.db, &user, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, pgx.ErrNoRows
		}
		return User{}, err
	}

	return user, nil
}

func (r *repository) getUserByID(ctx context.Context, id int64) (User, error) {
	query, args, err := psql.
		Select("id", "email", "username", "password_hash").
		From(userTableName).
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return User{}, err
	}

	var user User
	err = pgxscan.Get(ctx, r.db, &user, query, args...)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

// UpdateUser обновляет данные пользователя
func (r *repository) updateUser(ctx context.Context, user User) error {
	query, args, err := psql.
		Update(userTableName).
		Set("email", user.Email).
		Set("username", user.Username).
		Set("password_hash", user.PasswordHash).
		Where(squirrel.Eq{"id": user.ID}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, query, args...)
	return err
}

func (r *repository) deleteUser(ctx context.Context, id int64) error {
	query, args, err := psql.
		Delete(userTableName).
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, query, args...)
	return err
}

func isDuplicateKeyError(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505" // unique_violation
	}
	return false
}
