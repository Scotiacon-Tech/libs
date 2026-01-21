package lib

// Yes "Result" here not "Response".
// This doesn't come back from a http request.

type SendResult struct {
	NewKey    string
	MessageID string
}
