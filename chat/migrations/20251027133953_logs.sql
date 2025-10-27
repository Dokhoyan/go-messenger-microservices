-- +goose Up
-- +goose StatementBegin
create table chats_log
(
    id         serial primary key,
    chat_id    int       not null,
    action     varchar   not null,
    created_at timestamp not null default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table chats_log;
-- +goose StatementEnd