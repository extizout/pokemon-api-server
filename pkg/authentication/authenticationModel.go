package authentication

type (
	UserLoginRequestDto struct {
		Username string `json:"username" form:"username" validate:"required"`
		Password string `json:"password" form:"password" validate:"required"`
	}
	UserLoginResponseDto struct {
		AccessToken string `json:"accessToken"`
	}
	Claims struct {
		Username string `json:"username"`
	}
)
