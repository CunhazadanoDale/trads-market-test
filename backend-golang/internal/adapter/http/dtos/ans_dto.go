package dtos

import "github.com/CunhazadanoDale/trads-market-test/internal/core/domain"

type ANSMunicipalityMetricsResponse struct {
	IBGECode      int64   `json:"ibge_code"`
	Name          string  `json:"name"`
	UF            string  `json:"uf"`
	Region        string  `json:"region"`
	Beneficiaries int64   `json:"beneficiarios"`
	Population    int64   `json:"populacao"`
	Penetration   float64 `json:"penetracao"`
}

type ANSMetricsResponse struct {
	Year           int                              `json:"ano"`
	Beneficiaries  int64                            `json:"beneficiarios"`
	Population     int64                            `json:"populacao"`
	Penetration    float64                          `json:"penetracao"`
	Source         string                           `json:"fonte"`
	Municipalities []ANSMunicipalityMetricsResponse `json:"municipios"`
}

func NewANSMetricsResponse(metrics domain.ANSMetrics) ANSMetricsResponse {
	municipalities := make([]ANSMunicipalityMetricsResponse, 0, len(metrics.Municipalities))

	for _, item := range metrics.Municipalities {
		municipalities = append(municipalities, ANSMunicipalityMetricsResponse{
			IBGECode:      item.IBGECode,
			Name:          item.Name,
			UF:            item.UF,
			Region:        item.Region,
			Beneficiaries: item.Beneficiaries,
			Population:    item.Population,
			Penetration:   item.Penetration,
		})
	}

	return ANSMetricsResponse{
		Year:           metrics.Year,
		Beneficiaries:  metrics.Beneficiaries,
		Population:     metrics.Population,
		Penetration:    metrics.Penetration,
		Source:         metrics.Source,
		Municipalities: municipalities,
	}
}
