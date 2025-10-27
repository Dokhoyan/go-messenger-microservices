package chat

import (
	"context"
	"errors"

	
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"github.com/Dokhoyan/common/pkg/client/db"
	
	"github.com/Dokhoyan/go-messenger-microservices/chat_service/internal/repository"
	"github.com/Dokhoyan/go-messenger-microservices/chat_service/pkg/api/chat_v1"
)

const (
	tableChatsName      = "chats"
	tableChatsUsersName = "chat_users"

	idColumn     = "id"
	chatIDColumn = "chat_id"
	userIDColumn = "user_id"
	nameColumn   = "name"
)

type repo struct {
	db db.Client
}

// NewRepository creates new repository object
func NewRepository(db db.Client) repository.ChatRepository {
	return &repo{db: db}
}

// CreateChat creates new chat with specified parameters
func (r *repo) CreateChat(ctx context.Context, req *chat_v1.CreateChatRequest) (int64, error) {
	chatID, err := r.createNewChat(ctx, req)
	if err != nil {
		return 0, err
	}

	err = r.createChatUsers(ctx, req, chatID)
	if err != nil {
		return 0, err
	}

	return chatID, nil
}

// DeleteChat removes chat
func (r *repo) DeleteChat(ctx context.Context, id int64) error {
	err := r.deleteChatUsers(ctx, id)
	if err != nil {
		return err
	}

	err = r.deleteChat(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

// GetChatIDByName returns chat id by name
func (r *repo) GetChatIDByName(ctx context.Context, chatname string) (int64, error) {
	builderSelect := sq.Select(idColumn).
		From(tableChatsName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{nameColumn: chatname})

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return 0, err
	}

	q := db.Query{
		Name:     "chat_repository.GetChatIDByName",
		QueryRaw: query,
	}

	var chatID int64

	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&chatID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, errors.New("chat IDs not found")
		}
		return 0, err
	}

	return chatID, nil
}

func (r *repo) createNewChat(ctx context.Context, req *chat_v1.CreateChatRequest) (int64, error) {
	builderInsert := sq.Insert(tableChatsName).
		PlaceholderFormat(sq.Dollar).
		Columns(nameColumn).
		Values(req.Name).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return 0, err
	}

	q := db.Query{
		Name:     "chat_repository.CreateChat",
		QueryRaw: query,
	}

	var chatID int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&chatID)
	if err != nil {
		return 0, err
	}

	return chatID, nil
}

func (r *repo) createChatUsers(ctx context.Context, req *chat_v1.CreateChatRequest, chatID int64) error {
	builderInsertChatUser := sq.Insert(tableChatsUsersName).
		PlaceholderFormat(sq.Dollar).
		Columns(userIDColumn, chatIDColumn)

	for _, id := range req.UserIds {
		builderInsertChatUser = builderInsertChatUser.Values(id, chatID)
	}

	query, args, err := builderInsertChatUser.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "chat_repository.CreateChatUser",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) deleteChatUsers(ctx context.Context, id int64) error {
	builderDelete := sq.Delete(tableChatsUsersName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{chatIDColumn: id})

	query, args, err := builderDelete.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "chat_repository.DeleteChatUsers",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) deleteChat(ctx context.Context, id int64) error {
	builderChatDelete := sq.Delete(tableChatsName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id})

	query, args, err := builderChatDelete.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "chat_repository.DeleteChat",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return nil
}