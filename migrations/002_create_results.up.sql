CREATE TABLE results (
    id SERIAL PRIMARY KEY,
    url_id INT REFERENCES urls(id),
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
