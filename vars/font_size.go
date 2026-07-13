package vars

func FontSize(s SizeVar) string {
	return "var(--fs-" + sizeStrings[s] + ")"
}
