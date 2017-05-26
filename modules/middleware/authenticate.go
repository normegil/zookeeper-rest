package middleware

import (
	"context"
	"net/http"

	"github.com/pkg/errors"

	"github.com/Sirupsen/logrus"
	"github.com/normegil/zookeeper-rest/modules/database/mongo"
	errPkg "github.com/normegil/zookeeper-rest/modules/errors"
	"github.com/normegil/zookeeper-rest/modules/security"
)

const USER_CTX_KEY = "user"
const REQUEST_HEADER_AUTHORIZATION = "Authorization"

func RequestAuthenticator(log *logrus.Entry, mongoDB *mongo.Mongo, h http.Handler) http.Handler {

	dao := &mongo.MongoUserDAO{mongoDB}
	authenticator := security.Authenticator{dao}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := authenticator.AuthenticateRequest(r)
		if nil != err {
			errPkg.Handler{log}.Handle(w, errPkg.NewErrWithCode(40101, errors.Wrap(err, "Unable to authenticate user with 'Authentication' header content")))
		} else if "" == user {
			errPkg.Handler{log}.Handle(w, errPkg.NewErrWithCode(40102, errors.New("Unable to authenticate user with 'Authentication' header content (No error but user is empty)")))
		} else {
			h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), USER_CTX_KEY, user)))
		}
	})
}
