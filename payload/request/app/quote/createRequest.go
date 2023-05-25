package faq

type CreateRequest struct {
	FullName    string `json:"fullname"`
	Email       string `json:"email"`
	ServiceId   string `json:"serviceId"`
	ServiceType string `json:"serviceType"`
	Phone       string `json:"phone"`
	Country     string `json:"country"`
}
