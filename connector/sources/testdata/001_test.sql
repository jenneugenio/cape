BEGIN;

CREATE TABLE test (
    int2 smallint,
    int4 integer,
    int8 bigint,
    float8 double precision,
    float4 real,
    vchar varchar(20),
    ch char(20),
    txt text,
    ts timestamp,
    bool boolean,
    bytes bytea
);

INSERT INTO test VALUES (2, 4, 8, 8.8, 4.4, 'hello', 'thisisatest', 'andthis', NOW(), false, decode('DEADBEEF', 'hex'));
INSERT INTO test VALUES (2, 4, 8, 8.8, 4.4, 'hello', 'thisisatest', 'andthis', NOW(), false, decode('DEADBEEF', 'hex'));
INSERT INTO test VALUES (2, 4, 8, 8.8, 4.4, 'hello', 'thisisatest', 'andthis', NOW(), false, decode('DEADBEEF', 'hex'));
INSERT INTO test VALUES (2, 4, 8, 8.8, 4.4, 'hello', 'thisisatest', 'andthis', NOW(), false, decode('DEADBEEF', 'hex'));
INSERT INTO test VALUES (2, 4, 8, 8.8, 4.4, 'hello', 'thisisatest', 'andthis', NOW(), false, decode('DEADBEEF', 'hex'));
INSERT INTO test VALUES (2, 4, 8, 8.8, 4.4, 'hello', 'thisisatest', 'andthis', NOW(), false, decode('DEADBEEF', 'hex'));
INSERT INTO test VALUES (2, 4, 8, 8.8, 4.4, 'hello', 'thisisatest', 'andthis', NOW(), false, decode('DEADBEEF', 'hex'));
INSERT INTO test VALUES (2, 4, 8, 8.8, 4.4, 'hello', 'thisisatest', 'andthis', NOW(), false, decode('DEADBEEF', 'hex'));
INSERT INTO test VALUES (2, 4, 8, 8.8, 4.4, 'hello', 'thisisatest', 'andthis', NOW(), false, decode('DEADBEEF', 'hex'));

COMMIT;

---- create above / drop below ----

BEGIN;
DROP TABLE test;
COMMIT;