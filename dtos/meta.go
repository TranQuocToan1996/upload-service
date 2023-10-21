package dtos

type BaseResponse struct {
	Meta Meta `json:"meta"`
}

// Meta contains metadata that response in each request to client
type Meta struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

const (
	Success       = "Success"
	BindError     = "BindError"
	InternalError = "InternalError"
	UserExist     = "UserExist"
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
}

var success = Meta{
	Code:    "200",
	Message: "Success",
}

var bindError = Meta{
	Code:    "4000100",
	Message: "Bind error",
}

var internalError = Meta{
	Code:    "5000100",
	Message: "Internal Error",
}

var userExist = Meta{
	Code:    "4000101",
	Message: "User name exist",
}