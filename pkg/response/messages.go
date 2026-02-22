package response

const (
	// Success Messages
	SuccessLogin    = "Login successful"
	SuccessRegister = "Registration successful"
	SuccessFetch    = "Data retrieved successfully"
	SuccessUpdate   = "Data updated successfully"
	SuccessDelete   = "Data deleted successfully"
	SuccessRefresh  = "Session refreshed successfully"
	SuccessCreated  = "Data created successfully"

	// Error Messages
	ErrInvalidRequest     = "Invalid request parameters"
	ErrInvalidCredentials = "Invalid email or password"
	ErrEmailExists        = "Email is already registered"
	ErrUsernameExists     = "Username is already taken"
	ErrInternalServer     = "Internal server error, please try again later"
	ErrUnauthorized       = "Access denied, please login first"
	ErrForbidden          = "Forbidden access"
	ErrUserNotFound       = "User not found"
	ErrTaskNotFound       = "Task not found"
	ErrInvalidToken       = "Invalid or expired session"
	ErrResourceNotFound   = "Resource not found"
	ErrWeakPassword       = "Password must be at least 8 characters"
	ErrEmailNotFound      = "Email not found"
	ErrInvalidPassword    = "Incorrect password"
	ErrInternalDatabase   = "Database error occurred, please contact administrator"
	ErrUnauthorizedAccess = "You do not have access to this resource"
)
