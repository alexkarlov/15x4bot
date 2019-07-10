CREATE TYPE placetype AS ENUM ('for_event','for_repetition', 'for_all');
ALTER TABLE places ADD COLUMN IF NOT EXISTS type placetype DEFAULT 'for_repetition';
UPDATE places SET type='for_all' WHERE id = 1;
UPDATE places SET type='for_event' WHERE id = 2;
