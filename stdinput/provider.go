package stdinput

import (
	"fmt"
	"io"
	"sync"
)

type Provider struct {
	mu      sync.Mutex
	stopped bool
}

func (p *Provider) StopProviding() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.stopped = true
}

func (p *Provider) ScanInput() (string, error) {
	var input string

	if _, err := fmt.Scan(&input); err != nil {
		return "", err
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if p.stopped {
		return "", io.EOF
	}

	return input, nil
}
