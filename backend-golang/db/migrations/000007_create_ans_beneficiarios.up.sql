CREATE TABLE IF NOT EXISTS ans_beneficiarios (
    id BIGSERIAL PRIMARY KEY,
    ibge_code INTEGER NOT NULL,
    ano SMALLINT NOT NULL,
    beneficiarios_ativos BIGINT NOT NULL,
    fonte VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_ans_beneficiarios_ibge_ano
        UNIQUE (ibge_code, ano),

    CONSTRAINT chk_ans_beneficiarios_beneficiarios
        CHECK (beneficiarios_ativos >= 0),

    CONSTRAINT chk_ans_beneficiarios_ano
        CHECK (ano >= 1900)
);

CREATE INDEX idx_ans_beneficiarios_ibge_code
    ON ans_beneficiarios(ibge_code);

CREATE INDEX idx_ans_beneficiarios_ano
    ON ans_beneficiarios(ano);
