package security

type UserDAO interface {
	Load(user string) (*User, error)
}
