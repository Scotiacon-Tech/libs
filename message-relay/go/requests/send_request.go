package requests

type SendRequest struct {
	From    any    `json:"from"`
	To      any    `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}
