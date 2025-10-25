package user

import (
	"context"

	"github.com/Dokhoyan/common/pkg/client/db"
	"github.com/Dokhoyan/go-messenger-microservices/auth/internal/model"
	"github.com/Dokhoyan/go-messenger-microservices/auth/internal/repository"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

const (
	tableName       = "users"

	idColumn        = "id"
	usernameColumn  = "username"
	passwordColumn  = "password_hash"
	roleColumn      = "role"
	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.UserRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, userinfo *model.UserAuthData) (int64, error) {
	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(usernameColumn, passwordColumn, roleColumn).
		Values(userinfo.Username, userinfo.PasswordHash, userinfo.Role).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, err
	}

	q := db.Query{
		Name:     "user_repository.Create",
		QueryRaw: query,
	}

	var id int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *repo) Get(ctx context.Context, userName string) (*model.User, error) {
	builder := sq.Select(idColumn, usernameColumn, passwordColumn, roleColumn, createdAtColumn, updatedAtColumn).
		PlaceholderFormat(sq.Dollar).
		From(tableName).
		Where(sq.Eq{usernameColumn: userName}).
		Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "user_repository.Get",
		QueryRaw: query,
	}

	var user model.User
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(
    &user.ID,
    &user.Info.Username,
    &user.Info.PasswordHash,
    &user.Info.Role,
    &user.CreatedAt,
    &user.UpdatedAt,
)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}

		return nil, err
	}

	return &user, nil
}