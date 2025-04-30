package user

import "time"

type (
	UserProfile struct {
		Id        int       `json:"id"`
		Username  string    `json:"username"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	}

	CreateUserRequestDto struct {
		Username string `json:"username" form:"username" validate:"required"`
		Password string `json:"password" form:"password" validate:"required"`
	}

	CreateUserResponseDto struct {
		Id int `json:"id"`
	}
)
