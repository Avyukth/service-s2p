package usergrp

import (
	userCore "github.com/Avyukth/service3-clone/business/core/user"
	"github.com/Avyukth/service3-clone/business/sys/auth"
)

type Handlers struct {
	User userCore.Core
	Auth *auth.Auth
}
