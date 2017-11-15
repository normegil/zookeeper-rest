package middleware

import (
	"context"
	"net/http"

	"github.com/pkg/errors"

	"github.com/normegil/resterrors"
	definitions "github.com/normegil/zookeeper-rest/modules/errors"
	"github.com/normegil/zookeeper-rest/modules/model/dao"
	"github.com/sirupsen/logrus"
)

const USER_CTX_KEY = "user"
const REQUEST_HEADER_AUTHORIZATION = "Authorization"

func RequestAuthenticator(log logrus.FieldLogger, userDAO dao.UserDAO, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errHandler := resterrors.Handler{definitions.Definitions(), definitions.DEFAULT_ERROR_CODE}

		username, pass, ok := r.BasicAuth()
		if !ok {
			errHandler.Handle(w, resterrors.NewErrWithCode(40101, errors.New("Could not parse 'Authorization' header")))
			return
		}
		if "" == username {
			errHandler.Handle(w, resterrors.NewErrWithCode(40103, errors.New("User not specified")))
			return
		}

		user, err := userDAO.Load(username)
		if err != nil {
			errHandler.Handle(w, errors.Wrapf(err, "Loading %s", username))
			return
		}
		if nil == user {
			errHandler.Handle(w, resterrors.NewErrWithCode(40102, errors.New("User doesn't exist")))
			return
		}

		if user.Check(pass) {
			h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), USER_CTX_KEY, user)))
		} else {
			errHandler.Handle(w, resterrors.NewErrWithCode(40105, errors.New("User found but password doesn't match")))
		}
	})
}
