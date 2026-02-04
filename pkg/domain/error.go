package domain

import "fmt"

var (
	ErrNotAuth  = &PdaccessErr{code: "PDA-ERROR-1000", message: "Not Authentication"}
	ErrInternal = &PdaccessErr{code: "PDA-ERROR-9999", message: "Unknown Error"}
)

type PdaccessErr struct {
	code    string
	message string
}

func (t *PdaccessErr) Error() string {
	return fmt.Sprintf("%s: %s", t.code, t.message)
}

func (t *PdaccessErr) Code() string {
	return t.code
}

type ApiResponse struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestId string `json:"requestId"`
}
