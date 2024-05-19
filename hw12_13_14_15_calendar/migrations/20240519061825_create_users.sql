-- +goose Up
-- +goose StatementBegin
create table if not exists otus.users (
    user_id   uuid   not null primary key
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists otus.users;
-- +goose StatementEnd
