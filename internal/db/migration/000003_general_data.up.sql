
-- Create a table to store general_data
CREATE TABLE general_data (
    name TEXT PRIMARY KEY, -- Name as primary key since names are unique
    value TEXT NOT NULL,   -- The CSS value, e.g., #8b5cf6
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
