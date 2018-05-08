package ibredirect

import (
	"fmt"
	"net"
	"net/http"
	urllib "net/url"
	"strings"
	"html/template"

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

	return domain[:(i-1)], nil
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

// Takes original URL and redirects to a prefixed domain.
func (rd Redirect) redirectToSuffixedURL(w http.ResponseWriter, ard *authReqData) (int, error) {
	url, err := urllib.Parse(ard.OriginURL)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	host, port, err := splitHostPort(url.Host)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if port == "" {
		if url.Scheme == "http" {
			port = "80"
		} else {
			port = "443"
		}
	}
	url.Host = fmt.Sprintf("%s.%s:%s", host, rd.Suffix, port)

	return rd.redirectWithReferer(w, url.String())
}

func (rd Redirect) redirectWithReferer(w http.ResponseWriter, newURL string) (int, error) {
	fmt.Fprintf(w, redirectPage, newURL)
	return 0, nil
}

func showAuthPage(w http.ResponseWriter, data *authPageData) error {
	tpl, err := template.New("authPage").Parse(authPageTemplate)
	if err != nil {
		return err
	}
	err = tpl.Execute(w, data)
	if err != nil {
		return err
	}
	return nil
}
