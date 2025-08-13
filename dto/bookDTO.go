package dto

type BookInfoDTO struct {
	Title string `json:"title" validate:"required"`
	Count uint   `json:"count" validate:"required"`
	ISBN  string `json:"isbn" validate:"required"`
}

type BookUpdateDTO struct {
	ID uint `json:"id" validate:"required"`
	BookInfoDTO
}

type BookSearchDTO struct {
	Title string `json:"title"`
	ISBN  string `json:"isbn"`
}
