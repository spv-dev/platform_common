package closer

import (
	"log"
	"os"
	"os/signal"
	"sync"
)

// globalCloser глобальный объект - чтобы можно было применять из
// разных мест приложения
// Скрыт. Для доступа используются функции Add, Wait, CloseAll
var globalCloser = New()

// Add функция для доступа к globalCloser.Add()
func Add(f ...func() error) {
	globalCloser.Add(f...)
}

// Wait функция для доступа к globalCloser.Wait()
func Wait() {
	globalCloser.Wait()
}

// CloseAll функция для доступа к globalCloser.CloseAll()
func CloseAll() {
	globalCloser.CloseAll()
}

// Closer структура для закрытия всех инициализированных сервисов
type Closer struct {
	mu    sync.Mutex
	onse  sync.Once
	done  chan struct{}
	funcs []func() error
}

// New создаёт новый объект Closer
func New(sig ...os.Signal) *Closer {
	c := &Closer{done: make(chan struct{})}
	if len(sig) > 0 {
		go func() {
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, sig...)
			<-ch
			signal.Stop(ch)
			c.CloseAll()
		}()
	}
	return c
}

// Add добавляет новую функцию закрытия
func (c *Closer) Add(f ...func() error) {
	c.mu.Lock()
	c.funcs = append(c.funcs, f...)
	c.mu.Unlock()
}

// Wait блокирует вызов
func (c *Closer) Wait() {
	<-c.done
}

// CloseAll вызывает все функции закрытия
func (c *Closer) CloseAll() {
	c.onse.Do(func() {
		defer close(c.done)

		c.mu.Lock()
		funcs := c.funcs
		c.funcs = nil
		c.mu.Unlock()

		errs := make(chan error, len(funcs))
		for _, f := range funcs {
			go func(f func() error) {
				errs <- f()
			}(f)
		}

		for i := 0; i < cap(errs); i++ {
			if err := <-errs; err != nil {
				log.Println("error returned from Closer")
			}
		}
	})
}
