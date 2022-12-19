package models

type Like struct {
	Id     int  `json:"id"`
	PostId int  `json:"post_id"`
	UserId int  `json:"user_id"`
	Status bool `json:"status"`
}

type CreateOrUpdateLikeRequest struct {
	PostId int `json:"post_id" binding:"required"`
	Status bool  `json:"status"`
}
