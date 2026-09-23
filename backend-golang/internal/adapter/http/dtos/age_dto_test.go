package dtos

import "testing"

func TestAgeEntryPopulation(t *testing.T) {
	tests := []struct {
		nome    string
		serie   map[string]string
		ano     string
		want    int64
		wantErr bool
	}{
		{
			nome:  "ano presente devolve o valor",
			serie: map[string]string{"2022": "1234567"},
			ano:   "2022",
			want:  1234567,
		},
		{
			nome:  "traco devolve zero",
			serie: map[string]string{"2022": "-"},
			ano:   "2022",
			want:  0,
		},
		{
			nome:    "ano ausente devolve erro",
			serie:   map[string]string{"2022": "1234567"},
			ano:     "2021",
			wantErr: true,
		},
		{
			nome:    "valor invalido devolve erro",
			serie:   map[string]string{"2022": "nao-e-numero"},
			ano:     "2022",
			wantErr: true,
		},
		{
			nome:    "valor vazio devolve erro",
			serie:   map[string]string{"2022": ""},
			ano:     "2022",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			entry := AgeEntry{
				Localidade: Localidade{Nome: "Brasil"},
				Serie:      tt.serie,
			}

			got, err := entry.Population(tt.ano)

			if tt.wantErr {
				if err == nil {
					t.Fatal("err = nil, quero erro")
				}
				return
			}

			if err != nil {
				t.Fatalf("err inesperada: %v", err)
			}

			if got != tt.want {
				t.Errorf("population = %d, quero %d", got, tt.want)
			}
		})
	}
}
