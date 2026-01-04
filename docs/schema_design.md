# Database Schema Design

**Strategy:** Pure PostgreSQL (Relational). Optimized for single-tenant performance and data integrity.

## 1. Users & Authentication

```sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm"; -- For Search

-- ENUMs for Strict Typing
CREATE TYPE user_role AS ENUM ('admin', 'staff', 'customer');
CREATE TYPE order_status AS ENUM ('pending', 'processing', 'shipped', 'completed', 'cancelled');

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(100) NOT NULL,
    
    role user_role NOT NULL DEFAULT 'customer',
    
    -- RBAC for Staff
    permissions JSONB DEFAULT '[]', 
    
    phone VARCHAR(20),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Optimization: Index for login
CREATE INDEX idx_users_email ON users(email);

### Store Configuration (Key-Value Store)
-- PERFORMANCE: This table MUST be cached in MemoryStore (Read-Through pattern).
-- The app should almost never hit the DB for these keys during high traffic.
CREATE TABLE store_configs (
    key VARCHAR(100) PRIMARY KEY,
    value TEXT NOT NULL,
    is_encrypted BOOLEAN DEFAULT FALSE, -- App uses AES-GCM with 'APP_SECRET' to read/write
    group_name VARCHAR(50) NOT NULL, -- 'general', 'marketing', 'payment'
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

### Payment Gateways
CREATE TABLE payment_gateways (
    id VARCHAR(50) PRIMARY KEY, -- 'sslcommerz', 'amarpay', 'cod'
    name VARCHAR(100) NOT NULL,
    
    config JSONB NOT NULL DEFAULT '{}', -- Encrypted fields inside JSON
    is_test_mode BOOLEAN DEFAULT TRUE,
    is_active BOOLEAN DEFAULT FALSE,
    
    position INT DEFAULT 0
);

### Couriers (Shipping Providers)
CREATE TABLE couriers (
    id VARCHAR(50) PRIMARY KEY, -- 'pathao', 'steadfast', 'redx'
    name VARCHAR(100) NOT NULL,
    
    config JSONB NOT NULL DEFAULT '{}', -- Encrypted API Keys here
    is_active BOOLEAN DEFAULT FALSE,
    
    position INT DEFAULT 0
);

### Sessions (For Guest & Auth Tracking)
CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    token VARCHAR(255) UNIQUE NOT NULL, -- The Cookie Value
    user_id UUID REFERENCES users(id), -- Null if Guest
    
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

### Carts (Persistent)
CREATE TABLE carts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    session_id UUID REFERENCES sessions(id) ON DELETE CASCADE, -- Link to Guest/User Session
    user_id UUID REFERENCES users(id), -- Optional: Direct link for cross-device merge
    
    status VARCHAR(20) DEFAULT 'active', -- active, converted, abandoned
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE cart_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    cart_id UUID REFERENCES carts(id) ON DELETE CASCADE,
    product_id UUID REFERENCES products(id),
    variant_id UUID REFERENCES product_variants(id),
    
    quantity INT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

### Categories
```sql
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    parent_id UUID REFERENCES categories(id), -- Nested categories
    is_active BOOLEAN DEFAULT TRUE
);
```

### Products (The Core)
```sql
CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    base_price DECIMAL(10, 2) NOT NULL, -- Display price
    
    -- Digital Product Fields
    is_digital BOOLEAN DEFAULT TRUE,
    file_path VARCHAR(255), -- S3/R2 path for digital download
    
    -- Organization
    category_id UUID REFERENCES categories(id),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Search Optimization
CREATE INDEX idx_products_search ON products USING GIN (title gin_trgm_ops);
```

### Product Options (The "Definitions")
-- e.g. "Color" -> ["Red", "Blue"], "Size" -> ["S", "M"]
CREATE TABLE product_options (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID REFERENCES products(id) ON DELETE CASCADE,
    name VARCHAR(50) NOT NULL, -- "Color"
    position INT DEFAULT 0, -- Order of display
    values TEXT[] NOT NULL -- ["Red", "Blue", "Green"]
);

### Product Variations (The "SKUs")
-- e.g. Red-S, Red-M, Blue-S...
CREATE TABLE product_variants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID REFERENCES products(id) ON DELETE CASCADE,
    
    -- Specific values selected from Options
    -- e.g. "Red / M" (Stored as simple string for easy display)
    title VARCHAR(255) NOT NULL, 
    
    -- References for logic (Optional but cleaner)
    options JSONB, -- {"Color": "Red", "Size": "M"} 
    
    price DECIMAL(10, 2) NOT NULL,
    compare_at_price DECIMAL(10, 2), -- Original price
    sku VARCHAR(100) UNIQUE,
    stock_quantity INT DEFAULT 0,
    
    is_active BOOLEAN DEFAULT TRUE
);
-- Ensure SKU uniqueness within store
CREATE UNIQUE INDEX idx_variant_sku ON product_variants (sku);

## 3. Orders & Transactions

### Orders
```sql
CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id), -- Nullable for Guest Checkout
    guest_info JSONB, -- { "email": "...", "name": "..." } if guest
    
    total_amount DECIMAL(10, 2) NOT NULL,
    status VARCHAR(20) DEFAULT 'pending', -- pending, paid, completed, cancelled
    
    -- Marketing Intelligence
    traffic_source VARCHAR(50), -- 'facebook', 'google', 'direct'
    
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

### Order Items (Snapshot)
```sql
CREATE TABLE order_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id UUID REFERENCES orders(id) ON DELETE CASCADE,
    product_id UUID REFERENCES products(id),
    variation_id UUID REFERENCES product_variations(id), -- Nullable
    
    quantity INT NOT NULL,
    price_at_booking DECIMAL(10, 2) NOT NULL, -- Crucial: Price might change later
    
    -- Digital Fulfilment Status
    download_link_sent BOOLEAN DEFAULT FALSE
);
```

## 4. Builder & Analytics (The "Strict" System)

### Pages (Builder Config)
```sql
CREATE TABLE pages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    route VARCHAR(100) UNIQUE NOT NULL, -- e.g. '/', '/landing/xmas'
    name VARCHAR(100) NOT NULL,
    
    -- The "Strict" Configuration
    -- Array of: { "component_id": "hero", "props": {...} }
    sections JSONB NOT NULL DEFAULT '[]', 
    
    is_published BOOLEAN DEFAULT FALSE,
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
```

### A/B Testing (PDP Optimizer)
```sql
CREATE TABLE ab_tests (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL, -- e.g. "Product Page Layout Test Q1"
    target_route VARCHAR(100) NOT NULL, -- e.g. "/products/macbook-pro"
    
    variants JSONB NOT NULL, 
    -- { "A": { "layout_config":... }, "B": { "layout_config":... } }
    
    is_active BOOLEAN DEFAULT TRUE,
    start_date TIMESTAMPTZ DEFAULT NOW()
);

-- Analytics Events (CAPI Queue / Outbox)
-- Strategy: Worker picks up 'pending', sends to API, and DELETES row on success.
-- This table should remain small/empty most of the time.
CREATE TABLE analytics_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_name VARCHAR(50) NOT NULL, -- 'Purchase', 'AddToCart'
    payload JSONB NOT NULL, -- Full event data
    
    status VARCHAR(20) DEFAULT 'pending', -- pending, failed (retrying)
    providers_sent JSONB DEFAULT '[]', 
    
    created_at TIMESTAMPTZ DEFAULT NOW(),
    retry_count INT DEFAULT 0
);
-- Index for the worker to find pending jobs fast
CREATE INDEX idx_events_status ON analytics_events(status) WHERE status = 'pending';
```
