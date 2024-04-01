CREATE TABLE sight(
    id integer PRIMARY KEY,
    rating float NOT NULL,
    name text NOT NULL,
    description text,
    city_id integer REFERENCES city (id),
    country_id integer REFERENCES country (id)
);

	-- {
	-- 	ID:          12,
	-- 	Rating:      4.5,
	-- 	Name:        "Салтинский Подземный Водопад",
	-- 	Description: "Единственный в России подземный водопад.",
	-- 	City:        "Салта",
	-- 	Url:         "public/12.jpg",
	-- },
