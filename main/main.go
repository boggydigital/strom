package main

import (
	"fmt"
	"iter"
	"net/http"
	"strconv"
	"time"

	"github.com/boggydigital/strom"
)

func main() {
	http.HandleFunc("GET /test", GetTest)
	if err := http.ListenAndServe(":12345", nil); err != nil {
		panic(err)
	}
}

func belowTheFold() iter.Seq[strom.Element] {

	return func(yield func(strom.Element) bool) {
		for ii := range 100000 {

			iistr := strconv.Itoa(ii)

			if !yield(strom.CreateElement("div", "Node "+iistr).
				AddClass("test-deferred-class").
				SetAttribute("id", iistr).
				SetTextContent("Deferred Node " + iistr)) {
				return
			}
		}
	}
}

func GetTest(w http.ResponseWriter, r *http.Request) {

	start := time.Now()

	root := strom.CreateElement("html")

	for ii := range 100 {
		iistr := strconv.Itoa(ii)
		root.Append(strom.CreateElement("div", "Node "+iistr).
			AddClass("test-class").
			SetAttribute("id", iistr),
			strom.CreateElement("span").
				SetTextContent("Test Text"))
	}

	root.Defer(belowTheFold)

	if err := root.Write(w); err != nil {
		panic(err)
	}

	elapsed := time.Since(start)

	fmt.Println(elapsed.Milliseconds(), "ms")
}
