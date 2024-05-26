package response

import "encoding/json"

// ProblemDetail is a struct for problem detail
type ProblemDetail struct {
	Title     string `json:"title" example:"Invalid request body"`
	Status    int    `json:"status" example:"400"`
	Code      string `json:"code" example:"ERR-001"`
	Detail    string `json:"detail" example:"Content-Type header is missing"`
	Timestamp string `json:"timestamp" example:"2021-07-01T15:04:05.999999-07:00"`
	Instance  string `json:"instance" example:"/api/v1/author"`
}

func (p ProblemDetail) JSON() []byte {
	bytes, _ := json.Marshal(p) //nolint:errchkjson // i don't care about error here
	return bytes
}
