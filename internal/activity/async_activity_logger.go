package activity

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

type ActivityEvent struct {
	Type     string
	Message  string
	UserID   *int
	Metadata map[string]interface{}
}

type AsyncActivityLogger struct {
	base     ActivityLogger
	events   chan ActivityEvent
	done     chan struct{}
	stopOnce sync.Once
	wg       sync.WaitGroup
	mu       sync.Mutex
	stopped  bool
}

func NewAsyncActivityLogger(base ActivityLogger, bufferSize int) *AsyncActivityLogger {
	if bufferSize <= 0 {
		bufferSize = 100
	}

	return &AsyncActivityLogger{
		base:   base,
		events: make(chan ActivityEvent, bufferSize),
		done:   make(chan struct{}),
	}
}

func (l *AsyncActivityLogger) Start() {
	l.wg.Add(1)
	go func() {
		defer l.wg.Done()

		for {
			select {
			case event := <-l.events:
				l.write(event)
			case <-l.done:
				l.drain()
				return
			}
		}
	}()
}

func (l *AsyncActivityLogger) Stop() {
	l.stopOnce.Do(func() {
		l.mu.Lock()
		l.stopped = true
		l.mu.Unlock()

		close(l.done)
		l.wg.Wait()
	})
}

func (l *AsyncActivityLogger) Log(ctx context.Context, eventType string, message string, userID *int, metadata map[string]interface{}) error {
	l.mu.Lock()
	stopped := l.stopped
	l.mu.Unlock()

	if stopped {
		return errors.New("activity logger durduruldu")
	}

	event := ActivityEvent{
		Type:     eventType,
		Message:  message,
		UserID:   userID,
		Metadata: metadata,
	}

	select {
	case l.events <- event:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	default:
		return errors.New("activity log kuyruğu dolu")
	}
}

func (l *AsyncActivityLogger) FindLatest(ctx context.Context, limit int64) ([]ActivityLog, error) {
	return l.base.FindLatest(ctx, limit)
}

func (l *AsyncActivityLogger) drain() {
	for {
		select {
		case event := <-l.events:
			l.write(event)
		default:
			return
		}
	}
}

func (l *AsyncActivityLogger) write(event ActivityEvent) {
	err := l.base.Log(context.Background(), event.Type, event.Message, event.UserID, event.Metadata)
	if err != nil {
		fmt.Println("Activity log yazılamadı:", err)
	}
}
