package models

type User struct {
	ID              int64  `json:"id"`
	FirstName       string `json:"first_name" binding:"required,min=2,max=30"`
	LastName        string `json:"last_name" binding:"required,min=2,max=30"`
	PhoneNumber     string `json:"phone_number"`
	Email           string `json:"email" binding:"required,email"`
	Gender          string `json:"gender" binding:"oneof=male female"`
	Username        string `json:"username"`
	ProfileImageUrl string `json:"profile_image_url"`
	Type            string `json:"type" binding:"required,oneof=user superadmin"`
	Password        string `json:"password" binding:"required,min=6,max=16"`
	CreatedAt       string `json:"created_at"`
}

type CreateUserRequest struct {
	FirstName       string `json:"first_name" binding:"required,min=2,max=30"`
	LastName        string `json:"last_name" binding:"required,min=2,max=30"`
	PhoneNumber     string `json:"phone_number"`
	Email           string `json:"email" binding:"required,email"`
	Gender          string `json:"gender" binding:"oneof=male female"`
	Username        string `json:"username"`
	ProfileImageUrl string `json:"profile_image_url"`
	Type            string `json:"type" binding:"required,oneof=user superadmin"`
	Password        string `json:"password" binding:"required,min=6,max=16"`
}

type UpdateUserRequest struct {
	FirstName       string `json:"first_name" binding:"required,min=2,max=30"`
	LastName        string `json:"last_name" binding:"required,min=2,max=30"`
	PhoneNumber     string `json:"phone_number"`
	Gender          string `json:"gender" binding:"oneof=male female"`
	Username        string `json:"username"`
	ProfileImageUrl string `json:"profile_image_url"`
	Type            string `json:"type" binding:"required,oneof=superadmin user"`
}

type GetAllUsersResponse struct {
	Users []User `json:"users"`
	Count int32  `json:"count"`
}

type GetAllUserParams struct {
	Limit  int32  `json:"limit" binding:"required" default:"10"`
	Page   int32  `json:"page" binding:"required" default:"1"`
	Search string `json:"search"`
}

type GetAllParams struct {
	Limit  int32  `json:"limit" binding:"required" default:"10"`
	Page   int32  `json:"page" binding:"required" default:"1"`
	Search string `json:"search"`
}

