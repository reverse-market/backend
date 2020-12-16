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
    date        date   not null default current_date,
    finished    bool   not null default false
);

create table if not exists requests_tags
(
    request_id int not null references requests on delete cascade,
    tag_id     int not null references tags on delete cascade
);

create or replace view requests_view as
SELECT r.id,
       r.user_id,
       u.name as username,
       r.category_id,
       r.name,
       r.item_name,
       r.description,
       r.photos,
       r.price,
       r.quantity,
       r.date,
       r.finished,
       jsonb_agg(jsonb_build_object(
               'id', t.id,
               'name', t.name
           )) as tags
FROM requests r
         LEFT JOIN users u on r.user_id = u.id
         LEFT JOIN requests_tags rt on r.id = rt.request_id
         JOIN tags t on rt.tag_id = t.id
GROUP BY r.id, u.id;

create or replace function search_requests(page int, size int, category int, required_tags int[], price_from int,
                                           price_to int, sort_column text, sort_direction text, search text)
    returns setof requests_view as
$$
begin
    return query execute format('
        SELECT r.id,
               r.user_id,
               u.name as username,
               r.category_id,
               r.name,
               r.item_name,
               r.description,
               r.photos,
               r.price,
               r.quantity,
               r.date,
               r.finished,
               jsonb_agg(jsonb_build_object(
                       ''id'', t.id,
                       ''name'', t.name
                   )) as tags
        FROM requests r
                 LEFT JOIN users u on r.user_id = u.id
                 LEFT JOIN requests_tags rt on r.id = rt.request_id
                 JOIN tags t on rt.tag_id = t.id
        WHERE r.finished = false
          AND ($1 IS NULL OR r.category_id = $1)
          AND ($3 IS NULL OR r.price >= $3)
          AND ($4 IS NULL OR r.price <= $4)
          AND ($5 = '''' or lower(r.name) LIKE ''%%'' || $5 || ''%%''
                         or lower(r.item_name) LIKE ''%%'' || $5 || ''%%''
                         or lower(r.description) LIKE ''%%'' || $5 || ''%%'')
        GROUP BY r.id, u.id, date
        HAVING array_agg(t.id) @> $2
        ORDER BY %I %s LIMIT %s OFFSET %s', sort_column, sort_direction, size, size * page)
        USING category, required_tags, price_from, price_to, search;
end;
$$ LANGUAGE plpgsql;

create table if not exists proposals
(
    id           serial primary key,
    user_id      int    not null references users on delete cascade,
    request_id   int    not null references requests on delete cascade,
    description  text   not null,
    photos       text[] not null,
    price        int    not null,
    quantity     int    not null,
    date         date   not null                            default current_date,
    bought_by_id int    references users on delete set null default null
);

create or replace view proposals_view as
SELECT p.id,
       p.user_id,
       COALESCE(u.name, '') as username,
       p.request_id,
       r.name,
       r.item_name,
       p.description,
       p.photos,
       p.price,
       p.quantity,
       p.date,
       p.bought_by_id
FROM proposals p
         JOIN requests r on p.request_id = r.id
         LEFT JOIN users u on u.id = p.user_id;

INSERT INTO categories (id, name, photo)
VALUES (1, 'Недвижимость', '/images/property.png');
INSERT INTO categories (id, name, photo)
VALUES (2, 'Электроника', '/images/electronics.png');
INSERT INTO categories (id, name, photo)
VALUES (3, 'Хобби и отдых', '/images/hobby.png');
INSERT INTO categories (id, name, photo)
VALUES (4, 'Транспорт', '/images/transport.png');
INSERT INTO categories (id, name, photo)
VALUES (5, 'Одежда', '/images/clothes.png');
INSERT INTO categories (id, name, photo)
VALUES (6, 'Животные', '/images/pets.png');
INSERT INTO categories (id, name, photo)
VALUES (7, 'Для дома', '/images/house.png');
INSERT INTO categories (id, name, photo)
VALUES (8, 'Прочее', '/images/other.png');

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




