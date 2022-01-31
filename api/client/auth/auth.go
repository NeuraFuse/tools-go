package auth

import (
	"github.com/neurafuse/tools-go/config"
	"github.com/neurafuse/tools-go/crypto/jwt"
	"github.com/neurafuse/tools-go/logging"
	usersID "github.com/neurafuse/tools-go/users/id"
	"github.com/neurafuse/tools-go/vars"
)

type F struct{}

var jwtAuthStatus bool
var createAuthStatusLast string

func (f F) Check(jwtAuth bool) {
	jwtAuthStatus = jwtAuth
	var status string
	if jwtAuth {
		jwt.SigningKeyActive = []byte(config.Setting("get", "user", "Spec.Auth.JWT.SigningKey", ""))
		status = "active"
	} else {
		jwt.ResetSigningKey()
		status = "default"
	}
	if status != createAuthStatusLast {
		if status == "active" {
			logging.Log([]string{"", vars.EmojiAPI, vars.EmojiUser}, "Logging in as user "+usersID.F.GetActive(usersID.F{})+"..", 0)
		}
		createAuthStatusLast = status
	}
}

func (f F) GetJwtAuthStatus() bool {
	return jwtAuthStatus
}
