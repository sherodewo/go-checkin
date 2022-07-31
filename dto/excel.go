package dto

type Excel struct {
	Start string `json:"start" form:"start" param:"start"`
	End   string `json:"end" form:"end" param:"end"`
	ID    string `json:"id" form:"id" param:"id"`
}
