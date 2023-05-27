package reverse

type Action byte

const (
	Accept Action = 1 + iota
	Dial
)

type Request struct {
	Action  Action
	Network string
	Address string
}

type Response struct {
	Succeed bool
	Message error
}

func FailResponse(err error) Response {
	return Response{Succeed: false, Message: err}
}

var (
	SucceedResponse = Response{
		Succeed: true,
	}
)
