SET CHARACTER SET utf8mb4;

INSERT INTO users (id, name, api_key, created_at, updated_at)
VALUES ('01DXF6DT000000000000000000', 'minguu42', 'rAM9Fm9huuWEKLdCwHBcju9Ty_-TL2tDsAicmMrXmUnaCGp3RtywzYpMDPdEtYtR', '2020-01-01 00:00:00', '2020-01-01 00:00:00');

INSERT INTO projects (id, user_id, name, color, is_archived, created_at, updated_at)
VALUES ('01DXF6DT000000000000000000', '01DXF6DT000000000000000000', 'プロジェクト1', '#1A2B3C', false, '2020-01-01 00:00:00', '2020-01-01 00:00:00'),
       ('01DXF6DT000000000000000001', '01DXF6DT000000000000000000', 'プロジェクト2', '#1A2B3C', false, '2020-01-01 00:00:01', '2020-01-01 00:00:01');

INSERT INTO tasks (id, project_id, title, content, priority, due_on, completed_at, created_at, updated_at)
VALUES ('01DXF6DT000000000000000000', '01DXF6DT000000000000000000', 'タスク1', 'Hello, 世界!', 0, NULL, NULL, '2020-01-01 00:00:00', '2020-01-01 00:00:00'),
       ('01DXF6DT000000000000000001', '01DXF6DT000000000000000000', 'タスク2', '', 3, '2020-01-02', '2020-01-02 00:00:00', '2020-01-01 00:00:01', '2020-01-02 00:00:00');
