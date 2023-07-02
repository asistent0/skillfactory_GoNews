DROP TABLE IF EXISTS post, author;

CREATE TABLE author (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE post (
    id SERIAL PRIMARY KEY,
    author_id INTEGER REFERENCES author(id) NOT NULL,
    title TEXT  NOT NULL,
    content TEXT NOT NULL,
    created_at BIGINT NOT NULL
);

INSERT INTO author (id, name) VALUES (0, 'Дмитрий');
INSERT INTO post (id, author_id, title, content, created_at) VALUES (0, 0, 'Статья', 'Содержание статьи', 0);