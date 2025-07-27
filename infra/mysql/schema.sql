create table users (
    id              char(26)     not null primary key,
    email           varchar(254) not null,
    hashed_password char(60)     not null,
    created_at      datetime     not null default current_timestamp,
    updated_at      datetime     not null default current_timestamp on update current_timestamp,
    unique (email)
);

create table projects (
    id          char(26)     not null primary key,
    user_id     char(26)     not null,
    name        varchar(80)  not null,
    color       varchar(255) not null,
    is_archived tinyint(1)   not null default 0,
    created_at  datetime     not null default current_timestamp,
    updated_at  datetime     not null default current_timestamp on update current_timestamp,
    foreign key (user_id) references users (id) on delete cascade,
    check (color in ('blue', 'brown', 'default', 'gray', 'green', 'orange', 'pink', 'purple', 'red',
                     'yellow'))
);

create table tasks (
    id           char(26)         not null primary key,
    user_id      char(26)         not null,
    project_id   char(26)         not null,
    name         varchar(100)     not null,
    content      varchar(300)     not null,
    priority     tinyint unsigned not null,
    due_on       date,
    completed_at datetime,
    created_at   datetime         not null default current_timestamp,
    updated_at   datetime         not null default current_timestamp on update current_timestamp,
    foreign key (user_id) references users (id) on delete cascade,
    foreign key (project_id) references projects (id) on delete cascade,
    check (priority between 0 and 3)
);

create table steps (
    id           char(26)     not null primary key,
    user_id      char(26)     not null,
    task_id      char(26)     not null,
    name         varchar(100) not null,
    completed_at datetime,
    created_at   datetime     not null default current_timestamp,
    updated_at   datetime     not null default current_timestamp on update current_timestamp,
    foreign key (user_id) references users (id) on delete cascade,
    foreign key (task_id) references tasks (id) on delete cascade
);

create table tags (
    id         char(26)    not null primary key,
    user_id    char(26)    not null,
    name       varchar(20) not null,
    created_at datetime    not null default current_timestamp,
    updated_at datetime    not null default current_timestamp on update current_timestamp,
    foreign key (user_id) references users (id) on delete cascade
);

create table tasks_tags (
    task_id    char(26) not null references tasks (id) on delete cascade,
    tag_id     char(26) not null references tags (id) on delete cascade,
    created_at datetime not null default current_timestamp,
    primary key (task_id, tag_id),
    foreign key (task_id) references tasks (id) on delete cascade,
    foreign key (tag_id) references tags (id) on delete cascade
);
