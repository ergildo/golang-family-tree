package model

type Parent struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type Ascendancy struct {
	Id      int64     `json:"id"`
	Name    string    `json:"name"`
	Depth   int       `json:"depth"`
	Parents []*Parent `json:"parents"`
}
