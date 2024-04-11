CREATE TABLE message (
    id UUID PRIMARY KEY,
    ticket_id INT,
    author_id INT,
    content TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);