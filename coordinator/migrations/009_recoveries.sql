BEGIN;

CREATE TABLE recoveries (
  id char(29) primary key not null,
  user_id char(29) references users(id) on delete cascade not null,
  data jsonb not null,
  CONSTRAINT recoveries_id_check CHECK (data::jsonb#>>'{id}' = id),
  CONSTRAINT recoveries_user_id_check CHECK (data::jsonb#>>'{user_id}' = user_id)
);

CREATE TRIGGER recoveries_hoist_tgr
  BEFORE INSERT ON recoveries
  FOR EACH ROW EXECUTE PROCEDURE hoist_values('id', 'user_id');

COMMIT;

---- create above / drop below ----

BEGIN;
DROP TABLE recoveries;
COMMIT;
