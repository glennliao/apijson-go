package consts

const (
	ErrCode = 40001
)

type Err struct {
	code    int
	message string
}

func (e *Err) Code() int {
	return e.code
}

func (e *Err) Error() string {
	return e.message
}
