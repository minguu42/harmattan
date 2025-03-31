CREATE TABLE IF NOT EXISTS users (
    id                CHAR(26)     PRIMARY KEY,
    email             VARCHAR(254) NOT NULL,
    hashed_password   CHAR(60)     NOT NULL,
    created_at        DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at        DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE INDEX (email)
) ENGINE = InnoDB CHARSET = utf8mb4 COLLATE = utf8mb4_bin;
