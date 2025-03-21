package types


type User struct {
	ID         int    `json:"id" db:"user_id"`
	Email      string `json:"email" db:"email"`
	Password   string `json:"password" db:"password"`
	FirstName  string `json:"first_name" db:"first_name"`
	LastName   string `json:"last_name" db:"last_name"`
	LoggedInAt string `json:"logged_in_at" db:"logged_in_at"`
	CreatedAt  string `json:"created_at" db:"created_at"`
}

type SignInRequest struct {
	ID       int    `json:"id" db:"user_id"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}
