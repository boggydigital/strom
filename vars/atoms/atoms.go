package atoms

import "github.com/boggydigital/strom/vars/sizes"

type Atom int8

const (
	DisplayFlex Atom = iota

	FlexFlowRowWrap
	FlexFlowRowNoWrap
	FlexFlowColWrap
	FlexFlowColNoWrap

	FlexDirColumn
	FlexDirRow

	AlignItemsStart
	AlignItemsCenter
	AlignItemsEnd

	AlignContentStart
	AlignContentCenter
	AlignContentEnd

	JustifyItemsStart
	JustifyItemsCenter
	JustifyItemsEnd

	JustifyContentStart
	JustifyContentCenter
	JustifyContentEnd

	RowGapNormal
	RowGapSmall
	RowGapLarge

	ColGapNormal
	ColGapSmall
	ColGapLarge

	PaddingNormal
	PaddingSmall
	PaddingLarge

	MarginNormal
	MarginSmall
	MarginLarge

	FontSizeNormal
	FontSizeSmall
	FontSizeLarge

	FontWeightNormal
	FontWeightLight
	FontWeightBold

	BorderRadiusNormal
	BorderRadiusSmall
	BorderRadiusLarge
)

var atomicClasses = map[Atom]string{
	DisplayFlex: "d-f",

	FlexFlowRowWrap:   "ff-rw",
	FlexFlowRowNoWrap: "ff-rnw",
	FlexFlowColWrap:   "ff-cw",
	FlexFlowColNoWrap: "ff-cnw",

	FlexDirColumn: "fd-c",
	FlexDirRow:    "fd-r",

	AlignItemsStart:  "ai-s",
	AlignItemsCenter: "ai-c",
	AlignItemsEnd:    "ai-e",

	AlignContentStart:  "ac-s",
	AlignContentCenter: "ac-c",
	AlignContentEnd:    "ac-e",

	JustifyItemsStart:  "ji-s",
	JustifyItemsCenter: "ji-c",
	JustifyItemsEnd:    "ji-e",

	JustifyContentStart:  "jc-s",
	JustifyContentCenter: "jc-c",
	JustifyContentEnd:    "jc-e",

	RowGapNormal: "rg-n",
	RowGapSmall:  "rg-s",
	RowGapLarge:  "rg-l",

	ColGapNormal: "cg-n",
	ColGapSmall:  "cg-s",
	ColGapLarge:  "cg-l",

	PaddingNormal: "p-n",
	PaddingSmall:  "p-s",
	PaddingLarge:  "p-l",

	MarginNormal: "m-n",
	MarginSmall:  "m-s",
	MarginLarge:  "m-l",

	FontSizeNormal: "fs-n",
	FontSizeSmall:  "fs-s",
	FontSizeLarge:  "fs-l",

	FontWeightNormal: "fw-n",
	FontWeightLight:  "fw-l",
	FontWeightBold:   "fw-b",

	BorderRadiusNormal: "br-n",
	BorderRadiusSmall:  "br-s",
	BorderRadiusLarge:  "br-l",
}

func (a Atom) Class() string {
	return atomicClasses[a]
}

func FlexRowWrap(gap sizes.Size) []Atom {
	return flexFlow(FlexFlowRowWrap, gap)
}

func FlexRow(gap sizes.Size) []Atom {
	return flexFlow(FlexFlowRowNoWrap, gap)
}

func FlexColWrap(gap sizes.Size) []Atom {
	return flexFlow(FlexFlowColWrap, gap)
}

func FlexCol(gap sizes.Size) []Atom {
	return flexFlow(FlexFlowColNoWrap, gap)
}

func flexFlow(ff Atom, gap sizes.Size) []Atom {

	ffe := []Atom{DisplayFlex, ff}

	switch gap {
	case sizes.Normal:
		ffe = append(ffe, RowGapNormal, ColGapNormal)
	case sizes.Small:
		ffe = append(ffe, RowGapSmall, ColGapSmall)
	case sizes.Large:
		ffe = append(ffe, RowGapLarge, ColGapLarge)
	}

	return ffe
}
