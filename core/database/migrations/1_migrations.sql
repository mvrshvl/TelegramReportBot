-- +migrate Up

CREATE TABLE IF NOT EXISTS messages
(
    id INT AUTO_INCREMENT PRIMARY KEY,
    city VARCHAR(20),
    `text` TEXT,
    media TEXT,
    place TEXT,
    `from` TEXT,
    timestamp TIMESTAMP
);

-- +migrate Down
DROP TABLE IF EXISTS messages;
