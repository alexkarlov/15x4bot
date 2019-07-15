CREATE TABLE IF NOT EXISTS reminder_chat (
    id SERIAL PRIMARY KEY,
    tg_chat_id INT NOT NULL,
    task_id INT NOT NULL
);