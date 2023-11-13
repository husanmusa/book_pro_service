-- +goose Up
-- +goose StatementBegin
create table book_category
(
    id   uuid primary key,
    name varchar
);

create table book
(
    id               uuid primary key,
    name             varchar,
    author_name      varchar,
    pages            int,
    description      varchar(500),
    book_category_id uuid references book_category (id),
    created_at       timestamp default now(),
    updated_at       timestamp default now(),
    deleted_at       bigint
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table book;
drop table book_category;
-- +goose StatementEnd
