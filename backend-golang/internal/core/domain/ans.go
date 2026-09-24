package domain

type ANSKey struct {
	Year     int
	IBGECode int64
}

type ANSBeneficiaryRow struct {
	Year          int
	IBGECode      int64
	Beneficiaries int64
}

type ANSMunicipalityMetrics struct {
	IBGECode      int64
	Name          string
	UF            string
	Region        string
	Beneficiaries int64
	Population    int64
	Penetration   float64
}

type ANSMetrics struct {
	Year           int
	Beneficiaries  int64
	Population     int64
	Penetration    float64
	Source         string
	Municipalities []ANSMunicipalityMetrics
}
