create table if not exists users
(
    id                 serial primary key,
    name               text        not null,
    email              text unique not null,
    photo              text        not null default '',
    default_address_id int
)
