CREATE TABLE journey(
    id integer PRIMARY KEY,
    client_id integer REFERENCES "user"(id),
    description text
);
