BEGIN;

CREATE TABLE sources (
  id char(29) primary key not null,
  service_id char(29) references services(id),
  data jsonb not null,
  CONSTRAINT sources_id_check CHECK (data::jsonb#>>'{id}' = id),
  CONSTRAINT sources_service_id_check CHECK (data::jsonb#>>'{service_id}' = service_id)
);

CREATE UNIQUE INDEX source_label_idx ON sources((data::jsonb#>>'{label}'));

CREATE TRIGGER sources_hoist_tgr
  BEFORE INSERT ON sources
  FOR EACH ROW EXECUTE PROCEDURE hoist_values('id', 'service_id');

COMMIT;

---- create above / drop below ----

BEGIN;
DROP TABLE sources;
COMMIT;
