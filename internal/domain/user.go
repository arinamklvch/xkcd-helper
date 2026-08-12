package domain

const (
	UsersRoleAdmin = "admin"
	UsersRoleUser  = "user"
)

type User struct {
	Login    string
	Password string
	Role     string
}
