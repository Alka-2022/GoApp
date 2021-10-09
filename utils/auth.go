package utils

import (
	"encoding/json"
	"io"
	"net/url"
)


type OauthParams struct {
	ClientID    string `json:"clientId"`
	Code        string `json:"code"`
	RedirectURI string `json:"redirectUri"`
}

func MakeOauthQparams(body io.ReadCloser) (url.Values, OauthParams, error) {
	var oauthParams OauthParams

	
	defer body.Close()
	err := json.NewDecoder(body).Decode(&oauthParams)

	if err != nil {
		return nil, oauthParams, err
	}

	// make query params
	q := url.Values{}
	q.Set("client_id", oauthParams.ClientID)
	q.Add("redirect_uri", oauthParams.RedirectURI)
	q.Add("code", oauthParams.Code)
	q.Add("grant_type", "authorization_code")

	return q, oauthParams, nil
}
