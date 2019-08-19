ALTER TABLE lectures RENAME TO lections;
ALTER TABLE event_lectures RENAME TO event_lections;
ALTER TABLE event_lections RENAME COLUMN id_lecture TO id_lection;