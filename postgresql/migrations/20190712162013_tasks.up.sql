CREATE TYPE task_type AS ENUM ('reminder_lector', 'reminder_designer', 'reminder_grammar_nazi', 'reminder_fb_event', 'post_tg_chat', 'post_tg_channel');
CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    type task_type NOT NULL,
    execution_time  TIMESTAMP NOT NULL,
    status SMALLINT, 
    details JSON NOT NULL,
    cdate TIMESTAMP NOT NULL DEFAULT NOW(),
    udate TIMESTAMP NOT NULL DEFAULT NOW()
);
