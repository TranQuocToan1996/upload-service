package dtos

type RevokeTokenRequest struct {
	RevokeTokenAt int64 `json:"revoke_token_at" validate:"required,gt=0"`
}

type RevokeTokenResponse struct {
	Meta Meta `json:"meta"`
}
