CREATE TABLE IF NOT EXISTS reminder_user (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    task_id INT NOT NULL
);