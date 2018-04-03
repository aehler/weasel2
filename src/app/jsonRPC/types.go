package jsonRPC

type Params interface {
	IsValid() bool
	Url() string
	RequiresAuth() bool
	GetToken() string
}

type Client interface {
	Send(params Params, result interface{}) error
}

type client struct {
	name    string
	url     string
	address string
}

type Result interface{}

type Error error

type Response struct {
	Result Result `json:"r,omitempty"`
	Error  *Error `json:"e,omitempty"`
}