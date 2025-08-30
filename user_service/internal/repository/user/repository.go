package user

import (
	//"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/repository"
	"context"

	sq "github.com/Masterminds/squirrel"

	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/client/db"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/repository/user/conventer"

	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/model"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/repository"
	//converter "github.com/Dokhoyan/go-messenger-microservices/user_service/internal/repository/user/conventer"
	modelRepo "github.com/Dokhoyan/go-messenger-microservices/user_service/internal/repository/user/model"
)

const (
	tableName = "users"
	
	idColumn         = "id"
	nameColumn       = "name"
	usernameColumn   = "username"
	emailColumn      = "email"
	birth_dateColumn = "birth_date"
	avatar_urlColumn = "avatar_url"
	createdAtColumn  = "created_at"
	updatedAtColumn  = "updated_at"
)


type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.UserRepository {
	return &repo{db: db}
}


func (r *repo) Create(ctx context.Context, info *model.UserInfo) (int64, error) {
	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(nameColumn, usernameColumn, emailColumn, birth_dateColumn, avatar_urlColumn).
		Values(info.Name, info.Username, info.Email, info.Birth_date, info.Avatar_url).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, err
	}

	q:=db.Query{
		Name: "user Repository Create",
		QueryRaw: query,
	}
	var id int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}


func (r *repo) Get(ctx context.Context, id int64) (*model.User, error) {
	builder := sq.Select(idColumn, nameColumn, usernameColumn, emailColumn, birth_dateColumn, avatar_urlColumn).
		PlaceholderFormat(sq.Dollar).
		From(tableName).
		Where(sq.Eq{idColumn: id}).
		Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "user_repository.Get",
		QueryRaw: query,
	}

	var user modelRepo.User
	err = r.db.DB().ScanOneContext(ctx, &user, q, args...)
	if err != nil {
		return nil, err
	}

	return converter.ToUserFromRepo(&user), nil 

}