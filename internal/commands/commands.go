package commands

import (
	"context"
	"errors"
	"fmt"
)

type Input struct {
	Args []string
}

type Output struct {
	Response string
}

type Func func(context.Context, Input) (*Output, error)

var ErrAlreadyRegistered = errors.New("command already registered in manager")

type ErrNotSupported struct {
	Command string
}

func (e ErrNotSupported) Error() string {
	return fmt.Sprintf("command '%s' not supported", e.Command)
}

type Manager interface {
	Register(string, Func) error
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
		return ErrAlreadyRegistered
	}
	m.registry[cmd] = fn
	return nil
}

func (m *manager) Exec(ctx context.Context, cmd string, input Input) (*Output, error) {
	fn, exists := m.registry[cmd]
	if !exists {
		return nil, ErrNotSupported{Command: cmd}
	}
	return fn(ctx, input)
}
