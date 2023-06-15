package FlowX

import (
	"fmt"
	"net/http"
	"strings"
)

func GetCookiesFromResponse(resp *http.Response, cookieNames []string) ([]*http.Cookie, error) {
	cookies := make([]*http.Cookie, 0, len(cookieNames))

	for _, cookieName := range cookieNames {
		found := false
		for _, header := range resp.Header.Values("Set-Cookie") {
			cookie := strings.SplitN(header, ";", 2)[0]
			if strings.HasPrefix(cookie, cookieName+"=") {
				value := strings.TrimPrefix(cookie, cookieName+"=")
				cookies = append(cookies, &http.Cookie{Name: cookieName, Value: value})
				found = true
				break
			}
		}

		if !found {
			return nil, fmt.Errorf("cookie not found: %s", cookieName)
		}
	}

	return cookies, nil
}

func GetCookieFromResponse(resp *http.Response, cookieName string) (*http.Cookie, error) {
	for _, header := range resp.Header.Values("Set-Cookie") {
		cookie := strings.SplitN(header, ";", 2)[0]
		if strings.HasPrefix(cookie, cookieName+"=") {
			value := strings.TrimPrefix(cookie, cookieName+"=")
			return &http.Cookie{Name: cookieName, Value: value}, nil
		}
	}

	return nil, fmt.Errorf("cookie not found: %s", cookieName)
}

func SetCookieForRequest(req *http.Request, cookie *http.Cookie) {
	req.AddCookie(cookie)
}

func SetCookiesForRequest(req *http.Request, cookies []*http.Cookie) {
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
}

