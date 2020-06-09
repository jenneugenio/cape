BEGIN;

CREATE TABLE test (
    test_text text
);

INSERT INTO test (test_text) VALUES ('test');

COMMIT;


---- create above / drop below ----

BEGIN;
DROP TABLE test;
COMMIT;