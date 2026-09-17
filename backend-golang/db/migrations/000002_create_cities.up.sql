CREATE TABLE cities (
    id BIGSERIAL PRIMARY KEY,
    ibge_code INTEGER NOT NULL UNIQUE,
    state_id BIGINT NOT NULL,
    name VARCHAR(150) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_cities_state
        FOREIGN KEY (state_id)
        REFERENCES states(id)
        ON DELETE RESTRICT
);

CREATE INDEX idx_cities_state_id
    ON cities(state_id);

CREATE INDEX idx_cities_name
    ON cities(name);