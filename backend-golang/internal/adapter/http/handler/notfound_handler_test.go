package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNotFoundHandler(t *testing.T) {
	rotaInexistente := "/api/v1/rota-que-nao-existe"

	req := httptest.NewRequest(http.MethodGet, rotaInexistente, nil)
	rec := httptest.NewRecorder()

	NotFoundHandler(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, quero %d", rec.Code, http.StatusNotFound)
	}

	if contentType := rec.Header().Get("Content-Type"); contentType != "application/json" {
		t.Fatalf("Content-Type = %q, quero application/json", contentType)
	}

	var body struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("erro ao decodificar resposta: %v", err)
	}

	if body.Error.Code != CodeNotFound {
		t.Fatalf("code = %q, quero %q", body.Error.Code, CodeNotFound)
	}

	if !strings.Contains(body.Error.Message, rotaInexistente) {
		t.Fatalf("message = %q, deveria conter %q", body.Error.Message, rotaInexistente)
	}
}
