CREATE TABLE IF NOT EXISTS message
(
    id        BIGSERIAL PRIMARY KEY,
    title_id  BIGINT                  NOT NULL,
    text      VARCHAR(255)            NOT NULL,
    type      VARCHAR(255)            NOT NULL,
    timestamp TIMESTAMP DEFAULT NOW() NOT NULL,
    FOREIGN KEY (title_id) REFERENCES title (id) ON DELETE CASCADE ON UPDATE CASCADE
)