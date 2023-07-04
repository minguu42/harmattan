SET CHARACTER SET utf8mb4;

INSERT INTO users (id, name, api_key, created_at, updated_at)
VALUES (1, 'minguu42', 'rAM9Fm9huuWEKLdCwHBcju9Ty_-TL2tDsAicmMrXmUnaCGp3RtywzYpMDPdEtYtR', '2020-01-01 00:00:00', '2020-01-01 00:00:00');

INSERT INTO projects (id, user_id, name, created_at, updated_at)
VALUES (1, 1, 'プロジェクト1', '2020-01-01 00:00:00', '2020-01-01 00:00:00'),
       (2, 1, 'プロジェクト2', '2020-01-01 00:00:00', '2020-01-01 00:00:00');

INSERT INTO tasks (id, project_id, title, completed_at, created_at, updated_at)
VALUES (1, 1, 'タスク1', NULL, '2020-01-01 00:00:00', '2020-01-01 00:00:00'),
       (2, 1, 'タスク2', '2020-01-02 00:00:00', '2020-01-01 00:00:00', '2020-01-02 00:00:00');
