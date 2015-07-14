package lenz

import (
	"os"

	"github.com/HunanTV/eru-agent/common"
	"github.com/HunanTV/eru-agent/defines"
	"github.com/HunanTV/eru-agent/logs"
)

var Lenz *LenzForwarder

type LenzForwarder struct {
	Attacher *AttachManager
	Router   *RouteManager
	Routefs  RouteFileStore
}

func NewLenz(config defines.LenzConfig) *LenzForwarder {
	obj := &LenzForwarder{}
	obj.Attacher = NewAttachManager(common.Docker)
	obj.Router = NewRouteManager(obj.Attacher, config.Stdout)
	obj.Routefs = RouteFileStore(config.Routes)

	if len(config.Forwards) > 0 {
		logs.Info("Routing all to", config.Forwards)
		target := defines.Target{Addrs: config.Forwards}
		route := defines.Route{ID: "lenz_default", Target: &target}
		route.LoadBackends()
		obj.Router.Add(&route)
	}

	if _, err := os.Stat(config.Routes); err == nil {
		logs.Info("Loading and persisting routes in", config.Routes)
		logs.Assert(obj.Router.Load(obj.Routefs), "persistor")
	}
	return obj
}
