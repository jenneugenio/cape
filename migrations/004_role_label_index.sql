BEGIN;

CREATE UNIQUE INDEX roles_label_idx ON roles((data::jsonb#>>'{label}'));

COMMIT;

---- create above / drop below ----

BEGIN;
DROP INDEX roles_label_idx;
COMMIT;