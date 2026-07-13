package vars

type ColorVar int8

const (
	ColorRed ColorVar = iota
	ColorOrange
	ColorYellow
	ColorGreen
	ColorMint
	ColorTeal
	ColorCyan
	ColorBlue
	ColorIndigo
	ColorPurple
	ColorPink
	ColorBrown
	ColorGray
	ColorBackground
	ColorForeground
	ColorHighlight
)

var colorStrings = map[ColorVar]string{
	ColorRed:        "red",
	ColorOrange:     "orange",
	ColorYellow:     "yellow",
	ColorGreen:      "green",
	ColorMint:       "mint",
	ColorTeal:       "teal",
	ColorCyan:       "cyan",
	ColorBlue:       "blue",
	ColorIndigo:     "indigo",
	ColorPurple:     "purple",
	ColorPink:       "pink",
	ColorBrown:      "brown",
	ColorGray:       "gray",
	ColorBackground: "background",
	ColorForeground: "foreground",
	ColorHighlight:  "highlight",
}

func Color(c ColorVar) string {
	return "var(--c-" + colorStrings[c] + ")"
}
