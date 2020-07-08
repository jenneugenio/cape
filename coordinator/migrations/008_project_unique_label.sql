BEGIN;

CREATE UNIQUE INDEX projects_label_idx ON projects((data::jsonb#>>'{label}'));

COMMIT;

---- create above / drop below ----

BEGIN;

DROP INDEX projects_label_idx;

COMMIT;