create table if not exists users (
    id              char(26) primary key,
    email           varchar(254) not null unique,
    hashed_password char(60)     not null,
    created_at      datetime     not null default current_timestamp,
    updated_at      datetime     not null default current_timestamp on update current_timestamp
);

create table if not exists projects (
    id          char(26) primary key,
    user_id     char(26)     not null references users (id) on delete cascade,
    name        varchar(80)  not null,
    color       varchar(255) not null check (color in ('blue', 'brown', 'default', 'gray', 'green', 'orange', 'pink', 'purple', 'red', 'yellow')),
    is_archived tinyint(1)   not null default 0,
    created_at  datetime     not null default current_timestamp,
    updated_at  datetime     not null default current_timestamp on update current_timestamp
);

create table tasks (
    id           char(26) primary key,
    user_id      char(26)            not null references users (id) on delete cascade,
    project_id   char(26)            not null references projects (id) on delete cascade,
    name         varchar(100)        not null,
    content      varchar(300)        not null,
    priority     tinyint(2) unsigned not null check (priority between 0 and 3),
    due_on       date,
    completed_at datetime,
    created_at   datetime            not null default current_timestamp,
    updated_at   datetime            not null default current_timestamp on update current_timestamp
);

create table steps (
    id           char(26) primary key,
    user_id      char(26)     not null references users (id) on delete cascade,
    task_id      char(26)     not null references tasks (id) on delete cascade,
    name         varchar(100) not null,
    completed_at datetime,
    created_at   datetime    not null default current_timestamp,
    updated_at   datetime    not null default current_timestamp on update current_timestamp
);

create table tags (
    id         char(26) primary key,
    user_id    char(26)    not null references users (id) on delete cascade,
    name       varchar(20) not null,
    created_at datetime    not null default current_timestamp,
    updated_at datetime    not null default current_timestamp on update current_timestamp
);

create table tasks_tags (
    task_id    char(26) not null references tasks (id) on delete cascade,
    tag_id     char(26) not null references tags (id) on delete cascade,
    created_at datetime not null default current_timestamp,
    primary key (task_id, tag_id)
);
