
-- Create a table to store individual CSS variables
CREATE TABLE css_variables (
    name TEXT PRIMARY KEY, -- Name as primary key since names are unique
    value TEXT NOT NULL,   -- The CSS value, e.g., #8b5cf6
    variable_type TEXT NOT NULL, -- For organization, e.g., color, font, radius
    is_system BOOLEAN NOT NULL DEFAULT FALSE, -- System variables cannot be deleted
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Index for fast type lookups
CREATE INDEX idx_css_variables_type ON css_variables (variable_type);

-- Insert all variables from the previous 'default' theme
INSERT INTO css_variables (name, value, variable_type, is_system) VALUES
-- Colors
('--color-surface', '#ffffff', 'color', TRUE),
('--color-surface-alt', '#fbf5ff', 'color', TRUE),
('--color-on-surface', '#555555', 'color', TRUE),
('--color-on-surface-strong', '#212121', 'color', TRUE),
('--color-primary', '#7f00ff', 'color', TRUE),
('--color-on-primary', '#ffffff', 'color', TRUE),
('--color-secondary', '#ccff00', 'color', TRUE),
('--color-on-secondary', '#212121', 'color', TRUE),
('--color-outline', '#d4d4d4', 'color', TRUE),
('--color-outline-strong', '#212121', 'color', TRUE),
('--color-info', '#0284c7', 'color', TRUE),
('--color-on-info', '#000000', 'color', TRUE),
('--color-success', '#059669', 'color', TRUE),
('--color-on-success', '#000000', 'color', TRUE),
('--color-warning', '#ff4500', 'color', TRUE),
('--color-on-warning', '#ffffff', 'color', TRUE),
('--color-danger', '#fff5f5', 'color', TRUE),
('--color-on-danger', '#ef4444', 'color', TRUE),
-- Fonts
('--font-body', 'Inter, sans-serif', 'font', TRUE),
('--font-title', 'Nunito, sans-serif', 'font', TRUE),
-- Radius
('--radius-none', '0', 'radius', TRUE),
('--radius-radius', '1rem', 'radius', TRUE),
-- Shadows
('--shadow-sm', '0 1px 2px 0 rgb(0 0 0 / 0.05)', 'shadow', TRUE),
('--shadow-md', '0 4px 6px -1px rgb(0 0 0 / 0.1), 0 2px 4px -2px rgb(0 0 0 / 0.1)', 'shadow', TRUE),
('--shadow-lg', '0 10px 15px -3px rgb(0 0 0 / 0.1), 0 4px 6px -4px rgb(0 0 0 / 0.1)', 'shadow', TRUE),
('--shadow-xl', '0 20px 25px -5px rgb(0 0 0 / 0.1), 0 8px 10px -6px rgb(0 0 0 / 0.1)', 'shadow', TRUE);
