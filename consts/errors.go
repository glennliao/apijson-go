package consts

// action

var (
	ErrNoTag = Err{
		code:    400,
		message: "no tag",
	}
)

func NewStructureKeyNoFoundErr(k string) Err {
	return NewValidStructureErr("field " + k + " is not found")
}

func NewValidStructureErr(msg string) Err {
	return Err{
		code:    400,
		message: msg,
	}
}

func NewValidReqErr(msg string) Err {
	return Err{
		code:    400,
		message: msg,
	}
}

func NewMethodNotSupportErr(msg string) Err {
	return Err{
		code:    400,
		message: msg,
	}
}

// access

func NewDenyErr(key, role string) Err {
	return Err{
		code:    403,
		message: "deny node: " + key + " with " + role,
	}
}

func NewNoAccessErr(key, role string) Err {
	return Err{
		code:    403,
		message: "node not access: " + key + " with " + role,
	}
}

func NewAccessNoFoundErr(key string) Err {
	return Err{
		code:    404,
		message: "access no found: " + key,
	}
}

func NewSysErr(msg string) Err {
	return Err{
		code:    500,
		message: msg,
	}
}

type Err struct {
	code    int
	message string
}

func (e Err) Code() int {
	return e.code
}

func (e Err) Error() string {
	return e.message
}
