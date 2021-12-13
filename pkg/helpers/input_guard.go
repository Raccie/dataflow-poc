package helpers

// InputGuard counts number of closed inputs
type InputGuard struct {
	ports    map[string]bool
	complete int
}

// NewInputGuard returns a guard for a given number of inputs
func NewInputGuard(ports ...string) *InputGuard {
	portMap := make(map[string]bool, len(ports))
	for _, p := range ports {
		portMap[p] = false
	}
	return &InputGuard{portMap, 0}
}

// Complete is called when a port is closed and returns true when all the ports have been closed
func (g *InputGuard) Complete(port string) bool {
	if !g.ports[port] {
		g.ports[port] = true
		g.complete++
	}
	return g.complete >= len(g.ports)
}
