package models

// RegisterRequest struct
type RegisterRequest struct {
	Username     string `json:"username"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	FullName     string `json:"full_name"`
	MobileNumber string `json:"mobile_number"`
	BirthDate    string `json:"birth_date"`
	Gender       string `json:"gender"`
}
