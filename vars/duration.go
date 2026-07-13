package vars

type DurationVar int8

const (
	DurationNormal DurationVar = iota
	DurationFast
	DurationSlow
)

var durationStrings = map[DurationVar]string{
	DurationNormal: "n",
	DurationFast:   "f",
	DurationSlow:   "s",
}

func Duration(d DurationVar) string {
	return "var(--du-" + durationStrings[d] + ")"
}
