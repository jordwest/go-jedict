
CREATE TABLE jedict_entries (
    id SERIAL PRIMARY KEY,
    sequence_id integer
);

CREATE TABLE jedict_kanji (
    id SERIAL PRIMARY KEY,
    entry_id integer,
    kanji varchar(255)
);

CREATE TABLE jedict_meaning (
    id SERIAL PRIMARY KEY,
    entry_id integer,
    meaning varchar(1024)
);

CREATE TABLE jedict_reading (
    id SERIAL PRIMARY KEY,
    entry_id integer,
    reading varchar(1024)
);

