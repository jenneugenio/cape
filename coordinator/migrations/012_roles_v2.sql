BEGIN;

alter table assignments add column project_id char(29);
alter table assignments drop constraint assignments_user_id_role_id_key;
alter table assignments add constraint unique_assignment unique (user_id, role_id, project_id);

COMMIT;

---- create above / drop below ----

BEGIN;

alter table assignments drop constraint unique_assignment;
alter table assignments add constraint assignments_user_id_role_id_key unique(user_id, role_id);
alter table assignments drop column project_id;

COMMIT;