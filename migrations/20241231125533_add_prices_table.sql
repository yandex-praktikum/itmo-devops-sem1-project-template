-- +goose Up
-- +goose StatementBegin
CREATE TABLE prices (
    id integer primary key,
    name varchar(250),
    category varchar(250),
    price numeric(10, 2),
    create_date date
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE prices;
-- +goose StatementEnd
