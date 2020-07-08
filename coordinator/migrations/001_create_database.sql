BEGIN;

-- This block creates the hoist values function which is used to pull values
-- off of a data blob into a column for that value on the table.
CREATE EXTENSION IF NOT EXISTS hstore;
CREATE OR REPLACE FUNCTION hoist_values() RETURNS TRIGGER AS $$
  DECLARE
    value hstore;
    paths text[];
    path text;
    segments text[];
    segment text;
  BEGIN
    value = hstore(NEW);
    paths = TG_ARGV;

    FOREACH path IN ARRAY paths LOOP
      segments = string_to_array(path, '.')::text[];
      segment = segments[array_upper(segments, 1)];

      value := value || hstore(segment, NEW.data::jsonb#>>segments);
      NEW := NEW #= value;
    END LOOP;

    RETURN NEW;
  END;
$$ LANGUAGE plpgsql;

CREATE TABLE users (
  id char(29) primary key not null,
  data jsonb not null,
  CONSTRAINT users_id_check CHECK (data::jsonb#>>'{id}' = id)
);

CREATE UNIQUE INDEX users_email_idx ON users((data::jsonb#>>'{email}'));

CREATE TRIGGER users_hoist_tgr
  BEFORE INSERT ON users
  FOR EACH ROW EXECUTE PROCEDURE hoist_values('id');

CREATE TABLE roles (
  id char(29) primary key not null,
  data jsonb not null,
  CONSTRAINT roles_id_check CHECK (data::jsonb#>>'{id}' = id)
);

CREATE TRIGGER roles_hoist_tgr
  BEFORE INSERT ON roles
  FOR EACH ROW EXECUTE PROCEDURE hoist_values('id');

CREATE TABLE policies (
  id char(29) primary key not null,
  data jsonb not null,
  CONSTRAINT policies_id_check CHECK (data::jsonb#>>'{id}' = id)
);

CREATE TRIGGER policies_hoist_tgr
  BEFORE INSERT ON policies
  FOR EACH ROW EXECUTE PROCEDURE hoist_values('id');

CREATE TABLE attachments (
  id char(29) primary key not null,
  policy_id char(29) references policies(id) on delete cascade not null,
  role_id char(29) references roles(id) on delete cascade not null,
  data jsonb not null,
  CONSTRAINT attachments_id_check CHECK (data::jsonb#>>'{id}' = id),
  CONSTRAINT attachments_policy_id_check CHECK (data::jsonb#>>'{policy_id}' = policy_id),
  CONSTRAINT attachments_role_id_check CHECK (data::jsonb#>>'{role_id}' = role_id)
);

CREATE TRIGGER attachments_hoist_tgr
  BEFORE INSERT ON attachments
  FOR EACH ROW EXECUTE PROCEDURE hoist_values('id', 'policy_id', 'role_id');

CREATE TABLE assignments (
  id char(29) primary key not null,
  user_id char(29) references users(id) on delete cascade not null,
  role_id char(29) references roles(id) on delete cascade not null,
  data jsonb not null,
  UNIQUE(user_id, role_id),
  CONSTRAINT assignments_id_check CHECK (data::jsonb#>>'{id}' = id),
  CONSTRAINT assignments_user_id_check CHECK (data::jsonb#>>'{user_id}' = user_id),
  CONSTRAINT assignments_role_id_check CHECK (data::jsonb#>>'{role_id}' = role_Id)
);

CREATE TRIGGER assignments_hoist_tgr
  BEFORE INSERT ON assignments
  FOR EACH ROW EXECUTE PROCEDURE hoist_values('id', 'user_id', 'role_id');

CREATE TABLE tokens (
  id char(29) primary key not null,
  user_id char(29) references users(id) on delete cascade not null,
  data jsonb not null,
  CONSTRAINT tokens_id_check CHECK (data::jsonb#>>'{id}' = id),
  CONSTRAINT tokens_user_id_check CHECK (data::jsonb#>>'{user_id}' = user_id)
);

CREATE TRIGGER tokens_hoist_tgr
  BEFORE INSERT ON tokens
  FOR EACH ROW EXECUTE PROCEDURE hoist_values('id', 'user_id');

COMMIT;

---- create above / drop below ----

BEGIN;
DROP TABLE assignments;
DROP TABLE attachments;
DROP TABLE policies;
DROP TABLE roles;
DROP TABLE tokens;
DROP INDEX users_email_idx;
DROP TABLE users;
COMMIT;
