package vars

type SizeVar int8

const (
	SizeNormal SizeVar = iota
	SizeSmall
	SizeXSmall
	SizeXXSmall
	SizeXXXSmall
	SizeLarge
	SizeXLarge
	SizeXXLarge
	SizeXXXLarge
)

var sizeStrings = map[SizeVar]string{
	SizeNormal:   "n",
	SizeSmall:    "s",
	SizeXSmall:   "xs",
	SizeXXSmall:  "xxs",
	SizeXXXSmall: "xxxs",
	SizeLarge:    "l",
	SizeXLarge:   "xl",
	SizeXXLarge:  "xxl",
	SizeXXXLarge: "xxxl",
}

func Size(s SizeVar) string {
	return "var(--s-" + sizeStrings[s] + ")"
}
