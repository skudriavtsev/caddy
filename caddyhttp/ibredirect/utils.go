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

func splitHostPort(hostport string) (string, string, error) {
	i := strings.LastIndexByte(hostport, ':')
	j := strings.LastIndexByte(hostport, ']')
	if i == -1 || j != -1 && j > i {
		hostport += ":"
	}

	return net.SplitHostPort(hostport)
}

const unsuffixDomainError = "Wrong domain name to remove the suffix from."
func GetUnSuffixedDomain(domain, suffix string) (string, error) {
	res, err := hasSuffix(domain, suffix)
	if err != nil {
		return "", err
	}
	if !res {
		return "", fmt.Errorf(unsuffixDomainError)
	}
	i := strings.LastIndex(domain, suffix)
	// len("a.b.suffix") must be greater than len("a.b.")
	if i < 4 {
		return "", fmt.Errorf(unsuffixDomainError)
	}

	return domain[:(i-1)], nil
}

func (rd Redirect) getUnSuffixedDomain(domain string) (string, error) {
	return GetUnSuffixedDomain(domain, rd.Suffix)
}

const unsuffixURLerror = "Wrong domain name to remove the suffix from."
func (rd Redirect) getUnSuffixedURL(rawURL string) (string, error) {
	url, err := urllib.Parse(rawURL)
	if err != nil {
		return "", err
	}

	sHost, _, err := splitHostPort(url.Host)
	if err != nil {
		return "", err
	}

	tHost, err := rd.getUnSuffixedDomain(sHost)
	if err != nil {
		return "", err
	}

	url.Host = tHost

	return url.String(), nil
}

func hasSuffix(h, s string) (bool, error) {
	host, _, err := splitHostPort(h)
	if err != nil {
		return false, err
	}
	if host == "" {
		return false, fmt.Errorf("'Host' header is wrong or not found")
	}
	if host == s || strings.HasSuffix(host, "." + s) {
		return true, nil
	}

	return false, nil
}

func (rd Redirect) hasSuffix(h string) (bool, error) {
	return hasSuffix(h, rd.Suffix)
}

// Takes original URL and redirects to a prefixed domain.
func (rd Redirect) redirectToSuffixedURL(w http.ResponseWriter, r *http.Request, origURL string) (int, error) {
	url, err := urllib.Parse(origURL)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	host,_ , err := splitHostPort(url.Host)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	url.Host = fmt.Sprintf("%s.%s", host, rd.Suffix)
	return rd.redirect(w, r, url.String())
}

func (rd Redirect) redirect(w http.ResponseWriter, r *http.Request, newURL string) (int, error) {
	http.Redirect(w, r, newURL, http.StatusSeeOther)

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
