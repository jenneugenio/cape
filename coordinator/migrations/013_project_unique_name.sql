BEGIN;

CREATE UNIQUE INDEX projects_name_idx ON projects((data::jsonb#>>'{name}'));

COMMIT;

---- create above / drop below ----

BEGIN;

DROP INDEX projects_name_idx;

COMMIT;