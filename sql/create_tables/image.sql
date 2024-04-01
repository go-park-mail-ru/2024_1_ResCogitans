CREATE TABLE image(
    id integer PRIMARY KEY,
    "path" text NOT NULL,
    sight_id integer REFERENCES sight(id)
);
