BEGIN;

CREATE TABLE sources (
  id char(29) primary key not null,
  data jsonb not null,
  CONSTRAINT sources_id_check CHECK (data::jsonb#>>'{id}' = id)
);

CREATE TRIGGER sources_hoist_tgr
  BEFORE INSERT ON sessions
  FOR EACH ROW EXECUTE PROCEDURE hoist_values('id');

COMMIT;

---- create above / drop below ----

BEGIN;
DROP TABLE sources;
COMMIT;