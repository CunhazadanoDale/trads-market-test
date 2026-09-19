package dtos

type Aggregate struct {
	ID   string `json:"id"`
	Nome string `json:"nome"`
}

type Period struct {
	Periodo string `json:"periodo"`
}

type PopulationRecord struct {
	Localidade Localidade        `json:"localidade"`
	Serie      map[string]string `json:"serie"`
}

type Localidade struct {
	ID    string `json:"id"`
	Nivel Nivel  `json:"nivel"`
	Nome  string `json:"nome"`
}

type Nivel struct {
	ID   string `json:"id"`
	Nome string `json:"nome"`
}
