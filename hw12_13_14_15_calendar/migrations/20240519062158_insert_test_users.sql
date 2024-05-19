-- +goose Up
-- +goose StatementBegin
insert into otus.users values('6a09a9aa-6f09-4ea9-be76-5f9c8470e7ea');
insert into otus.users values('ef6acd9c-7385-420d-b408-f0029a9decc3');
insert into otus.users values('3f65e509-7a2a-4aed-adb3-90cf08b2917a');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
delete from otus.users;
-- +goose StatementEnd
