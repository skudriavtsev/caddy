package ibredirect

import (
	"fmt"
	"net/http"
	"time"

	// log "github.com/sirupsen/logrus"
)

type authPageData struct {
	ActionURL string
	Message string
}

type authReqData struct {
	Token string
}

type authItem struct {
	Domains []string
	Expires time.Time
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

	return aReqData, nil
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

	if time.Now().After(aItem.Expires) {
		return false
	}

	domain, _, err := splitHostPort(r.Host)
	if err != nil {
		return false
	}

	uDomain, err := rd.getUnSuffixedDomain(domain)
	if err != nil {
		return false
	}

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

	aPageData := &authPageData{
		ActionURL: r.URL.String(),
	}
	token := aReqData.Token
	if token == "" {
		err = showAuthPage(w, aPageData)
		if err != nil {
			return http.StatusInternalServerError, err
		}
		return 0, nil
	}

	aItem, err := authenticate(token)
	if err != nil {
		aPageData.Message = err.Error()
		err = showAuthPage(w, aPageData)
		if err != nil {
			return http.StatusInternalServerError, err
		}
		return 0, nil
	}

	if aItem == nil {
		aPageData.Message = "An invalid token has been entered"
		err = showAuthPage(w, aPageData)
		if err != nil {
			return http.StatusInternalServerError, err
		}
		return 0, nil
	}

	aStorage[token] = aItem

	aCookie := http.Cookie{
		Path: "/",
		Domain: "." + rd.Suffix,
		Name: authTokenCookieName,
		Value: token,
		Expires: aItem.Expires,
	}
	http.SetCookie(w, &aCookie)

	target := fmt.Sprintf("http://%s%s", r.Host, r.RequestURI)
	return rd.redirect(w, r, target)
}
