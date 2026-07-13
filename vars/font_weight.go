package vars

type WeightVar int8

const (
	WeightNormal WeightVar = iota
	WeightLight
	WeightBold
)

var weightStrings = map[WeightVar]string{
	WeightNormal: "n",
	WeightLight:  "l",
	WeightBold:   "b",
}

func FontWeight(w WeightVar) string {
	return "var(--fw-" + weightStrings[w] + ")"
}
