ALTER TABLE lections RENAME TO lectures;
ALTER TABLE event_lections RENAME TO event_lectures;
ALTER TABLE event_lectures RENAME COLUMN id_lection TO id_lecture;