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
