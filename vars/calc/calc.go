package calc

import "strconv"

func Mult(v string, m float64) string {
	return "calc(" + v + "*" + strconv.FormatFloat(m, 'f', 4, 64) + ")"
}
