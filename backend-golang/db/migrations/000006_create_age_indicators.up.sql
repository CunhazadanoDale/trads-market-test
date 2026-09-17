CREATE TABLE IF NOT EXISTS age_indicators (
    id BIGSERIAL PRIMARY KEY,
    city_id BIGINT NOT NULL,
    year SMALLINT NOT NULL,
    age_group VARCHAR(30) NOT NULL,
    population BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_age_city
        FOREIGN KEY (city_id)
        REFERENCES cities(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_age_city_year_group
        UNIQUE (city_id, year, age_group),

    CONSTRAINT chk_age_population
        CHECK (population >= 0),

    CONSTRAINT chk_age_year
        CHECK (year >= 1900)
);

CREATE INDEX idx_age_city
    ON age_indicators(city_id);

CREATE INDEX idx_age_year
    ON age_indicators(year);

CREATE INDEX idx_age_group
    ON age_indicators(age_group);