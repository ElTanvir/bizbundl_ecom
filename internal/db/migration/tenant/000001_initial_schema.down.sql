-- Drop Analytics
DROP TABLE IF EXISTS analytics_events;
DROP TABLE IF EXISTS ab_tests;

-- Drop Builder
DROP TABLE IF EXISTS pages;

-- Drop Orders
DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS orders;

-- Drop Carts
DROP TABLE IF EXISTS cart_items;
DROP TABLE IF EXISTS carts;

-- Drop Catalog
DROP TABLE IF EXISTS product_variants;
DROP TABLE IF EXISTS product_options;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS categories;

-- Drop Sessions
DROP TABLE IF EXISTS sessions;

-- Drop Configs
DROP TABLE IF EXISTS couriers;
DROP TABLE IF EXISTS payment_gateways;
DROP TABLE IF EXISTS store_configs;

-- Drop Users
DROP TABLE IF EXISTS users;

-- Drop Types
DROP TYPE IF EXISTS order_status;
DROP TYPE IF EXISTS user_role;

-- Drop Extensions (Optional, usually keep them)
-- DROP EXTENSION IF EXISTS "pg_trgm";
-- DROP EXTENSION IF EXISTS "uuid-ossp";