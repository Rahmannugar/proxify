package loadbalancer

import "errors"

// ErrNoHealthyUpstreams indicates that no healthy upstreams are available to receive traffic.
var ErrNoHealthyUpstreams = errors.New("no healthy upstreams available")
