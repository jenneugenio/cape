BEGIN;

CREATE UNIQUE INDEX services_email_idx ON services((data::jsonb#>>'{email}'));
CREATE UNIQUE INDEX users_email_idx ON users((data::jsonb#>>'{email}'));

COMMIT;

---- create above / drop below ----

BEGIN;

DROP INDEX services_email_idx;
DROP INDEX users_email_idx;

COMMIT;
