package vo

type BookInfoVO struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
	Count uint   `json:"count"`
	ISBN  string `json:"isbn"`
}
