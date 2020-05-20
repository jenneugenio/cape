BEGIN;
CREATE TABLE source_schema
(
    id char(29) primary key not null,
    -- source_id char(29) references source(id) on delete cascade not null,
    data jsonb not null,

    CONSTRAINT source_schema_id_check CHECK(data::jsonb#>>'{id}' = id)
);

CREATE TRIGGER source_schema_hoist_tgr
  BEFORE INSERT ON source_schema
  FOR EACH ROW EXECUTE PROCEDURE hoist_values('id');

COMMIT;

---- create above / drop below ----

BEGIN;
DROP TRIGGER source_schema_hoist_tgr ON source_schema;
DROP TABLE source_schema;
COMMIT;