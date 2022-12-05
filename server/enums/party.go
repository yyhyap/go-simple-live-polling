package enums

import "strings"

type Party string

var (
	partiesMap = map[string]Party{
		"ph": PakatanHarapan,
		"bn": BarisanNasional,
		"pn": PerikatanNasional,
	}
)

const (
	PakatanHarapan    Party = "PakatanHarapan"
	BarisanNasional   Party = "BarisanNasional"
	PerikatanNasional Party = "PerikatanNasional"
)

func (p Party) String() string {
	switch p {
	case PakatanHarapan:
		return "PakatanHarapan"
	case BarisanNasional:
		return "BarisanNasional"
	case PerikatanNasional:
		return "PerikatanNasional"
	default:
		return "unknown"
	}
}

func ParseString(s string) (Party, bool) {
	p, exist := partiesMap[strings.ToLower(s)]
	return p, exist
}
