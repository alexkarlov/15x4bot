CREATE TABLE IF NOT EXISTS reminder_channel (
    id SERIAL PRIMARY KEY,
    tg_channel_id INT NOT NULL,
    task_id INT NOT NULL
);