package xqCookie

import (
	"net/http"
)

func GetCookie(client *http.Client) (string, error) {
	req, err := http.NewRequest("GET", "https://xueqiu.com", nil)
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	cookie := resp.Cookies()[1]
	return cookie.Name + "=" + cookie.Value, nil
}
