package routers

import (
	"go-learning/internal/routers/manage"
	"go-learning/internal/routers/user"
)

type RouterGroup struct {
	Manage manage.ManageRouterGroup
	User   user.UserRouterGroup
}

var RouterGroupApp = new(RouterGroup)
