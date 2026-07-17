package strom

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"iter"
	"net/http"
	"strings"
)

func WriteResponse(w http.ResponseWriter, root Element, scripts ...[]byte) error {

	if root.GetTagName() != htmlTagName {
		return errors.New("writing response requires html element")
	}

	var body Element
	for body = range root.GetElementsByTagName("body") {
		break
	}

	var getScripts = func() iter.Seq[Element] {
		return func(yield func(Element) bool) {
			for _, script := range scripts {
				if !yield(Script(script)) {
					return
				}
			}
		}
	}

	if body != nil {
		body.Append(OnDemand(getScripts))
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	csp, err := contentSecurityPolicy(scripts...)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Security-Policy", csp)

	return root.Write(w)
}

func contentSecurityPolicy(scripts ...[]byte) (string, error) {

	var scriptCsp string
	scriptDigests := make([]string, 0, len(scripts))

	for _, script := range scripts {

		sh, err := computeSha256(bytes.NewReader(script))
		if err != nil {
			return "", err
		}

		scriptDigests = append(scriptDigests, "'sha256-"+sh+"'")
	}

	if len(scriptDigests) > 0 {
		scriptCsp = "script-src " + strings.Join(scriptDigests, " ")
	}

	objectCsp := "object-src 'none'"
	frameAncestorsCsp := "frame-ancestors 'self'"
	baseUriCsp := "base-uri 'self'"
	formActionCsp := "form-action 'self'"

	return strings.Join([]string{scriptCsp, objectCsp, frameAncestorsCsp, baseUriCsp, formActionCsp}, "; "), nil
}

func computeSha256(reader io.Reader) (string, error) {
	h := sha256.New()
	var err error
	if _, err = io.Copy(h, reader); err == nil {
		hs := h.Sum(nil)
		return base64.StdEncoding.EncodeToString(hs), nil
	}
	return "", err
}
