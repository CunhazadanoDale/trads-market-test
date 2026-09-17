CREATE TABLE IF NOT EXISTS income_indicators (
    id BIGSERIAL PRIMARY KEY,
    city_id BIGINT NOT NULL,
    year SMALLINT NOT NULL,
    average_income NUMERIC(12, 2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_income_city
        FOREIGN KEY (city_id)
        REFERENCES cities(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_income_city_year
        UNIQUE (city_id, year),

    CONSTRAINT chk_income_positive
        CHECK (average_income >= 0),

    CONSTRAINT chk_income_year
        CHECK (year >= 1900)
);

CREATE INDEX idx_income_city
    ON income_indicators(city_id);

CREATE INDEX idx_income_year
    ON income_indicators(year);