package strom

import (
	"errors"
	"net/http"
)

func WriteResponse(w http.ResponseWriter, root Element) error {

	if root.GetTagName() != htmlTagName {
		return errors.New("writing response requires html element")
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return root.Write(w)

}
