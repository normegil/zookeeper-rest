package security

type User struct {
	Username     string `json:"name" bson:"name"`
	Pass string `json:"password" bson:"password"`
}

func (u User) Name() string {
	return u.Username
}

func (u User) Password() string {
	return u.Pass
}

func (u User) Check(password string) bool {
	return u.Password() == password
}
