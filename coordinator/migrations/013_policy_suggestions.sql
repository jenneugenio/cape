BEGIN;

create table suggestions(
    id text primary key not null,
    project_id char(29) references projects(id) on delete cascade not null,
    project_spec_id char(29) references project_specs(id) on delete cascade not null,
    data jsonb not null,
    constraint suggestion_id_check check (data::jsonb#>>'{id}' = id),
    constraint suggestion_project_id_check check (data::jsonb#>>'{project_id}' = project_id),
    constraint suggestion_project_spec_id_check check (data::jsonb#>>'{project_spec_id}' = project_spec_id)
);

CREATE TRIGGER suggestions_hoist_tgr
    BEFORE INSERT ON suggestions
    FOR EACH ROW EXECUTE PROCEDURE hoist_values('id', 'project_id', 'project_spec_id');

COMMIT;

---- create above / drop below ----

BEGIN;

drop trigger suggestions_hoist_tgr on suggestions;
drop table suggestions;


COMMIT;