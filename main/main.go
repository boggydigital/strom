package main

import (
	"iter"
	"net/http"
	"strconv"

	"github.com/boggydigital/strom"
	"github.com/boggydigital/strom/vars"
)

func main() {
	http.HandleFunc("GET /", GetTest)
	if err := http.ListenAndServe(":12345", nil); err != nil {
		panic(err)
	}
}

func GetTest(w http.ResponseWriter, r *http.Request) {

	root := strom.Page("test")

	var body strom.Element
	for body = range root.GetElementsByTagName("body") {
		break
	}

	body.Append(
		strom.Defer(aboveTheFold),
		strom.Defer(belowTheFold))

	//for ii := range 255 {
	//	iistr := strconv.Itoa(ii)
	//	root.Append(strom.Create("div", "Node "+iistr).SetStyles(map[string]string{
	//		"color": "rgb(" + strconv.FormatInt(int64(ii), 10) + ",0,0)",
	//	}))
	//	root.Append(strom.Create("span").
	//		SetTextContent("Test Text"))
	//}
	//
	//for ii := range 10000 {
	//
	//	iistr := strconv.Itoa(ii)
	//
	//	root.Append(strom.Create("div", "Node "+iistr).
	//		AddClass("test-deferred-class").
	//		SetAttribute("id", iistr).
	//		SetTextContent("Deferred Node " + iistr))
	//}

	if err := root.Write(w); err != nil {
		panic(err)
	}
}

func aboveTheFold() iter.Seq[strom.Element] {
	return func(yield func(strom.Element) bool) {
		for ii := range 255 {
			iistr := strconv.Itoa(ii)
			if !yield(strom.CreateText("div", "Node "+iistr).SetStyle(map[string]string{
				"color":  vars.Color(vars.ColorRed),
				"height": vars.Size(vars.SizeLarge),
			})) {
				return
			}
			if !yield(strom.Comment("span").
				SetTextContent("Test Text")) {
				return
			}
		}
	}
}

func belowTheFold() iter.Seq[strom.Element] {

	return func(yield func(strom.Element) bool) {
		for ii := range 10000 {

			iistr := strconv.Itoa(ii)

			if !yield(strom.CreateText("div", "Node "+iistr).
				AddClass("test-deferred-class").
				SetAttribute("id", iistr).
				SetTextContent("Deferred Node " + iistr)) {
				return
			}
		}
	}
}
