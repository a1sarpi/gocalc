package constants

import "math"

const (
	Pi = math.Pi 
	E  = math.E
)

func GetConstant(name string) (float64, bool) {
	switch name {
	case "pi":
		return Pi, true
	case "e":
		return E, true
	default:
		return 0, false
	}
}
