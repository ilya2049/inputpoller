package stdinput

import (
	"errors"
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

var errTest = errors.New("test error")

func (p *Provider) ScanInput() (string, error) {
	var input string

	_, err := fmt.Scan(&input)

	p.mu.Lock()
	defer p.mu.Unlock()

	if p.stopped {
		return "", io.EOF
	}

	if err != nil {
		return "", err
	}

	if input == "err" {
		return "", errTest
	}

	return input, nil
}
