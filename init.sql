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
    id      serial primary key,
    user_id int   not null references users on delete cascade,
    info    jsonb not null
);
alter table users
    add foreign key (default_address_id) references addresses on delete set null;

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

create table if not exists requests
(
    id          serial primary key,
    user_id     int    not null references users on delete cascade,
    category_id int    not null references categories on delete cascade,
    name        text   not null,
    item_name   text   not null,
    description text   not null,
    photos      text[] not null,
    price       int    not null,
    quantity    int    not null,
    date        text   not null
);

create table if not exists requests_tags
(
    request_id int not null references requests on delete cascade,
    tag_id     int not null references tags on delete cascade
);

create or replace view requests_view as
SELECT r.id,
       r.user_id,
       r.category_id,
       r.name,
       r.item_name,
       r.description,
       r.photos,
       r.price,
       r.quantity,
       r.date,
       jsonb_agg(jsonb_build_object(
               'id', t.id,
               'name', t.name
           )) as tags
FROM requests r
         LEFT JOIN requests_tags rt on r.id = rt.request_id
         JOIN tags t on rt.tag_id = t.id
GROUP BY r.id;

INSERT INTO categories (id, name, photo)
VALUES (1, 'Недвижимость', '/images/property.jpeg');
INSERT INTO categories (id, name, photo)
VALUES (2, 'Электроника', '/images/electronics.jpeg');
INSERT INTO categories (id, name, photo)
VALUES (3, 'Хобби и отдых', '/images/hobby.jpeg');
INSERT INTO categories (id, name, photo)
VALUES (4, 'Транспорт', '/images/transport.png');
INSERT INTO categories (id, name, photo)
VALUES (5, 'Одежда', '/images/clothes.png');
INSERT INTO categories (id, name, photo)
VALUES (6, 'Животные', '/images/pets.jpeg');
INSERT INTO categories (id, name, photo)
VALUES (7, 'Для дома', '/images/house.png');
INSERT INTO categories (id, name, photo)
VALUES (8, 'Прочее', '/images/other.jpeg');

INSERT INTO tags (category_id, name)
VALUES (1, 'Студия');
INSERT INTO tags (category_id, name)
VALUES (1, 'Однокомнатная');
INSERT INTO tags (category_id, name)
VALUES (1, 'Двухкомнатная');
INSERT INTO tags (category_id, name)
VALUES (1, 'Трёхкомнатная');
INSERT INTO tags (category_id, name)
VALUES (2, 'Ноутбук');
INSERT INTO tags (category_id, name)
VALUES (2, 'Телефон');
INSERT INTO tags (category_id, name)
VALUES (2, 'Наушники');
INSERT INTO tags (category_id, name)
VALUES (2, 'Видеокарта');
INSERT INTO tags (category_id, name)
VALUES (3, 'Лыжи');
INSERT INTO tags (category_id, name)
VALUES (3, 'Ролики');
INSERT INTO tags (category_id, name)
VALUES (3, 'Велосипед');
INSERT INTO tags (category_id, name)
VALUES (3, 'Коньки');
INSERT INTO tags (category_id, name)
VALUES (4, 'Автомобиль');
INSERT INTO tags (category_id, name)
VALUES (4, 'Автобус');
INSERT INTO tags (category_id, name)
VALUES (4, 'Электросамокат');
INSERT INTO tags (category_id, name)
VALUES (4, 'Самолёт');
INSERT INTO tags (category_id, name)
VALUES (5, 'Обувь');
INSERT INTO tags (category_id, name)
VALUES (5, 'Верхняя одежда');
INSERT INTO tags (category_id, name)
VALUES (5, 'Рубашки');
INSERT INTO tags (category_id, name)
VALUES (5, 'Штаны');
INSERT INTO tags (category_id, name)
VALUES (6, 'Игрушка');
INSERT INTO tags (category_id, name)
VALUES (6, 'Одежда');
INSERT INTO tags (category_id, name)
VALUES (6, 'Корм');
INSERT INTO tags (category_id, name)
VALUES (6, 'Ошейник');
INSERT INTO tags (category_id, name)
VALUES (7, 'Стул');
INSERT INTO tags (category_id, name)
VALUES (7, 'Стол');
INSERT INTO tags (category_id, name)
VALUES (7, 'Кровать');
INSERT INTO tags (category_id, name)
VALUES (7, 'Шкаф');
INSERT INTO tags (category_id, name)
VALUES (null, 'Синий');
INSERT INTO tags (category_id, name)
VALUES (null, 'Красный');
INSERT INTO tags (category_id, name)
VALUES (null, 'Зелёный');
INSERT INTO tags (category_id, name)
VALUES (null, 'Белый');
INSERT INTO tags (category_id, name)
VALUES (null, 'Маленький');
INSERT INTO tags (category_id, name)
VALUES (null, 'Средний');
INSERT INTO tags (category_id, name)
VALUES (null, 'Большой');
INSERT INTO tags (category_id, name)
VALUES (null, 'Огромный');




