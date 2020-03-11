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

-- This block creates the create_identity trigger which is used to extract a
-- user or service id into the identities table
CREATE OR REPLACE FUNCTION create_identity() RETURNS TRIGGER AS $$
  BEGIN
    IF TG_TABLE_NAME = 'users' THEN
      INSERT INTO identities(id, user_id) VALUES(NEW.id, NEW.id);
    ELSIF TG_TABLE_NAME = 'services' THEN
      INSERT INTO identities(id, service_id) VALUES(NEW.id, NEW.id);
    ELSE
      RAISE EXCEPTION 'Trigger function must be run on users or services table, not %', TG_TABLE_NAME;
    END IF;
     RETURN NEW;
  END;
$$ LANGUAGE plpgsql;

CREATE TABLE users (
  id char(29) primary key not null,
  data jsonb not null,
  CONSTRAINT users_id_check CHECK (data::jsonb#>>'{id}' = id)
);

CREATE TRIGGER users_hoist_tgr
  BEFORE INSERT ON users
  FOR EACH ROW EXECUTE PROCEDURE hoist_values('id');

CREATE TRIGGER create_identity
  AFTER INSERT
  ON users
  FOR EACH ROW
  EXECUTE PROCEDURE create_identity();

CREATE TABLE services (
  id char(29) primary key not null,
  data jsonb not null,
  CONSTRAINT services_id_check CHECK (data::jsonb#>>'{id}' = id)
);

CREATE TRIGGER create_identity
  AFTER INSERT
  ON services
  FOR EACH ROW
  EXECUTE PROCEDURE create_identity();


CREATE TRIGGER services_hoist_tgr
  BEFORE INSERT ON services
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

-- identities table is required so that we don't have to have
-- user_id and service_id columns on every table that depends on
-- either a user or a service. When a user and service gets inserted
-- into their respective tables a trigger adds a corresponding row into
-- the identities table which is then used by all the tables which use
-- a user or service e.g. attachments and tokens.
CREATE TABLE identities (
  id char(29) primary key not null,
  user_id char(29) references users(id) on delete cascade,
  service_id char(29) references services(id) on delete cascade,
  CONSTRAINT identities_id_user_service_unique UNIQUE (id, user_id, service_id),
  CONSTRAINT identities_user_service_check CHECK (user_id != null AND service_id != null),
  CONSTRAINT identities_id_check CHECK (id = service_id OR id = user_id)
);

CREATE TABLE assignments (
  id char(29) primary key not null,
  identity_id char(29) references identities(id) on delete cascade not null,
  role_id char(29) references roles(id) on delete cascade not null,
  data jsonb not null,
  UNIQUE(identity_id, role_id),
  CONSTRAINT assignments_id_check CHECK (data::jsonb#>>'{id}' = id),
  CONSTRAINT assignments_identity_id_check CHECK (data::jsonb#>>'{identity_id}' = identity_id),
  CONSTRAINT assignments_role_id_check CHECK (data::jsonb#>>'{role_id}' = role_Id)
);

CREATE TRIGGER assignments_hoist_tgr
  BEFORE INSERT ON assignments
  FOR EACH ROW EXECUTE PROCEDURE hoist_values('id', 'identity_id', 'role_id');

CREATE TABLE tokens (
  id char(29) primary key not null,
  identity_id char(29) references identities(id) on delete cascade not null,
  data jsonb not null,
  CONSTRAINT tokens_id_check CHECK (data::jsonb#>>'{id}' = id),
  CONSTRAINT tokens_identity_id_check CHECK (data::jsonb#>>'{identity_id}' = identity_id)
);

CREATE TRIGGER tokens_hoist_tgr
  BEFORE INSERT ON tokens
  FOR EACH ROW EXECUTE PROCEDURE hoist_values('id', 'identity_id');

COMMIT;

---- create above / drop below ----

BEGIN;
DROP TRIGGER create_identity ON users;
DROP TRIGGER create_identity ON services;
DROP FUNCTION create_identity();
DROP TABLE assignments;
DROP TABLE attachments;
DROP TABLE policies;
DROP TABLE roles;
DROP TABLE tokens;
DROP TABLE identities;
DROP TABLE services;
DROP TABLE users;
COMMIT;
