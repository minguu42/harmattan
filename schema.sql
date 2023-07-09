DROP TABLE IF EXISTS tasks;
DROP TABLE IF EXISTS projects;
DROP TABLE IF EXISTS users;

CREATE TABLE users (
    id         CHAR(26)  PRIMARY KEY,
    name       CHAR(255) NOT NULL UNIQUE,
    api_key    CHAR(64)  NOT NULL UNIQUE,
    created_at DATETIME  NOT NULL,
    updated_at DATETIME  NOT NULL ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE projects (
    id         CHAR(26)  PRIMARY KEY,
    user_id    CHAR(26)  NOT NULL,
    name       CHAR(255) NOT NULL,
    created_at DATETIME  NOT NULL,
    updated_at DATETIME  NOT NULL ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE tasks (
    id           CHAR(26)  PRIMARY KEY,
    project_id   CHAR(26)  NOT NULL,
    title        CHAR(255) NOT NULL,
    completed_at DATETIME,
    created_at   DATETIME  NOT NULL,
    updated_at   DATETIME  NOT NULL ON UPDATE CURRENT_TIMESTAMP
);
