package dtos

type BaseResponse struct {
	Meta Meta `json:"meta"`
}

// Meta contains metadata that response in each request to client
type Meta struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Limit   int64  `json:"limit,omitempty"`
	Offset  int64  `json:"offset,omitempty"`
	Total   int64  `json:"total,omitempty"`
}

func (m *Meta) SetLimit(limit int64) *Meta {
	m.Limit = limit
	return m
}

func (m *Meta) SetOffset(offset int64) *Meta {
	m.Offset = offset
	return m
}

func (m *Meta) SetTotal(total int64) *Meta {
	m.Total = total
	return m
}

const (
	Success       = "Success"
	BindError     = "BindError"
	InternalError = "InternalError"
	UserExist     = "UserExist"
	UserNotExist  = "UserNotExist"
	PasswordWrong = "PasswordWrong"
	TokenRevoke   = "TokenRevoke"
)

func GetMeta(metaType string) Meta {
	meta, ok := mapMeta[metaType]
	if ok {
		return meta
	}
	return internalError
}

var mapMeta = map[string]Meta{
	Success:       success,
	BindError:     bindError,
	InternalError: internalError,
	UserExist:     userExist,
	UserNotExist:  userNotExist,
	PasswordWrong: passwordWrong,
	TokenRevoke:   tokenRevoke,
}

var success = Meta{
	Code:    "200",
	Message: "Success",
}

var bindError = Meta{
	Code:    "4000100",
	Message: "Bind error",
}

var tokenRevoke = Meta{
	Code:    "4010100",
	Message: "Token revoke",
}

var internalError = Meta{
	Code:    "5000100",
	Message: "Internal Error",
}

var userExist = Meta{
	Code:    "4000101",
	Message: "User name exist",
}

var userNotExist = Meta{
	Code:    "4000102",
	Message: "User not exist",
}

var passwordWrong = Meta{
	Code:    "4000102",
	Message: "Password Wrong",
}
