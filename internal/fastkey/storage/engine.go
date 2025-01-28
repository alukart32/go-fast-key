package storage

// Engine defines a key-value data store.
type Engine struct {
	m map[string]string
}

// NewEngine creates a new Engine.
func NewEngine(cap int) *Engine {
	if cap == 0 {
		cap = 128
	}
	return &Engine{
		m: make(map[string]string, cap),
	}
}

// Set sets a new key-value pair.
func (e *Engine) Set(k, v string) error {
	if e == nil {
		return ErrStandByEngine
	}
	if len(k) == 0 {
		return ErrInvalidEntityID
	}
	if len(v) == 0 {
		return ErrInvalidEntityData
	}

	e.m[k] = v
	return nil
}

// Get finds and returns a value by key.
func (e *Engine) Get(k string) (string, error) {
	if e == nil {
		return "", ErrStandByEngine
	}
	if len(k) == 0 {
		return "", ErrInvalidEntityID
	}

	if val, found := e.m[k]; !found {
		return "", ErrNotFound
	} else {
		return val, nil
	}
}

// Del deletes the value by key.
func (e *Engine) Del(k string) error {
	if e == nil {
		return ErrStandByEngine
	}
	if len(k) == 0 {
		return ErrInvalidEntityID
	}

	delete(e.m, k)
	return nil
}
