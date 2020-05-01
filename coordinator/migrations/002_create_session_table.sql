BEGIN;

CREATE TABLE sessions (
  id char(29) primary key not null,
  identity_id char(29) references identities(id) on delete cascade not null,
  data jsonb not null,
  CONSTRAINT sessions_id_check CHECK (data::jsonb#>>'{id}' = id),
  CONSTRAINT sessions_identity_id_check CHECK (data::jsonb#>>'{identity_id}' = identity_id)
);

CREATE UNIQUE INDEX sessions_token_idx ON sessions((data::jsonb#>>'{token}'));

CREATE TRIGGER sessions_hoist_tgr
  BEFORE INSERT ON sessions
  FOR EACH ROW EXECUTE PROCEDURE hoist_values('id', 'identity_id');

COMMIT;

---- create above / drop below ----

BEGIN;
DROP TABLE sessions;
COMMIT;
