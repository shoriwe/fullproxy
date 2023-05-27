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
	Message string
}

func FailResponse(err error) Response {
	return Response{Succeed: false, Message: err.Error()}
}

var (
	SucceedResponse = Response{
		Succeed: true,
		Message: "Succeed",
	}
)
