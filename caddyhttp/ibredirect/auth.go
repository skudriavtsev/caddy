package ibredirect

import (
	"fmt"
	"net/http"
	"time"
	// "html/template"

	// log "github.com/sirupsen/logrus"
)

type authPageData struct {
	OriginURL string
	Message string
}

type authReqData struct {
	OriginURL string
	Token string
}

type authItem struct {
	Domains []string
	Expire time.Time
}

var (
	aStorage = make(map[string]*authItem)
	fakeAuthTokens = make(map[string]*authItem)
)

func getAuthReqData(r *http.Request) (*authReqData, error)  {
	aReqData := &authReqData{}
	token, found := r.Form[tokenParamName]
	if found {
		aReqData.Token = token[0]
	}

	origin, found := r.Form[originUrlParamName]
	if found {
		aReqData.OriginURL = origin[0]
	}

	return aReqData, nil
}

func getReferer(r *http.Request, ard *authReqData) string {
	if ard.OriginURL != "" {
		return ard.OriginURL
	}
	referer, found := r.Header["Referer"]
	if found {
		return referer[0]
	}
	return noRefererFound
}

// for now it is a fake auth
func authenticate(token string) (*authItem, error) {
	if len(token) == 5 {
		// Generate a fake error for tokens of 5 symbols.
		return nil, fmt.Errorf("Cannot verify the token: network error")
	}

	aItem, found := fakeAuthTokens[token]
	if !found {
		return nil, nil
	}

	return aItem, nil
}

func (rd Redirect) hasAuth(r *http.Request) bool {
	c, err := r.Cookie(authTokenCookieName)
	if err != nil {
		return false
	}
	aItem, found := aStorage[c.Value] // .Value is actually an auth token
	if !found {
		return false
	}

	domain, _, err := splitHostPort(r.Host)
	if err != nil {
		return false
	}

	uDomain, err := rd.getUnSuffixedDomain(domain)
	for _, d := range aItem.Domains {
		if uDomain == d {
			return true
		}
	}

	return false
}

func (rd Redirect) authPage(w http.ResponseWriter, r *http.Request) (int, error) {
	aReqData, err := getAuthReqData(r)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	aPageData := &authPageData{}
	aPageData.OriginURL = getReferer(r, aReqData)
	token := aReqData.Token
	if token == "" {
		// apply template
		// show page and stop processing
		return 0, nil
	}

	aItem, err := authenticate(token)
	if err != nil {
		aPageData.Message = err.Error()
		// apply template
		// show page and stop processing
		return 0, nil
	}

	if aItem == nil {
		aPageData.Message = "An invalid token has been entered"
		// apply template
		// show page and stop processing
		return 0, nil
	}

	aStorage[token] = aItem

	aCookie := http.Cookie{
		Name: authTokenCookieName,
		Value: token,
	}
	http.SetCookie(w, &aCookie)

	return rd.redirectToSuffixedURL(w, r)
}
