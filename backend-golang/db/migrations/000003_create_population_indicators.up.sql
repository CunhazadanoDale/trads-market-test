CREATE TABLE population_indicators (
    id BIGSERIAL PRIMARY KEY,
    city_id BIGINT NOT NULL,
    year SMALLINT NOT NULL,
    population BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_population_city
        FOREIGN KEY (city_id)
        REFERENCES cities(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_population_city_year
        UNIQUE (city_id, year),

    CONSTRAINT chk_population_positive
        CHECK (population >= 0),

    CONSTRAINT chk_population_year
        CHECK (year >= 1900)
);

CREATE INDEX idx_population_city
    ON population_indicators(city_id);

CREATE INDEX idx_population_year
    ON population_indicators(year);