ALTER TABLE products ADD COLUMN is_featured BOOLEAN DEFAULT false;
CREATE INDEX idx_products_featured ON products(is_featured) WHERE is_featured = TRUE;
