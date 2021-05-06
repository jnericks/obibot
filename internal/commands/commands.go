package commands

import (
	"context"
	"fmt"
)

type Input struct {
	Args []string
}

type Output struct {
	Response string
	Markdown bool
}

type Func func(context.Context, Input) (*Output, error)

type Manager interface {
	Register(string, Func) error
	CanExec(string) bool
	Exec(context.Context, string, Input) (*Output, error)
}

func NewManager() Manager {
	return &manager{
		registry: make(map[string]Func),
	}
}

type manager struct {
	registry map[string]Func
}

func (m *manager) Register(cmd string, fn Func) error {
	if _, exists := m.registry[cmd]; exists {
		return fmt.Errorf("command '%s' already registered", cmd)
	}
	m.registry[cmd] = fn
	return nil
}

func (m *manager) CanExec(cmd string) bool {
	_, registered := m.registry[cmd]
	return registered
}

func (m *manager) Exec(ctx context.Context, cmd string, input Input) (*Output, error) {
	fn, registered := m.registry[cmd]
	if !registered {
		return nil, fmt.Errorf("command '%s' not supported", cmd)
	}
	return fn(ctx, input)
}
