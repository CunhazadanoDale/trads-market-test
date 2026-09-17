CREATE TABLE IF NOT EXISTS states (
    id BIGSERIAL PRIMARY KEY,
    ibge_code INTEGER NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    uf CHAR(2) NOT NULL UNIQUE,
    region VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
    
CREATE INDEX idx_states_region
    ON states(region);