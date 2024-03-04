package response

const (
	StatusSuccess = "OK"
	StatusError   = "Error"
)

type Response struct {
	Status string      `json:"status"`
	Error  string      `json:"error,omitempty"`
	Data   interface{} `json:"data,omitempty"`
}

func (r Response) GetStatus() string {
	return r.Status
}

func (r Response) GetMessage() string {
	return r.Error
}

func (r Response) GetData() interface{} {
	return r.Data
}

func OK(data interface{}) Response {
	return Response{
		Status: StatusSuccess,
		Data:   data,
	}
}

func GetError(msg string, data interface{}) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
		Data:   data,
	}
}
