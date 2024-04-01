CREATE TABLE journey_sight(
    id integer PRIMARY KEY,
    journey_id integer REFERENCES journey(id),
    sight_id integer REFERENCES sight(id),
    priority integer NOT NULL,
    note text
);
