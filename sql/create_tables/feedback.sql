CREATE TABLE feedback(
    id integer PRIMARY KEY,
    user_id integer REFERENCES "user"(id),
    sight_id integer REFERENCES sight(id),
    rating integer NOT NULL,
    feedback text
);
