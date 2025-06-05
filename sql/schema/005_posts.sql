-- +goose Up
CREATE TABLE posts(
     id int PRIMARY KEY,
     created_at TIMESTAMP NOT NULL,
     updated_at TIMESTAMP NOT NULL,
     title TEXT NOT NULL,
     url TEXT UNIQUE NOT NULL,
     description TEXT,
     published_at TIMESTAMP,
     feed_id int NOT NULL 
     REFERENCES feeds
     ON DELETE CASCADE
);
-- +goose Down
DROP TABLE posts;