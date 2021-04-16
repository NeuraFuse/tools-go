package auth

import (
	"../../../config"
	"../../../crypto/jwt"
	"../../../logging"
	"../../../vars"
	"../../../users"
)

type F struct{}

var jwtAuthStatus bool
var createAuthStatusLast string = ""

func (f F) Check(jwtAuth bool) {
	jwtAuthStatus = jwtAuth
	status := ""
	if jwtAuth {
		jwt.SigningKeyActive = []byte(config.Setting("get", "user", "Spec.Auth.JWT.SigningKey", ""))
		status = "active"
	} else {
		jwt.SigningKeyActive = jwt.SigningKeyDefault
		status = "default"
	}
	if status != createAuthStatusLast {
		if status == "active" {
			logging.Log([]string{"", vars.EmojiAPI, vars.EmojiUser}, "Logging in as user "+users.GetIDActive()+"..", 0)
		}
		createAuthStatusLast = status
	}
}

func (f F) GetJwtAuthStatus() bool {
	return jwtAuthStatus
}
