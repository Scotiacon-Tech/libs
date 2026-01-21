package requests

type TokenRequest struct {
	GrantType    string   `json:"grant_type"`
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	Audience     []string `json:"audience"`
	Scope        string   `json:"scope"`
}
