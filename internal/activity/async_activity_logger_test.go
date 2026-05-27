package activity

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

type fakeActivityLogger struct {
	logs             []ActivityEvent
	logErr           error
	latest           []ActivityLog
	findLatestCalled bool
	findLatestLimit  int64
	mu               sync.Mutex
}

func newFakeActivityLogger() *fakeActivityLogger {
	return &fakeActivityLogger{}
}

func (l *fakeActivityLogger) Log(ctx context.Context, eventType string, message string, userID *int, metadata map[string]interface{}) error {
	if l.logErr != nil {
		return l.logErr
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	l.logs = append(l.logs, ActivityEvent{
		Type:     eventType,
		Message:  message,
		UserID:   userID,
		Metadata: metadata,
	})
	return nil
}

func (l *fakeActivityLogger) FindLatest(ctx context.Context, limit int64) ([]ActivityLog, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.findLatestCalled = true
	l.findLatestLimit = limit
	return l.latest, nil
}

func (l *fakeActivityLogger) eventsLen() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.logs)
}

func (l *fakeActivityLogger) eventAt(index int) ActivityEvent {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.logs[index]
}

func TestAsyncActivityLogger_LogProcessesEvent(t *testing.T) {
	base := newFakeActivityLogger()
	logger := NewAsyncActivityLogger(base, 10)
	logger.Start()
	defer logger.Stop()

	userID := 7
	metadata := map[string]interface{}{
		"source": "test",
	}

	err := logger.Log(context.Background(), "test_event", "Test event", &userID, metadata)
	if err != nil {
		t.Fatalf("Log returned error: %v", err)
	}

	waitFor(t, func() bool {
		return base.eventsLen() == 1
	})

	event := base.eventAt(0)
	if event.Type != "test_event" {
		t.Fatalf("expected event type test_event, got %q", event.Type)
	}
	if event.Message != "Test event" {
		t.Fatalf("expected message Test event, got %q", event.Message)
	}
	if event.UserID == nil || *event.UserID != userID {
		t.Fatalf("expected userID %d, got %#v", userID, event.UserID)
	}
	if event.Metadata["source"] != "test" {
		t.Fatalf("expected metadata source test, got %#v", event.Metadata)
	}
}

func TestAsyncActivityLogger_FindLatestDelegatesToBase(t *testing.T) {
	base := newFakeActivityLogger()
	base.latest = []ActivityLog{{Type: "latest_event"}}
	logger := NewAsyncActivityLogger(base, 10)

	logs, err := logger.FindLatest(context.Background(), 5)
	if err != nil {
		t.Fatalf("FindLatest returned error: %v", err)
	}

	if len(logs) != 1 || logs[0].Type != "latest_event" {
		t.Fatalf("unexpected logs: %#v", logs)
	}

	base.mu.Lock()
	defer base.mu.Unlock()
	if !base.findLatestCalled {
		t.Fatal("expected base FindLatest to be called")
	}
	if base.findLatestLimit != 5 {
		t.Fatalf("expected limit 5, got %d", base.findLatestLimit)
	}
}

func TestAsyncActivityLogger_QueueFullReturnsError(t *testing.T) {
	base := newFakeActivityLogger()
	logger := NewAsyncActivityLogger(base, 1)

	err := logger.Log(context.Background(), "first", "First", nil, nil)
	if err != nil {
		t.Fatalf("first Log returned error: %v", err)
	}

	err = logger.Log(context.Background(), "second", "Second", nil, nil)
	if err == nil {
		t.Fatal("expected queue full error")
	}
}

func TestAsyncActivityLogger_StopCanBeCalledTwice(t *testing.T) {
	base := newFakeActivityLogger()
	logger := NewAsyncActivityLogger(base, 10)
	logger.Start()

	defer func() {
		if recovered := recover(); recovered != nil {
			t.Fatalf("Stop panicked: %v", recovered)
		}
	}()

	logger.Stop()
	logger.Stop()
}

func TestAsyncActivityLogger_LogErrorDoesNotStopWorker(t *testing.T) {
	base := newFakeActivityLogger()
	base.logErr = errors.New("write failed")
	logger := NewAsyncActivityLogger(base, 10)
	logger.Start()
	defer logger.Stop()

	err := logger.Log(context.Background(), "test_event", "Test event", nil, nil)
	if err != nil {
		t.Fatalf("Log returned error: %v", err)
	}
}

func waitFor(t *testing.T, condition func() bool) {
	t.Helper()

	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		if condition() {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}

	t.Fatal("condition was not met before deadline")
}
