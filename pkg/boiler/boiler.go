package boiler

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

const (
	NoAuth AuthType = iota
	BasicAuth
)

type AuthType int

func (a AuthType) String() string {
	return [...]string{"No Auth", "Basic Auth"}[a]

}
func getCookie(c echo.Context) (*http.Cookie, bool) {
	cookie, err := c.Cookie(cookieName)
	if err != nil {
		return nil, false
	}
	return cookie, true
}

func setCookie(c echo.Context, id string) {
	cookie := http.Cookie{
		Value:   id,
		Name:    cookieName,
		Expires: time.Now().Add(24 * time.Hour),
	}
	c.SetCookie(&cookie)

}

type Auth struct {
	AuthType AuthType
	Username string
	Password string
}

type GetParams struct {
	Auth Auth
	Url  string
}

func (gp *GetParams) SetUrl(url string) *GetParams {
	gp.Url = url

	return gp
}

func (gp *GetParams) SetAuth(auth Auth) *GetParams {
	gp.Auth = auth
	return gp
}

func (gp *GetParams) SetAuthType(authType AuthType) *GetParams {
	gp.Auth.AuthType = authType
	return gp
}

func (gp *GetParams) SetUserName(username string) *GetParams {
	gp.Auth.Username = username
	return gp
}

func (gp *GetParams) SetPassword(password string) *GetParams {
	gp.Auth.Password = password
	return gp
}

func NewGetParams() *GetParams {
	return &GetParams{}
}

func GetJsonFromHttp(params *GetParams) ([]byte, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", params.Url, nil)
	if err != nil {
		fmt.Println("Error creating request: ", err)
		return nil, err
	}

	switch params.Auth.AuthType {
	case BasicAuth:
		fmt.Println("setting basic auth")
		req.SetBasicAuth(params.Auth.Username, params.Auth.Password)
	case NoAuth:
		// don't do anything
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if resp.StatusCode != 200 {

		statusError := fmt.Sprintf("non success status returned: %d", resp.StatusCode)
		fmt.Println(statusError)
		return nil, fmt.Errorf("%s", statusError)
	}
	if err != nil {
		fmt.Println("Error making request: ", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response: ", err)
		return nil, err
	}

	return body, nil

}
