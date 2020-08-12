-- Write your migrate up statements here

BEGIN;

CREATE TABLE test (
	id char(29) not null primary key,
	data jsonb not null,
	CONSTRAINT test_id_equals CHECK (data::jsonb#>>'{id}' = id)
);

CREATE TRIGGER test_hoist_tgr
	BEFORE INSERT ON test
	FOR EACH ROW EXECUTE PROCEDURE hoist_values('id');

CREATE TABLE test_mutable (
	id char(29) not null primary key,
	data jsonb not null,
	CONSTRAINT test_mutable_id_equals CHECK (data::jsonb#>>'{id}' = id)
);

CREATE TRIGGER test_mutable_hoist_tgr
	BEFORE INSERT ON test_mutable
	FOR EACH ROW EXECUTE PROCEDURE hoist_values('id');

COMMIT;

---- create above / drop below ----

BEGIN;

DROP TRIGGER test_hoist_tgr on test;
DROP TABLE test;

DROP TRIGGER test_mutable_hoist_tgr on test_mutable;
DROP TABLE test_mutable;

COMMIT;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
