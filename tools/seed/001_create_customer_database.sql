BEGIN;

CREATE TABLE transactions (
    id serial primary key,
    processor text,
    -- random time within the last year-ish
    timestamp timestamp default NOW() - ('1 year'::INTERVAL * ROUND(RANDOM()) - - '1 day'::INTERVAL * ROUND(RANDOM() * 365) - '1 hour'::INTERVAL * ROUND(RANDOM() * 24)),
    card_id int,
    card_number bigint,
    value float,
    ssn integer,
    vendor text
);

COMMIT;


---- create above / drop below ----

BEGIN;
DROP TABLE transactions;
COMMIT;
