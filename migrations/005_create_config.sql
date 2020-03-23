BEGIN;

CREATE OR REPLACE FUNCTION enforce_one_config() RETURNS TRIGGER AS $$
    DECLARE
        config_count INTEGER := 0;
    BEGIN
        -- Lock the table to prevent concurrent `setup` commands from both
        -- making it into the DB.
        LOCK TABLE config IN EXCLUSIVE MODE;

        SELECT INTO config_count COUNT(*)
        FROM config;

        IF config_count > 0 THEN
            RAISE EXCEPTION 'Cannot insert more than 1 config object in the database';
        END IF;

        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE TABLE config (
    id char(29) primary key not null,
    data jsonb not null,

    -- If the config is in the table, the setup flag must be true
    CONSTRAINT config_setup CHECK ((data::jsonb#>>'{setup}')::BOOLEAN = true),
    CONSTRAINT config_id_check CHECK (data::jsonb#>>'{id}' = id)
);

CREATE TRIGGER config_hoist_tgr
  BEFORE INSERT ON config
  FOR EACH ROW EXECUTE PROCEDURE hoist_values('id');

CREATE TRIGGER config_check_insert_tgr
    BEFORE INSERT ON config
    FOR EACH ROW EXECUTE PROCEDURE enforce_one_config();

COMMIT;

---- create above / drop below ----

BEGIN;
DROP TRIGGER config_check_insert_tgr ON config;
DROP TRIGGER config_hoist_tgr ON config;
DROP FUNCTION enforce_one_config();
DROP TABLE CONFIG;
COMMIT;