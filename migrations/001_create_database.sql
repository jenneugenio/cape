BEGIN;

CREATE TABLE users (
  id char(29) primary key not null,
  data jsonb
);

CREATE TABLE services (
  id char(29) primary key not null,
  data jsonb
);

CREATE TABLE roles (
  id char(29) primary key not null,
  data jsonb
);

CREATE TABLE policies (
  id char(29) primary key not null,
  data jsonb
);

CREATE TABLE attachments (
  id char(29) primary key not null,
  policy_id char(29) references policies(id) on delete cascade not null,
  role_id char(29) references roles(id) on delete cascade not null,
  data jsonb
);

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
  UNIQUE (id, user_id, service_id),
  CHECK (user_id != null AND service_id != null)
);

CREATE TABLE assignments (
  id char(29) primary key not null,
  identity_id char(29) references identities(id) on delete cascade not null,
  role_id char(29) references roles(id) on delete cascade not null,
  data jsonb,
  UNIQUE(identity_id, role_id)
);

CREATE TABLE tokens (
  id char(29) primary key not null,
  identity_id char(29) references identities(id) on delete cascade,
  data jsonb
);

CREATE OR REPLACE FUNCTION create_identity()
  RETURNS trigger AS
$$
BEGIN
  IF TG_TABLE_NAME = 'users' THEN
    INSERT INTO identities(id, user_id)
        VALUES(NEW.id, NEW.id);
  ELSIF TG_TABLE_NAME = 'services' THEN
    INSERT INTO identities(id, service_id)
        VALUES(NEW.id, NEW.id);
  ELSE
    RAISE EXCEPTION 'Trigger function must be run on users or services table, not %', TG_TABLE_NAME;
  END IF;

   RETURN NEW;
END;
$$ language plpgsql;

CREATE TRIGGER create_identity
  AFTER INSERT
  ON users
  FOR EACH ROW
  EXECUTE PROCEDURE create_identity();

CREATE TRIGGER create_identity
  AFTER INSERT
  ON services
  FOR EACH ROW
  EXECUTE PROCEDURE create_identity();

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
