package types

type Unauthorized struct {
	Message string
}

func (e *Unauthorized) Error() string {
	if e.Message == "" {
		return "unauthorized"
	}
	return e.Message
}

type Forbidden struct {
	Message string
}

func (e *Forbidden) Error() string {
	if e.Message == "" {
		return "forbidden"
	}
	return e.Message
}

type NotFound struct {
	Message string
}

func (e *NotFound) Error() string {
	if e.Message == "" {
		return "not found"
	}
	return e.Message
}

type BadRequest struct {
	Message string
}

func (e *BadRequest) Error() string {
	if e.Message == "" {
		return "bad request"
	}
	return e.Message
}
