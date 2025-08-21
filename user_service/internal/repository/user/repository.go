package user

import (
	//"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/repository"
	"context"

	sq "github.com/Masterminds/squirrel"

	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/model"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/repository"
	//converter "github.com/Dokhoyan/go-messenger-microservices/user_service/internal/repository/user/conventer"

	//modelRepo "github.com/Dokhoyan/go-messenger-microservices/user_service/internal/repository/user/model"
	"github.com/jackc/pgx/v4/pgxpool"
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
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) repository.UserRepository {
	return &repo{db: db}
}


func (r *repo) Create(ctx context.Context, info *model.UserInfo) (int64, error) {
	builder := sq.Insert(tableName).
		//PlaceholderFormat(sq.Dollar).
		Columns(nameColumn, usernameColumn, emailColumn, birth_dateColumn, avatar_urlColumn).
		Values(info.Name, info.Username, info.Email, info.Birth_date, info.Avatar_url).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, err
	}

	var id int64
	err = r.db.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}


func (r *repo) Get(ctx context.Context, id int64) (*model.User, error) {
	builder := sq.Select(idColumn, nameColumn, usernameColumn, emailColumn, birth_dateColumn, avatar_urlColumn, createdAtColumn, updatedAtColumn).
		//PlaceholderFormat(sq.Dollar).
		From(tableName).
		Where(sq.Eq{idColumn: id}).
		Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	var user model.User
	err = r.db.QueryRow(ctx, query, args...).Scan(&user.ID, &user.Info.Name, &user.Info.Username, &user.Info.Email,
		&user.Info.Birth_date, &user.Info.Avatar_url, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	//return converter.ToUserFromRepo(&user), nil 
	return &user, nil
}