package model

type Person struct {
	Name     string  `json:"name"`
	Parent   int64  `json:"parent"`
	Gender   string  `json:"gender"`
	Children []int64 `json:"children"`
}

func (p Person) GetGender() Gender {
	return ParseGender(p.Gender)
}

func (p Person) IsFather() bool {
	return Male == p.GetGender()
}

func (p Person) IsMother() bool {
	return Female == p.GetGender()
}
