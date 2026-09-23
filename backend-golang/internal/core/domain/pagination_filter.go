package domain


type PaginacaoFilter struct {
	Page int `json:"pagina"`
	Size int `json:"tamanho"`

	// Nome filtra por trecho do nome do município (case-insensitive).
	Nome string `json:"nome"`
	// Ordenar escolhe a coluna de ordenação: "", populacao, renda ou pib.
	Ordenar string `json:"ordenar"`
	// Ordem define a direção: "", asc ou desc.
	Ordem string `json:"ordem"`
}


type PaginacaoResponse[T any] struct {
	Dados []T `json:"dados"`
	Page int `json:"pagina"`
	Size int `json:"tamanho"`
	Total int `json:"total"`
}