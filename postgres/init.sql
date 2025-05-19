CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE task_statuses (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    code VARCHAR(20) NOT NULL UNIQUE
);

CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deadline TIMESTAMP WITH TIME ZONE,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status_id INT NOT NULL REFERENCES task_statuses(id) ON DELETE RESTRICT
);

INSERT INTO task_statuses (name, code) VALUES
    ('New', 'new'),
    ('In Progress', 'in_progress'),
    ('Completed', 'completed')
ON CONFLICT (code) DO NOTHING;