package lib

import "errors"

var KeyInvalidError = errors.New("Key Invalid")
var RequestFailedError = errors.New("Request failed")
var JSONDecodeError = errors.New("Failed to decode JSON")
