package security

import (
	"net/http"

	"github.com/pkg/errors"
)

type Authenticator struct {
	DAO UserDAO
}

func (a Authenticator) AuthenticateRequest(r *http.Request) (string, error) {
	user, pass, ok := r.BasicAuth()
	if !ok {
		return "", errors.New("Could not parse 'Authorization' header")
	}

	if "" == user {
		return "", errors.New("User not specified")
	}

	authenticated, err := a.checkUser(user, pass)
	if nil != err {
		if errNotExisting == err {
			return "", nil
		}
		return "", errors.Wrap(err, "Error while attempting to authenticate "+user)
	}
	if !authenticated {
		return "", errors.New("User/Password don't correspond to each others")
	}
	return user, nil
}

var errNotExisting = errors.New("User doesn't exist")

func (a Authenticator) checkUser(username, password string) (bool, error) {
	usr, err := a.DAO.Load(username)
	if nil != err {
		return false, err
	}
	if nil == usr {
		return false, errNotExisting
	}

	return usr.Check(password), nil
}
