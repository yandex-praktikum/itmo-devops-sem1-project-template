-- +goose Up
-- +goose StatementBegin
create table if not exists prices
(
    id               int,
    name             text,
    category         text,
    price            numeric(10,2),
    create_date      date,
    constraint pk_products primary key (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table prices;
-- +goose StatementEnd
