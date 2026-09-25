package config

import (
	"strings"
	"testing"
)

func TestLoadConfigDefaults(t *testing.T) {
	t.Setenv("PORT", "")
	t.Setenv("APP_ENV", "")

	cfg := LoadConfig()

	if cfg.Port != defaultPort {
		t.Errorf("port = %q, quero %q", cfg.Port, defaultPort)
	}

	if cfg.AppEnv != defaultAppEnv {
		t.Errorf("appEnv = %q, quero %q", cfg.AppEnv, defaultAppEnv)
	}
}

func TestLoadConfigRespeitaVariaveisDoAmbiente(t *testing.T) {
	t.Setenv("PORT", "9090")
	t.Setenv("APP_ENV", "production")

	cfg := LoadConfig()

	if cfg.Port != "9090" {
		t.Errorf("port = %q, quero %q", cfg.Port, "9090")
	}

	if cfg.AppEnv != "production" {
		t.Errorf("appEnv = %q, quero %q", cfg.AppEnv, "production")
	}
}

func TestValidateAPI(t *testing.T) {
	tests := []struct {
		nome     string
		cfg      Config
		wantErro string
	}{
		{
			nome: "com DATABASE_URL a configuração é válida",
			cfg:  Config{DatabaseUrl: "postgres://trads:trads@localhost/trads"},
		},
		{
			nome:     "sem DATABASE_URL aponta a variável",
			cfg:      Config{},
			wantErro: "DATABASE_URL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			err := tt.cfg.ValidateAPI()

			if tt.wantErro == "" {
				if err != nil {
					t.Fatalf("err = %v, quero nil", err)
				}
				return
			}

			if err == nil {
				t.Fatal("err = nil, quero erro mencionando a variável")
			}
			if !strings.Contains(err.Error(), tt.wantErro) {
				t.Errorf("err = %v, quero menção a %s", err, tt.wantErro)
			}
		})
	}
}

func TestValidateImport(t *testing.T) {
	ibge := Config{
		DatabaseUrl:        "postgres://trads",
		BaseUrlIBGE:        "https://servicodados.ibge.gov.br/api/v3",
		BaseUrlLocalidades: "https://servicodados.ibge.gov.br/api/v1",
	}
	ans := Config{
		DatabaseUrl: "postgres://trads",
		BaseUrlANS:  "https://dadosabertos.ans.gov.br/pda.csv",
		AnsFonte:    "ANS PDA-047",
	}
	all := Config{
		DatabaseUrl:        "postgres://trads",
		BaseUrlIBGE:        ibge.BaseUrlIBGE,
		BaseUrlLocalidades: ibge.BaseUrlLocalidades,
		BaseUrlANS:         ans.BaseUrlANS,
		AnsFonte:           ans.AnsFonte,
	}

	tests := []struct {
		nome     string
		cfg      Config
		target   string
		wantErro string
	}{
		{
			nome:   "ibge com as duas bases do ibge é válido",
			cfg:    ibge,
			target: "ibge",
		},
		{
			nome:   "ans com url e fonte é válido",
			cfg:    ans,
			target: "ans",
		},
		{
			nome:   "all com as duas configurações é válido",
			cfg:    all,
			target: "all",
		},
		{
			nome:     "ibge sem BASE_URL_IBGE aponta a variável",
			cfg:      Config{DatabaseUrl: "postgres://trads", BaseUrlLocalidades: ibge.BaseUrlLocalidades},
			target:   "ibge",
			wantErro: "BASE_URL_IBGE",
		},
		{
			nome:     "ans sem BASE_URL_ANS aponta a variável",
			cfg:      Config{DatabaseUrl: "postgres://trads", AnsFonte: ans.AnsFonte},
			target:   "ans",
			wantErro: "BASE_URL_ANS",
		},
		{
			nome:     "ans sem ANS_FONTE aponta a variável",
			cfg:      Config{DatabaseUrl: "postgres://trads", BaseUrlANS: ans.BaseUrlANS},
			target:   "ans",
			wantErro: "ANS_FONTE",
		},
		{
			nome:     "sem DATABASE_URL aponta a variável",
			cfg:      Config{BaseUrlIBGE: ibge.BaseUrlIBGE, BaseUrlLocalidades: ibge.BaseUrlLocalidades},
			target:   "ibge",
			wantErro: "DATABASE_URL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			err := tt.cfg.ValidateImport(tt.target)

			if tt.wantErro == "" {
				if err != nil {
					t.Fatalf("err = %v, quero nil", err)
				}
				return
			}

			if err == nil {
				t.Fatal("err = nil, quero erro mencionando a variável")
			}
			if !strings.Contains(err.Error(), tt.wantErro) {
				t.Errorf("err = %v, quero menção a %s", err, tt.wantErro)
			}
		})
	}
}
