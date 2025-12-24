-- #############################################################################
-- ## DOWN MIGRATION
-- #############################################################################

-- Drop tables in reverse order of creation to respect foreign key constraints.
-- Using CASCADE automatically drops dependent objects like indexes, constraints, and triggers.
DROP TABLE IF EXISTS "leads";