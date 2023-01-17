CREATE SEQUENCE IF NOT EXISTS article_id_seq;

CREATE TABLE IF NOT EXISTS articles (
    id INT NOT NULL DEFAULT nextval('article_id_seq'),
    title TEXT NOT NULL,
    author TEXT,
    body TEXT NOT NULL,
    created TIMESTAMP NOT NULL
);
