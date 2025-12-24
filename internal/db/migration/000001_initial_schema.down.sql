-- #############################################################################
-- ## DOWN MIGRATION
-- #############################################################################
DROP TABLE IF EXISTS "users" CASCADE;
-- Drop the trigger function.
DROP FUNCTION IF EXISTS update_updated_at_column();
-- Drop the custom type
DROP TYPE IF EXISTS "user_role";
-- Drop extensions
DROP EXTENSION IF EXISTS "pg_trgm";