package outletsecurity

// User type represents logged in user
type User struct {
	Username string
	Password string
}

// CheckUser to authenticate User
func CheckUser(user User) bool {
	if user.Username == "user" && user.Password == "password" {
		return true
	}

	return false
}
