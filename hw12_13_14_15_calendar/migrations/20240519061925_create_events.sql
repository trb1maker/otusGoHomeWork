-- +goose Up
-- +goose StatementBegin
create table if not exists otus.events (
    event_id    uuid      not null default gen_random_uuid() primary key,
    title       text      not null,
    start_time  timestamp not null,
    end_time    timestamp not null check(end_time >= start_time),
    description text,
    user_id     uuid      references otus.users(user_id),
    notify      interval,
    is_deleted  boolean   not null default false
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists otus.events;
-- +goose StatementEnd
