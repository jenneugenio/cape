BEGIN;

CREATE TABLE contributors (
    id char(29) primary key not null,
    user_id char(29) references users(id) on delete cascade not null,
    project_id char(29) references projects(id) on delete cascade not null,

    data jsonb not null,
    constraint user_project unique (user_id, project_id)
);

create trigger contributors_hoist_tgr
    before insert on contributors
    for each row execute procedure hoist_values('id', 'user_id', 'project_id');

COMMIT;

---- create above / drop below ----

BEGIN;

DROP TRIGGER contributors_hoist_tgr on contributors;
DROP TABLE contributors;

COMMIT;