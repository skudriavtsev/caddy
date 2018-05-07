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

	return rd.redirectWithReferer(w, r, target)
}

const stubPage = `
<html xmlns="http://www.w3.org/1999/xhtml">    
  <head>      
    <title>Test</title>    
  </head>    
  <body>Debug</body>  
</html>     
`
func (rd Redirect) serve(w http.ResponseWriter, r *http.Request) (int, error) {
	if !rd.hasAuth(r) {
		return rd.authPage(w, r)
	}

	// the main processing
	fmt.Fprintf(w, stubPage)

	return 0, nil
}
