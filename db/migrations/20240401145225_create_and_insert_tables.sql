-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

CREATE TABLE city(
    id integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY ,
    city text NOT NULL
);

CREATE TABLE country(
    id integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY ,
    country text NOT NULL UNIQUE
);

CREATE TABLE user_data (
    id integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY ,
    email text NOT NULL UNIQUE,
    passwrd text NOT NULL
);

CREATE TABLE profile (
    user_id integer REFERENCES user_data(id),
    username text UNIQUE,
    avatar text,
	bio text
);

CREATE TABLE sight(
    id integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY ,
    rating float NOT NULL CHECK (rating > 0 AND rating <= 5),
    name text NOT NULL,
    description text,
    city_id integer REFERENCES city (id),
    country_id integer REFERENCES country (id),
	UNIQUE (name, city_id)
);

CREATE TABLE image_data(
    id integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY ,
    "path" text NOT NULL UNIQUE,
    sight_id integer REFERENCES sight(id)
);

CREATE TABLE journey(
    id integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY ,
	name text NOT NULL UNIQUE,
    user_id integer REFERENCES user_data(id),
    description text
);

CREATE TABLE journey_sight(
    id integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY ,
    journey_id integer REFERENCES journey(id),
    sight_id integer REFERENCES sight(id),
    priority integer NOT NULL,
	UNIQUE (journey_id, sight_id)
);

CREATE TABLE feedback(
    id integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY ,
    user_id integer REFERENCES user_data(id),
    sight_id integer REFERENCES sight(id),
    rating integer NOT NULL CHECK (rating > 0 AND rating <= 5),
    feedback text NOT NULL
);


INSERT INTO city (city) VALUES 
('Москва'),
('Вольск'),
('Тамбов'),
('Бахчисарай'),
('Евпатория'),
('Балаклава'),
('Казань'),
('Салта'),
('Мир'),
('Гудаута');

INSERT INTO country(country) VALUES
('Россия'),
('Беларусь'),
('Татарстан'),
('Крым'),
('Дагестан'),
('Абхазия');


INSERT INTO sight(rating, name, description, city_id, country_id) VALUES 
    (
        2.1, 
        'У дяди Вани',	
        'Ресторан с видом на Сталинскую высотку.', 
        1,
        1
    ),
	(
		3.1,
		'Государственный музей изобразительных искусств имени А.С. Пушкина',
		'Музей.',
		1,
		1
	),
	(
		4.99,
		'МГТУ им. Н. Э. Баумана',
		'Хороший вуз.',
		1,
		1
	),
	(
		3.2,
		'Вкусно - и точка',
		'Неплохое кафе, вызывает гастрит.',
		1,
		1
	),
	(
		4.1,
		'Краеведческий музей',
		'Один из самых больших провинциальных музеев краеведческого профиля.',
		2,
		1
	),
	(
		4.3,
		'Спасо-Преображенский кафедральный собор',
		'Спасо-Преображенский кафедральный собор расположен в центре города и является первым каменным храмом Тамбова и старейшим в Тамбовской обл.',
		3,
		1
	),
	(
		3.9,
		'Мирский замок',
		'Памятник архитектуры, внесён в список Всемирного наследия ЮНЕСКО.',
		9,
		2
	),
	(
		4.9,
		'Чуфут-Кале',
		'Пещерный город в Крыму. Топ.',
		4,
		4
	),
	(
		3.5,
		'Сасык-Сиваш',
		'Розовое озеро. Оно реально розовое.',
		5,
		4
	),
	(
		4.7,
		'Крепость Чембело',
		'Остатки крепости.',
		6,
        4
	),
	(
		4.0,
		'Мечеть Кул Шариф',
		'Главная джума-мечеть республики Татарстан и города Казани.',
		7,
		3
	),
	(
		4.5,
		'Салтинский Подземный Водопад',
		'Единственный в России подземный водопад.',
		8,
        5
	);


INSERT INTO image(path, sight_id) VALUES 
('public/1.jpg', 1),
('public/2.jpg', 2),
('public/3.jpg', 3),
('public/4.jpg', 4),
('public/5.jpg', 5),
('public/6.jpg', 6),
('public/7.jpg', 7),
('public/8.jpg', 8),
('public/9.jpg', 9),
('public/10.jpg', 10),
('public/11.jpg', 11),
('public/12.jpg', 12);


CREATE OR REPLACE FUNCTION create_profile()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO profile (user_id, username, bio, avatar)
    VALUES (NEW.id, NEW.email, '', '');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER create_profile_trigger
AFTER INSERT ON user_data
FOR EACH ROW
EXECUTE FUNCTION create_profile();


-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

DROP TABLE IF EXISTS user_data CASCADE;
DROP TABLE IF EXISTS city CASCADE;
DROP TABLE IF EXISTS country CASCADE;
DROP TABLE IF EXISTS sight CASCADE;
DROP TABLE IF EXISTS journey CASCADE;
DROP TABLE IF EXISTS journey_sight CASCADE;
DROP TABLE IF EXISTS image_data CASCADE;
DROP TABLE IF EXISTS feedback CASCADE;
DROP TABLE IF EXISTS profile CASCADE;
