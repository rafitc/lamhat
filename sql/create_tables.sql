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

-- Gallery Table 
CREATE TABLE app.gallery (
  id SERIAL PRIMARY KEY,
  user_id INT NOT NULL,
  gallery_name TEXT NOT NULL,
  gallery_status_id INT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  last_updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_user
    FOREIGN KEY (user_id) REFERENCES app.users(id),
  CONSTRAINT fk_gallery_status
    FOREIGN KEY (gallery_status_id) REFERENCES app.gallery_status(id)
);

-- Indexes for optimized queries
CREATE INDEX idx_gallery_user_id ON app.gallery(user_id);
CREATE INDEX idx_gallery_status_id ON app.gallery(gallery_status_id);
CREATE INDEX idx_gallery_created_at ON app.gallery(created_at);
ALTER TABLE app.gallery OWNER TO me; -- change owner to superuser

CREATE TABLE app.gallery_status (
  id SERIAL PRIMARY KEY,
  status TEXT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  last_updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for optimized queries
CREATE INDEX idx_gallery_status_status ON app.gallery_status(status);
CREATE INDEX idx_gallery_status_created_at ON app.gallery_status(created_at);
ALTER TABLE app.gallery_status OWNER TO me; -- change owner to superuser

-- Files table 
CREATE TABLE app.gallery_files (
  id SERIAL PRIMARY KEY,
  gallery_id INT NOT NULL,
  file_path TEXT NOT NULL,
  bucket_name TEXT NOT NULL,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  last_updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_gallery
    FOREIGN KEY (gallery_id) REFERENCES app.gallery(id)
);

-- Indexes for optimized queries
CREATE INDEX idx_gallery_files_gallery_id ON app.gallery_files(gallery_id);
CREATE INDEX idx_gallery_files_is_active ON app.gallery_files(is_active);
CREATE INDEX idx_gallery_files_created_at ON app.gallery_files(created_at);
ALTER TABLE app.gallery_files OWNER TO me; -- change owner to superuser

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

