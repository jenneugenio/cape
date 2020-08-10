BEGIN;

create table secrets(
    name text primary key not null,
    data jsonb not null,
    constraint secret_id_check check (data::jsonb#>>'{name}' = name)
);

CREATE TRIGGER secrets_hoist_tgr
    BEFORE INSERT ON secrets
    FOR EACH ROW EXECUTE PROCEDURE hoist_values('name');

COMMIT;

---- create above / drop below ----

BEGIN;

drop trigger secrets_hoist_tgr on secrets;
drop table secrets;


COMMIT;