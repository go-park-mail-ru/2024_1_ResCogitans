-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

CREATE TABLE city(
    id integer PRIMARY KEY,
    name text NOT NULL
);

CREATE TABLE country(
    id integer PRIMARY KEY,
    name text NOT NULL
);

CREATE TABLE "user" (
    id integer PRIMARY KEY,
    email text NOT NULL,
    passwrd text NOT NULL
);

CREATE TABLE sight(
    id integer PRIMARY KEY,
    rating float NOT NULL,
    name text NOT NULL,
    description text,
    city_id integer REFERENCES city (id),
    country_id integer REFERENCES country (id)
);

CREATE TABLE image(
    id integer PRIMARY KEY,
    "path" text NOT NULL,
    sight_id integer REFERENCES sight(id)
);

CREATE TABLE journey(
    id integer PRIMARY KEY,
    client_id integer REFERENCES "user"(id),
    description text
);

CREATE TABLE journey_sight(
    id integer PRIMARY KEY,
    journey_id integer REFERENCES journey(id),
    sight_id integer REFERENCES sight(id),
    priority integer NOT NULL,
    note text
);

CREATE TABLE feedback(
    id integer PRIMARY KEY,
    user_id integer REFERENCES "user"(id),
    sight_id integer REFERENCES sight(id),
    rating integer NOT NULL,
    feedback text
);


INSERT INTO city (id, name) VALUES 
(1, 'Москва'),
(2, 'Вольск'),
(3, 'Тамбов'),
(4, 'Бахчисарай'),
(5, 'Евпатория'),
(6, 'Балаклава'),
(7, 'Казань'),
(8, 'Салта'),
(9, 'Мир'),
(10, 'Гудаута');

INSERT INTO country(id, name) VALUES
(1, 'Россия'),
(2, 'Беларусь'),
(3, 'Татарстан'),
(4, 'Крым'),
(5, 'Дагестан'),
(6, 'Абхазия');


INSERT INTO sight(id, rating, name, description, city_id, country_id) VALUES 
    (
        1, 
        2.1, 
        'У дяди Вани',	
        'Ресторан с видом на Сталинскую высотку.', 
        1,
        1
    ),
	(
		2,
		3.1,
		'Государственный музей изобразительных искусств имени А.С. Пушкина',
		'Музей.',
		1,
		1
	),
	(
		3,
		4.99,
		'МГТУ им. Н. Э. Баумана',
		'Хороший вуз.',
		1,
		1
	),
	(
		4,
		3.2,
		'Вкусно - и точка',
		'Неплохое кафе, вызывает гастрит.',
		1,
		1
	),
	(
		5,
		4.1,
		'Краеведческий музей',
		'Один из самых больших провинциальных музеев краеведческого профиля.',
		2,
		1
	),
	(
		6,
		4.3,
		'Спасо-Преображенский кафедральный собор',
		'Спасо-Преображенский кафедральный собор расположен в центре города и является первым каменным храмом Тамбова и старейшим в Тамбовской обл.',
		3,
		1
	),
	(
		7,
		3.9,
		'Мирский замок',
		'Памятник архитектуры, внесён в список Всемирного наследия ЮНЕСКО.',
		9,
		2
	),
	(
		8,
		4.9,
		'Чуфут-Кале',
		'Пещерный город в Крыму. Топ.',
		4,
		4
	),
	(
		9,
		3.5,
		'Сасык-Сиваш',
		'Розовое озеро. Оно реально розовое.',
		5,
		4
	),
	(
		10,
		4.7,
		'Крепость Чембело',
		'Остатки крепости.',
		6,
        4
	),
	(
		11,
		4.0,
		'Мечеть Кул Шариф',
		'Главная джума-мечеть республики Татарстан и города Казани.',
		7,
		3
	),
	(
		12,
		4.5,
		'Салтинский Подземный Водопад',
		'Единственный в России подземный водопад.',
		8,
        5
	);


INSERT INTO image(id, path, sight_id) VALUES 
(1, 'public/1.jpg', 1),
(2, 'public/2.jpg', 2),
(3, 'public/3.jpg', 3),
(4, 'public/4.jpg', 4),
(5, 'public/5.jpg', 5),
(6, 'public/6.jpg', 6),
(7, 'public/7.jpg', 7),
(8, 'public/8.jpg', 8),
(9, 'public/9.jpg', 9),
(10, 'public/10.jpg', 10),
(11, 'public/11.jpg', 11),
(12, 'public/12.jpg', 12);


INSERT INTO "user"(id, email, passwrd) VALUES
(1, 'djafarovemil04@mail.ru', 246858);

INSERT INTO feedback(id, user_id, sight_id, rating, feedback) VALUES 
(1, 1, 1, 4, 'Все понравилось');


ALTER TABLE city RENAME COLUMN name TO city;
ALTER TABLE country RENAME COLUMN name TO country;


-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
