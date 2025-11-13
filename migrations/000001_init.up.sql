-- Команда
CREATE TABLE IF NOT EXISTS team (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

-- Пользователь
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    team_id INTEGER REFERENCES team(id) ON DELETE SET NULL
);

-- Статус Pull Request
CREATE TABLE IF NOT EXISTS pr_status (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

-- Pull Request
CREATE TABLE IF NOT EXISTS pull_request (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    author_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status_id INTEGER NOT NULL REFERENCES pr_status(id),
    merged_at TIMESTAMP NULL
);

-- Список назначенных ревьюверов
CREATE TABLE IF NOT EXISTS assigned_pr (
    pr_id INTEGER NOT NULL REFERENCES pull_request(id) ON DELETE CASCADE,
    reviewer_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (pr_id, reviewer_id),
    CHECK (pr_id IS NOT NULL AND reviewer_id IS NOT NULL)
);
