-- Users Table 

CREATE TABLE app.users (
	id serial4 NOT NULL,
	email_id varchar(255) NOT NULL,
	is_email_valid bool NOT NULL DEFAULT false,
	first_name varchar(100) NULL,
	last_name varchar(100) NULL,
	auth_key_hash bpchar(64) NULL,
	otp bpchar(6) NULL,
	otp_generated_at timestamptz NULL,
	is_user_active bool NULL DEFAULT true,
	created_at timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
	last_updated_at timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT users_auth_key_hash_key UNIQUE (auth_key_hash),
	CONSTRAINT users_email_id_key UNIQUE (email_id),
	CONSTRAINT users_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_user_id ON app.users USING btree (id);
CREATE INDEX idx_users_email_id ON app.users USING btree (email_id);
CREATE INDEX idx_users_is_user_active ON app.users USING btree (is_user_active);

ALTER TABLE app.users OWNER TO me; -- change owner to superuser

-- Give access to the backend user 
-- To all tables
DO $$
DECLARE
    r RECORD;
BEGIN
    FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = 'app') LOOP
        EXECUTE 'GRANT ALL PRIVILEGES ON TABLE app.' || quote_ident(r.tablename) || ' TO lamhat_backend_dev;';
    END LOOP;
END $$;

-- to all sequences
DO $$
DECLARE
    r RECORD;
BEGIN
    FOR r IN (SELECT sequencename FROM pg_sequences WHERE schemaname = 'app') LOOP
        EXECUTE 'GRANT ALL PRIVILEGES ON SEQUENCE app.' || quote_ident(r.sequencename) || ' TO lamhat_backend_dev;';
    END LOOP;
END $$;

-- to all future tables 
ALTER DEFAULT PRIVILEGES IN SCHEMA app GRANT ALL PRIVILEGES ON TABLES TO lamhat_backend_dev;
ALTER DEFAULT PRIVILEGES IN SCHEMA app GRANT ALL PRIVILEGES ON SEQUENCES TO lamhat_backend_dev;

