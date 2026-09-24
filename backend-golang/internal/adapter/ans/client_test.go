package ans

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const testHeader = "PERIODO;CD_MUNICIPIO;NM_MUNICIPIO;CD_UF;SG_UF;CD_RM;NM_RM;SEXO;FAIXA_ETARIA;BENEF_ASSISTENCIA_MEDICA;BENEF_EXCLUS_ODONTOLOGICO;BENEF_TOTAL;POPULACAO;TX_COBERT_ASSISTENCIA_MEDICA;TX_COBERT_EXCLUSIVAMENTE_ODONTOLOGICO;TX_COBERT_TOTAL"

func TestClientFetchParsesRowsAndSkipsInvalidOnes(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(testHeader + "\n"))
		w.Write([]byte("\"2026\";\"120005\";\"Assis Brasil\";\"12\";\"AC\";\"00000\";\"Fora\";\"FEMININO\";\"Até 1 ano\";\"150\";\"200\";\"350\";\"0\";0,0000;0,0000;0,0000\n"))
		w.Write([]byte("\"2026\";\"120005\";\"Assis Brasil\";\"12\";\"AC\";\"00000\";\"Fora\";\"MASCULINO\";\"Até 1 ano\";\"1.500\";\"200\";\"1.700\";\"0\";0,0000;0,0000;0,0000\n"))
		w.Write([]byte("\"2026\";\"-1\";\"Não informado\";\"\";\"\";\"\";\"\";\"FEMININO\";\"Até 1 ano\";\"10\";\"0\";\"10\";\"0\";0,0000;0,0000;0,0000\n"))
		w.Write([]byte("\"1800\";\"120005\";\"Assis Brasil\";\"12\";\"AC\";\"00000\";\"Fora\";\"FEMININO\";\"Até 1 ano\";\"99\";\"0\";\"99\";\"0\";0,0000;0,0000;0,0000\n"))
		w.Write([]byte("linha-quebrada\n"))
	}))
	defer srv.Close()

	client := NewClient(srv.URL)
	rows, err := client.Fetch(context.Background())
	if err != nil {
		t.Fatalf("Fetch() error = %v", err)
	}

	if len(rows) != 2 {
		t.Fatalf("len(rows) = %d, want 2", len(rows))
	}
	if rows[0].Year != 2026 || rows[0].IBGECode != 120005 {
		t.Errorf("rows[0] = %+v, want {Year: 2026 IBGECode: 120005}", rows[0])
	}
	if rows[0].Beneficiaries != 150 {
		t.Errorf("rows[0].Beneficiaries = %d, want 150", rows[0].Beneficiaries)
	}
	if rows[1].Beneficiaries != 1500 {
		t.Errorf("rows[1].Beneficiaries = %d, want 1500", rows[1].Beneficiaries)
	}
}

func TestClientFetchRejectsInvalidHeader(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("COLUNA_A;COLUNA_B\n1;2\n"))
	}))
	defer srv.Close()

	client := NewClient(srv.URL)
	_, err := client.Fetch(context.Background())
	if err == nil {
		t.Fatal("Fetch() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "inesperado") {
		t.Errorf("error = %v, want menção a cabeçalho inesperado", err)
	}
}

func TestClientFetchErrorsWhenNoValidRows(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(testHeader + "\n"))
		w.Write([]byte("linha-quebrada\n"))
	}))
	defer srv.Close()

	client := NewClient(srv.URL)
	_, err := client.Fetch(context.Background())
	if err == nil {
		t.Fatal("Fetch() error = nil, want error")
	}
}

func TestClientFetchErrorsWhenURLNotConfigured(t *testing.T) {
	client := NewClient("")
	_, err := client.Fetch(context.Background())
	if err == nil {
		t.Fatal("Fetch() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "BASE_URL_ANS") {
		t.Errorf("error = %v, want menção a BASE_URL_ANS", err)
	}
}

func TestClientFetchFailsOnHTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}))
	defer srv.Close()

	client := NewClient(srv.URL)
	_, err := client.Fetch(context.Background())
	if err == nil {
		t.Fatal("Fetch() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "HTTP 404") {
		t.Errorf("error = %v, want menção a HTTP 404", err)
	}
}
