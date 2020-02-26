BEGIN;

CREATE TABLE users (
  id char(29) primary key not null,
  data jsonb
);

CREATE TABLE services (
  id char(29) primary key not null,
  data jsonb
);

CREATE TABLE tokens (
  id char(29) primary key not null,
  user_id char(29) references users(id) on delete cascade,
  service_id char(29) references services(id) on delete cascade,
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

COMMIT;
