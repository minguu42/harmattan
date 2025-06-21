create table if not exists users (
    id              char(26) primary key,
    email           varchar(254) not null,
    hashed_password char(60)     not null,
    created_at      datetime     not null default current_timestamp,
    updated_at      datetime     not null default current_timestamp on update current_timestamp,
    unique index (email)
);

create table if not exists projects (
    id          char(26) primary key,
    user_id     char(26)    not null references users(id) on delete cascade,
    name        varchar(26) not null,
    color       char(7)     not null,
    is_archived tinyint(1)  not null default 0,
    created_at  datetime    not null default current_timestamp,
    updated_at  datetime    not null default current_timestamp on update current_timestamp
)
