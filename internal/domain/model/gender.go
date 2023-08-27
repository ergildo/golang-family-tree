package model

type Gender string

const (
	Male   = "M"
	Female = "F"
	Others = ""
)

func ParseGender(value string) Gender {
	switch value {
	case "M":
		return Male
	case "F":
		return Female
	default:
		return Others
	}
}
