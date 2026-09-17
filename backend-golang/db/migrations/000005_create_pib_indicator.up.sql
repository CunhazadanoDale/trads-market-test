CREATE TABLE IF NOT EXISTS gdp_indicators (
    id BIGSERIAL PRIMARY KEY,
    city_id BIGINT NOT NULL,
    year SMALLINT NOT NULL,
    gdp NUMERIC(18, 2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_gdp_city
        FOREIGN KEY (city_id)
        REFERENCES cities(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_gdp_city_year
        UNIQUE (city_id, year),

    CONSTRAINT chk_gdp_positive
        CHECK (gdp >= 0),

    CONSTRAINT chk_gdp_year
        CHECK (year >= 1900)
);

CREATE INDEX idx_gdp_city
    ON gdp_indicators(city_id);

CREATE INDEX idx_gdp_year
    ON gdp_indicators(year);

