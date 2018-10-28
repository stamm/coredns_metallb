package metallb

import (
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/mholt/caddy"
)

var log = clog.NewWithPlugin("root")

func init() {
	caddy.RegisterPlugin("metallb", caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	err := metallbParse(c)

	if err != nil {
		return plugin.Error("metallb", err)
	}
	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return NewMetallb(next)
	})

	return nil
}

func metallbParse(c *caddy.Controller) error {
	return nil
	// for c.Next() {
	// 	if !c.NextArg() {
	// 		return plugin.Error("metallb", c.ArgErr())
	// 	}
	// }
	// return nil
}
