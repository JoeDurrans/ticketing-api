CREATE TABLE IF NOT EXISTS message (
    id UUID,
    ticket_id INT,
    author_id INT,
    content TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    PRIMARY KEY ((ticket_id), created_at, id)
) WITH CLUSTERING ORDER BY (created_at DESC, id ASC);