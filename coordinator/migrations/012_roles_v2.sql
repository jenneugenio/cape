BEGIN;

alter table assignments add column project_id char(29);
alter table assignments drop constraint assignments_user_id_role_id_key;

-- TODO when project_id is null, this index doesn't work properly (doesn't get enforced for global roles)
alter table assignments add constraint unique_assignment unique (user_id, project_id);

DROP TRIGGER assignments_hoist_tgr on assignments;
CREATE TRIGGER assignments_hoist_tgr
    BEFORE INSERT ON assignments
    FOR EACH ROW EXECUTE PROCEDURE hoist_values('id', 'user_id', 'role_id', 'project_id');

COMMIT;

---- create above / drop below ----

BEGIN;

DROP TRIGGER assignments_hoist_tgr on assignments;
CREATE TRIGGER assignments_hoist_tgr
    BEFORE INSERT ON assignments
    FOR EACH ROW EXECUTE PROCEDURE hoist_values('id', 'user_id', 'role_id');

alter table assignments drop constraint unique_assignment;
alter table assignments add constraint assignments_user_id_role_id_key unique(user_id, role_id);
alter table assignments drop column project_id;

COMMIT;