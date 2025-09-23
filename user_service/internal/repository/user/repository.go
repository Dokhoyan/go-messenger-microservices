package user

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Dokhoyan/common/pkg/filter"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/client/db"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/model"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/repository"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

const (
	tableName = "users"
	
	idColumn         = "id"
	nameColumn       = "name"
	usernameColumn   = "username"
	emailColumn      = "email"
	birth_dateColumn = "birth_date"
	avatar_urlColumn = "avatar_url"
	roleColumn      = "role"
	passwordColumn  = "password"
	createdAtColumn  = "created_at"
	updatedAtColumn  = "updated_at"
)


type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.UserRepository {
	return &repo{db: db}
}


func (r *repo) Create(ctx context.Context, params *model.UserCreate) (int64, error) {
	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(nameColumn, usernameColumn, emailColumn, birth_dateColumn, avatar_urlColumn, roleColumn, passwordColumn).
		Values(params.Info.Name, params.Info.Username, params.Info.Email, params.Info.Birth_date, params.Info.Avatar_url, params.Info.Role, params.Password).
		Suffix(fmt.Sprintf("RETURNING %s", idColumn))

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, errors.Errorf("error at parse sql builder: %v", err)
	}

	q:=db.Query{
		Name: "user_Repository Create",
		QueryRaw: query,
	}

	var id int64

	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}


func (r *repo) Get(ctx context.Context, filters filter.Filter) (*model.User, error) {
	builder := sq.Select(idColumn, nameColumn, usernameColumn, emailColumn, birth_dateColumn, avatar_urlColumn, roleColumn, createdAtColumn, updatedAtColumn, passwordColumn).
		PlaceholderFormat(sq.Dollar).
		From(tableName)


	for _, condition := range filters.Conditions {
		builder = builder.Where(sq.Eq{condition.Key: condition.Value})
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, errors.Errorf("error at parse sql builder: %v", err)
	}

	q := db.Query{
		Name:     "user_repository.Get",
		QueryRaw: query,
	}

	var user model.User

	err = r.db.DB().ScanOneContext(ctx, &user, q, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.Errorf("user not found")
		}
		return nil, errors.Errorf("error at query to database: %v", err)
	}

	return &user, nil 

}