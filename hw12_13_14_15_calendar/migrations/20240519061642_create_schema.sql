-- +goose Up
-- +goose StatementBegin
create schema if not exists otus;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop schema if exists otus;
-- +goose StatementEnd
