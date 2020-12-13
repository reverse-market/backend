create table if not exists users
(
    id                 serial primary key,
    name               text        not null,
    email              text unique not null,
    photo              text        not null default '',
    default_address_id int
);

create table if not exists addresses
(
    id serial primary key,
    user_id int not null references users on delete cascade,
    info jsonb not null
);
alter table users add foreign key (default_address_id) references addresses on delete set null;

create table if not exists categories
(
    id    serial primary key,
    name  text not null,
    photo text not null
);

create table if not exists tags
(
    id          serial primary key,
    category_id int REFERENCES categories ON DELETE CASCADE,
    name        text not null
);




