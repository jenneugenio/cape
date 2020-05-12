BEGIN;
DROP INDEX sessions_token_idx;
COMMIT;

---- create above / drop below ----

BEGIN;
CREATE UNIQUE INDEX sessions_token_idx ON sessions((data::jsonb#>>'{token}'));
COMMIT;