-- #############################################################################
-- ## UP MIGRATION (optimized)
-- #############################################################################


-- Enable the uuid-ossp extension to use uuid_generate_v4()
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- Enable pg_trgm extension for efficient text search
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- -----------------------------------------------------------------------------
-- -- Function to automatically update 'updated_at' timestamps
-- -----------------------------------------------------------------------------
-- This trigger function is designed to be called before any update on a table.
-- It sets the 'updated_at' column of the row being updated to the current time.
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = now(); 
   RETURN NEW;
END;
$$ language 'plpgsql';

-- -----------------------------------------------------------------------------
-- -- Enums
-- -----------------------------------------------------------------------------
CREATE TYPE "user_role" AS ENUM ('user', 'moderator', 'admin');

-- -----------------------------------------------------------------------------
-- -- Users Table
-- -----------------------------------------------------------------------------
-- This table stores user account information.
CREATE TABLE "users" (
    "id"            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "username"      VARCHAR(255) UNIQUE,
    "hashed_password" VARCHAR(255) NOT NULL,
    "first_name"    VARCHAR(100),
    "last_name"     VARCHAR(100),
    "full_name"     VARCHAR(201) GENERATED ALWAYS AS (COALESCE(first_name, '') || ' ' || COALESCE(last_name, '')) STORED,
    "email"         VARCHAR(255) UNIQUE,
    "phone"         VARCHAR(20) UNIQUE,
    "role"          "user_role" NOT NULL DEFAULT 'user',
    "is_email_verified" BOOLEAN NOT NULL DEFAULT FALSE,
    "is_active"     BOOLEAN NOT NULL DEFAULT TRUE,
    "created_at"    TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at"    TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deleted_at"    TIMESTAMPTZ
);

-- Indexing for the 'users' table
-- Note: Unique columns (username, email, phone) automatically get indexes, so we don't need to explicitly create them.

-- Index for filtering by role, useful for admin dashboards
CREATE INDEX ON "users" ("role");

-- Index for sorting by creation time, useful for "newest users" queries
CREATE INDEX ON "users" ("created_at");

-- Partial index for soft deletes. This improves performance for queries that
-- filter out deleted users, which is a very common operation.
CREATE INDEX ON "users" ("deleted_at") WHERE "deleted_at" IS NULL;

-- GIN indexes for efficient text search using pg_trgm
CREATE INDEX ON "users" USING GIN ("username" gin_trgm_ops);
CREATE INDEX ON "users" USING GIN ("email" gin_trgm_ops);
CREATE INDEX ON "users" USING GIN ("full_name" gin_trgm_ops);

-- Trigger to automatically update the 'updated_at' field on user record changes.
CREATE TRIGGER update_users_updated_at
BEFORE UPDATE ON "users"
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
