package dto

type OAuthResponse struct {
	Body    string            `json:"body"`
	Path    string            `json:"path"`
	Headers map[string]string `json:"headers"`
}
