package dtos

type StateRecord struct {
	ID        int64  `json:"id"`
	Sigla     string `json:"sigla"`
	Nome      string `json:"nome"`
	Regiao    RegionRecord `json:"regiao"`
}

type RegionRecord struct {
	ID   int64  `json:"id"`
	Sigla string `json:"sigla"`
	Nome string `json:"nome"`
}

type CityRecord struct {
	ID     int64  `json:"id"`
	Nome   string `json:"nome"`
	Microrregiao MicroregionRecord `json:"microrregiao"`
}

type MicroregionRecord struct {
	Mesorregiao MesoregionRecord `json:"mesorregiao"`
}

type MesoregionRecord struct {
	UF UFRecord `json:"UF"`
}

type UFRecord struct {
	ID     int64  `json:"id"`
	Sigla  string `json:"sigla"`
	Nome   string `json:"nome"`
	Regiao RegionRecord `json:"regiao"`
}