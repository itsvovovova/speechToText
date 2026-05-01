package closer

import (
	"errors"
	"log"
	"sync"
)

// Closer collects cleanup functions and runs them in reverse (LIFO) order on Close.
type Closer struct {
	mu    sync.Mutex
	once  sync.Once
	funcs []func() error
}

func New() *Closer {
	return &Closer{}
}

func (c *Closer) Add(f func() error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.funcs = append(c.funcs, f)
}

func (c *Closer) Close() {
	c.once.Do(func() {
		c.mu.Lock()
		funcs := c.funcs
		c.funcs = nil
		c.mu.Unlock()

		errs := make([]error, 0, len(funcs))
		for i := len(funcs) - 1; i >= 0; i-- {
			if err := funcs[i](); err != nil {
				errs = append(errs, err)
			}
		}
		if err := errors.Join(errs...); err != nil {
			log.Printf("closer: shutdown errors: %v", err)
		}
	})
}
