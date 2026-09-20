package domain


type City struct {
	ID   int64    `json:"id"`
	IBGECode int64      `json:"ibge_code"`
	Name string   `json:"name"`
	StateID int64  `json:"state_id"`
}