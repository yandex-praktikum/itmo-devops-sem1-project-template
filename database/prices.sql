CREATE TABLE IF NOT EXISTS prices
(
    id integer not null primary key,
    name text not null,
    category text not null,
    price numeric(10, 2) not null,
    create_date date not null
);