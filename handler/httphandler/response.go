package httphandler

const (
	CodeOK          = "0000"
	DescOK          = "success"
	CodeBadReq      = "0001"
	DescBadReq      = "bad request"
	CodeInternalErr = "9999"
	DescInternalErr = "internal server error"
)

type Res[T any] struct {
	Code string `json:"code"`
	Desc string `json:"desc"`
	Data *T     `json:"data,omitempty"`
}
