package metallb

import (
	"context"
	"net"
	"sync"
	"time"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

// Metallb for k8s
type Metallb struct {
	Next plugin.Handler
	mu   sync.RWMutex
	Map  map[string]net.IP
}

// NewMetallb returns new object
func NewMetallb(next plugin.Handler) *Metallb {
	plug := &Metallb{
		Next: next,
		Map:  map[string]net.IP{},
	}
	go plug.Run()
	return plug
}

// ServeDNS does all the job
func (h *Metallb) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}
	log.Infof("log name %s", state.Name())
	if state.QType() != dns.TypeA {
		return plugin.NextOrFailure(h.Name(), h.Next, ctx, w, r)
	}
	ip, ok := h.Map[state.Name()]
	if !ok {
		return plugin.NextOrFailure(h.Name(), h.Next, ctx, w, r)
	}

	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative = true

	hdr := dns.RR_Header{
		Name:   state.QName(),
		Rrtype: state.QType(),
		Class:  state.QClass(),
		Ttl:    0,
	}
	m.Answer = []dns.RR{&dns.A{
		Hdr: hdr,
		A:   ip,
	}}
	state.SizeAndDo(m)
	w.WriteMsg(m)
	return 0, nil
}

// Run prepare
func (h *Metallb) Run() {
	h.Add("stamm.", "9.9.9.9")
	h.Add("redis.stamm.", "9.9.9.10")
	<-time.After(5 * time.Second)
	h.Add("new.stamm.", "9.9.9.11")
	<-time.After(5 * time.Second)
	h.Delete("new.stamm.")
}

// Add ip
func (h *Metallb) Add(name, ip string) {
	ipv4 := net.ParseIP(ip).To4()
	log.Infof("add ip %s for name %s", ipv4, name)
	h.mu.Lock()
	h.Map[name] = ipv4
	h.mu.Unlock()
}

// Delete ip
func (h *Metallb) Delete(name string) {
	log.Infof("delete name %s", name)
	h.mu.Lock()
	delete(h.Map, name)
	h.mu.Unlock()
}

// Name returns name of plugins
func (h *Metallb) Name() string {
	return "metallb"
}
