package domain


type PaginacaoFilter struct {
	Page int `json:"pagina"`
	Size int `json:"tamanho"`

	Nome    string `json:"nome"`
	Ordenar string `json:"ordenar"`
	Ordem   string `json:"ordem"`
}


type PaginacaoResponse[T any] struct {
	Dados []T `json:"dados"`
	Page int `json:"pagina"`
	Size int `json:"tamanho"`
	Total int `json:"total"`
}