-- +goose Up
create table chats
(
    id         serial primary key,
    name       varchar   not null,
    created_at timestamp not null default now()
);

create table chat_users
(
    chat_id    int not null,
    user_id    int not null,
    created_at timestamp not null default now()
);

-- +goose Down
drop table chats;
drop table chat_users;