package faq

type CreateRequest struct {
	FullName string `json:"fullname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Country  string `json:"country"`
	Message  string `json:"message"`
}
