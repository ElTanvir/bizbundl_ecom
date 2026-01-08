-- Enable Extensions for Public Schema
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 1. Platform Users (Owners)
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 2. Shops
CREATE TABLE shops (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    owner_id UUID REFERENCES users(id),
    name VARCHAR(100) NOT NULL,
    subdomain VARCHAR(100) UNIQUE NOT NULL, -- e.g. "neon-vibes"
    custom_domain VARCHAR(255) UNIQUE,      -- e.g. "neonvibes.com"
    
    tenant_id VARCHAR(100) UNIQUE NOT NULL, -- e.g. "shop_neon_vibes" (Schema Name)
    
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_shops_subdomain ON shops(subdomain);
CREATE INDEX idx_shops_custom_domain ON shops(custom_domain);

-- 3. Subscriptions (Simple MVP)
CREATE TABLE subscriptions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    shop_id UUID REFERENCES shops(id),
    plan_name VARCHAR(50) NOT NULL, -- "starter", "pro"
    status VARCHAR(20) DEFAULT 'active',
    current_period_end TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
