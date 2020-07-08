BEGIN;

CREATE UNIQUE INDEX policies_label_idx ON policies((data::jsonb#>>'{label}'));
ALTER TABLE attachments ADD CONSTRAINT attachment_unique UNIQUE (policy_id, role_id);

COMMIT;

---- create above / drop below ----

BEGIN;
ALTER TABLE attachments DROP CONSTRAINT attachment_unique;
DROP INDEX policies_label_idx;
COMMIT;
