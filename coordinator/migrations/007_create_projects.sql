BEGIN;

CREATE TABLE projects (
  id char(29) primary key not null,
  data jsonb not null,
  CONSTRAINT projects_id_check CHECK (data::jsonb#>>'{id}' = id)
);

CREATE TRIGGER projects_hoist_tgr
  BEFORE INSERT ON projects
  FOR EACH ROW EXECUTE PROCEDURE hoist_values('id', 'current_spec');

CREATE TRIGGER projects_hoist_tgr
    BEFORE UPDATE ON projects
    FOR EACH ROW EXECUTE PROCEDURE hoist_values('id', 'current_spec');

CREATE TABLE project_specs (
  id char(29) primary key not null,
  project_id char(29) references projects(id),
  parent_id char(29) references project_specs(id),
  data jsonb not null,
  CONSTRAINT project_spec_id_check CHECK (data::jsonb#>>'{id}' = id),
  CONSTRAINT project_spec_parent_id_check CHECK (data::jsonb#>>'{parent_id}' = parent_id)
);

CREATE TRIGGER project_specs_hoist_tgr
  BEFORE INSERT ON project_specs
  FOR EACH ROW EXECUTE PROCEDURE hoist_values('id', 'parent_id', 'project_id');

ALTER TABLE projects add column current_spec char(29) references project_specs(id);

COMMIT;

---- create above / drop below ----

BEGIN;
ALTER TABLE projects drop column current_spec;
DROP TRIGGER project_specs_hoist_tgr ON project_specs;
DROP TABLE project_specs;
DROP TRIGGER projects_hoist_tgr ON projects;
DROP TABLE projects;
COMMIT;