package ibredirect

import (
	"fmt"
	"net"
	"net/http"
	urllib "net/url"
	"strings"

	// log "github.com/sirupsen/logrus"
)

const (
	noRefererFound = "noRefererFound"
)

func splitHostPort(hostport string) (string, string, error) {
	i := strings.LastIndexByte(hostport, ':')
	j := strings.LastIndexByte(hostport, ']')
	if i == -1 || j != -1 && j > i {
		hostport += ":"
	}

	return net.SplitHostPort(hostport)
}

func getDefaultPort(r *http.Request) string {
	if r.TLS == nil {
		return "80"
	}

	return "443"
}


const unsuffixDomainError = "Wrong domain name to remove the suffix from."
func (rd Redirect) getUnSuffixedDomain(domain string) (string, error) {
	res, err := rd.hasSuffix(domain)
	if err != nil {
		return "", err
	}
	if !res {
		return "", fmt.Errorf(unsuffixDomainError)
	}
	i := strings.LastIndex(domain, rd.Suffix)
	// len("a.b.suffix") must be greater than len("a.b.")
	if i < 4 {
		return "", fmt.Errorf(unsuffixDomainError)
	}

	return domain[:i], nil
}

const unsuffixURLerror = "Wrong domain name to remove the suffix from."
func (rd Redirect) getUnSuffixedURL(rawURL string) (string, error) {
	url, err := urllib.Parse(rawURL)
	if err != nil {
		return "", err
	}

	sHost, port, err := splitHostPort(url.Host)
	if err != nil {
		return "", err
	}

	tHost, err := rd.getUnSuffixedDomain(sHost)
	if err != nil {
		return "", err
	}

	if port != "" {
		tHost = net.JoinHostPort(tHost, port)
	}
	url.Host = tHost

	return url.String(), nil
}

func (rd Redirect) hasSuffix(h string) (bool, error) {
	host, _, err := splitHostPort(h)
	if err != nil {
		return false, err
	}
	if host == "" {
		return false, fmt.Errorf("'Host' header is wrong or not found")
	}
	if host == rd.Suffix || strings.HasSuffix(host, "." + rd.Suffix) {
		return true, nil
	}

	return false, nil
}

func (rd Redirect) redirectToSuffixedURL(w http.ResponseWriter, r *http.Request) (int, error) {
	var secure string

	host, port, err := splitHostPort(r.Host)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if port == "" {
		port = getDefaultPort(r)
	}
	if r.TLS != nil {
		secure = "s"
	}
	target := fmt.Sprintf("http%s://%s.%s:%s/%s", secure, host, rd.Suffix, port, r.RequestURI)

	return rd.redirectWithReferer(w, r, target)
}

const redirectPage = `
<html xmlns="http://www.w3.org/1999/xhtml">    
  <head>      
    <meta http-equiv="refresh" content="0;URL='%s'" />    
  </head>    
  <body>Redirecting, please wait...</body>  
</html>     
`
func (rd Redirect) redirectWithReferer(w http.ResponseWriter, r *http.Request, newURL string) (int, error) {
	fmt.Fprintf(w, redirectPage, newURL)
	return 0, nil
}
