CREATE TABLE pages (
    id SERIAL PRIMARY KEY,
    url TEXT NOT NULL,
    url_id INT,
    title TEXT,
    summary TEXT,
    word_count INT,
    fetched_at TIMESTAMP
);
