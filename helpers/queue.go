package helpers

import (
	"context"
	"github.com/maragudk/goqite"
	"time"
)

// SendToQueue sends a message to the queue.
func SendToQueue(ctx context.Context, q *goqite.Queue, body []byte) error {
	return q.Send(ctx, goqite.Message{Body: body})
}

// ReceiveFromQueue receives a message from the queue.
func ReceiveFromQueue(ctx context.Context, q *goqite.Queue) ([]byte, error) {
	m, err := q.Receive(ctx)
	if err != nil {
		return nil, err
	}
	return m.Body, nil
}

// ExtendMessageTimeout extends the timeout of a message in the queue.
func ExtendMessageTimeout(ctx context.Context, q *goqite.Queue, id string, timeout time.Duration) error {
	return q.Extend(ctx, goqite.ID(id), timeout)
}

// DeleteMessageFromQueue deletes a message from the queue.
func DeleteMessageFromQueue(ctx context.Context, q *goqite.Queue, id string) error {
	return q.Delete(ctx, goqite.ID(id))
}
