package handler

import "net/http"

// validRegions espelha os valores reais de states.region no banco.
var validRegions = map[string]bool{
	"Norte":        true,
	"Nordeste":     true,
	"Centro-Oeste": true,
	"Sudeste":      true,
	"Sul":          true,
}

// parseRegion lê e valida o query param ?regiao=.
// Retorna "" quando o parâmetro não foi enviado (comportamento sem filtro).
// Escreve 400 em JSON e retorna ok=false quando a região é inválida.
func parseRegion(w http.ResponseWriter, r *http.Request) (regiao string, ok bool) {
	regiao = r.URL.Query().Get("regiao")
	if regiao == "" {
		return "", true
	}

	if !validRegions[regiao] {
		writeError(
			w,
			http.StatusBadRequest,
			CodeInvalidRequest,
			"regiao deve ser uma das 5 regiões do Brasil: Norte, Nordeste, Centro-Oeste, Sudeste, Sul",
		)
		return "", false
	}

	return regiao, true
}
