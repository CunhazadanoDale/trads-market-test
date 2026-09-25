package ibge

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

const fixtureSIDRA = `[{"resultados":[{"series":[{"localidade":{"id":"4314506","nivel":{"id":"N6","nome":"Município"},"nome":"Pinheiro Machado - RS"},"serie":{"2022":"11214"}}]}]}]`

func TestGetPopulation2022ComFixtureSIDRA(t *testing.T) {
	servidor := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/agregados/4709/periodos/2022/variaveis/93" {
			t.Errorf("path = %q, quero o agregado 4709", r.URL.Path)
		}
		if !strings.Contains(r.URL.RawQuery, "localidades") {
			t.Errorf("query = %q, quero o filtro de localidades", r.URL.RawQuery)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(fixtureSIDRA))
	}))
	defer servidor.Close()

	cliente := NewIbgeClient(servidor.URL, servidor.URL, servidor.Client())

	registros, err := cliente.GetPopulation2022(context.Background())
	if err != nil {
		t.Fatalf("err = %v, quero nil", err)
	}

	if len(registros) != 1 {
		t.Fatalf("registros = %d, quero 1", len(registros))
	}
	if registros[0].Localidade.ID != "4314506" {
		t.Errorf("localidade = %q, quero 4314506", registros[0].Localidade.ID)
	}

	populacao, err := PopulationByYear(registros[0], "2022")
	if err != nil {
		t.Fatalf("err = %v, quero nil", err)
	}
	if populacao != 11214 {
		t.Errorf("populacao = %d, quero 11214", populacao)
	}
}

func TestGetPopulation2022RetentaAteSucesso(t *testing.T) {
	var tentativas atomic.Int32

	servidor := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tentativas.Add(1) < 3 {
			http.Error(w, "erro interno", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(fixtureSIDRA))
	}))
	defer servidor.Close()

	cliente := NewIbgeClient(servidor.URL, servidor.URL, servidor.Client())

	registros, err := cliente.GetPopulation2022(context.Background())
	if err != nil {
		t.Fatalf("err = %v, quero nil depois das retentativas", err)
	}
	if tentativas.Load() != 3 {
		t.Errorf("tentativas = %d, quero 3", tentativas.Load())
	}
	if len(registros) != 1 {
		t.Fatalf("registros = %d, quero 1", len(registros))
	}
}

func TestGetPopulation2022NaoRetentaErroDeCliente(t *testing.T) {
	var tentativas atomic.Int32

	servidor := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tentativas.Add(1)
		http.Error(w, "não encontrado", http.StatusNotFound)
	}))
	defer servidor.Close()

	cliente := NewIbgeClient(servidor.URL, servidor.URL, servidor.Client())

	_, err := cliente.GetPopulation2022(context.Background())
	if err == nil {
		t.Fatal("err = nil, quero erro de status 404")
	}
	if !strings.Contains(err.Error(), "404") {
		t.Errorf("err = %v, quero menção ao status 404", err)
	}
	if tentativas.Load() != 1 {
		t.Errorf("tentativas = %d, quero 1", tentativas.Load())
	}
}

func TestGetPopulation2022InterrompeNoTimeout(t *testing.T) {
	bloqueado := make(chan struct{})

	servidor := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-bloqueado
	}))
	defer func() {
		close(bloqueado)
		servidor.Close()
	}()

	cliente := NewIbgeClient(servidor.URL, servidor.URL, servidor.Client())

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err := cliente.GetPopulation2022(ctx)
	if err == nil {
		t.Fatal("err = nil, quero erro de timeout")
	}
	if !strings.Contains(err.Error(), "interrompida") {
		t.Errorf("err = %v, quero menção à tentativa interrompida", err)
	}
}
