package model

type Person struct {
	Name     string  `json:"name"`
	Parent   int64   `json:"parent"`
	Children []int64 `json:"children"`
}
