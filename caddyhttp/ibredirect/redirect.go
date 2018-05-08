package ibredirect

import (
	"fmt"
	"net/http"

	// log "github.com/sirupsen/logrus"
)

// ServeHTTP implements the httpserver.Handler interface.
func (rd Redirect) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	err := r.ParseForm()
	if err != nil {
		return http.StatusInternalServerError, err
	}
	found, err := rd.hasSuffix(r.Host)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if found {
		return rd.serve(w, r)
	}

	target := fmt.Sprintf("http://%s/", rd.Suffix)

	return rd.redirectWithReferer(w, target)
}

func (rd Redirect) serve(w http.ResponseWriter, r *http.Request) (int, error) {
	if !rd.hasAuth(r) {
		return rd.authPage(w, r)
	}

	// the main processing
	fmt.Fprintf(w, stubPage)

	return 0, nil
}
