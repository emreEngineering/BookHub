package activity

import (
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
	go func() {
		for {
			select {
			case event := <-l.events:
				err := l.base.Log(event.Type, event.Message, event.UserID, event.Metadata)
				if err != nil {
					fmt.Println("Activity log yazılamadı:", err)
				}
			case <-l.done:
				return
			}
		}
	}()
}

func (l *AsyncActivityLogger) Stop() {
	l.stopOnce.Do(func() {
		close(l.done)
	})
}

func (l *AsyncActivityLogger) Log(eventType string, message string, userID *int, metadata map[string]interface{}) error {
	event := ActivityEvent{
		Type:     eventType,
		Message:  message,
		UserID:   userID,
		Metadata: metadata,
	}

	select {
	case l.events <- event:
		return nil
	default:
		return errors.New("activity log kuyruğu dolu")
	}
}

func (l *AsyncActivityLogger) FindLatest(limit int64) ([]ActivityLog, error) {
	return l.base.FindLatest(limit)
}
