package utils

const (
	ASC  string = "asc"
	DESC        = "desc"
)

type SortItem struct {
	Field string `json:"field"`
	Dir   string `json:"dir"`
}
