BEGIN;

DROP TRIGGER roles_hoist_tgr ON roles;

ALTER TABLE roles DROP COLUMN data;

ALTER TABLE roles ADD COLUMN label text;
ALTER TABLE roles ADD COLUMN system bool;

CREATE UNIQUE INDEX roles_label_idx ON roles(label);

COMMIT;

---- create above / drop below ----

BEGIN;

CREATE TRIGGER roles_hoist_tgr
  BEFORE INSERT ON roles
  FOR EACH ROW EXECUTE PROCEDURE hoist_values('id');

ALTER TABLE roles DROP label;
ALTER TABLE roles DROP system;

ALTER TABLE roles ADD COLUMN data jsonb;

CREATE UNIQUE INDEX roles_label_idx ON roles((data::jsonb#>>'{label}'));

COMMIT;