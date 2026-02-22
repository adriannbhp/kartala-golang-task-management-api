package users

// UpdateProfileParameter holds user profile update data
type UpdateProfileParameter struct {
	Username string `json:"username" binding:"omitempty,min=3,max=50" example:"johndoe"`
	Email    string `json:"email" binding:"omitempty,email" example:"john@example.com"`
}
