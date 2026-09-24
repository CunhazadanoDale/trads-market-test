package domain

type State struct {
	ID       int64  `db:"id"`
	IBGECode int64  `db:"ibge_code"`
	Name     string `db:"name"`
	UF       string `db:"uf"`
	Region   string `db:"region"`
}
