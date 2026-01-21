package requests

type SendRequest struct {
	To      any    `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}
