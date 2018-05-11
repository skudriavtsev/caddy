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
		if !rd.hasAuth(r) {
			return rd.authPage(w, r)
		}
		return rd.Next.ServeHTTP(w, r)
	}

	target := fmt.Sprintf("http://%s%s", r.Host, r.RequestURI)
	return rd.redirectToSuffixedURL(w, r, target)
}
