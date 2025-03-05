package engine

// MemEngine defines a key-value data store.
type MemEngine struct {
	m map[string]string
}

// NewMemEngine creates a new Engine.
func NewMemEngine(cap int) *MemEngine {
	if cap == 0 {
		cap = 128
	}
	return &MemEngine{
		m: make(map[string]string, cap),
	}
}

// Set sets a new key-value pair.
func (e *MemEngine) Set(k, v string) error {
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
func (e *MemEngine) Get(k string) (string, error) {
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
func (e *MemEngine) Del(k string) error {
	if len(k) == 0 {
		return ErrInvalidEntityID
	}

	delete(e.m, k)
	return nil
}
