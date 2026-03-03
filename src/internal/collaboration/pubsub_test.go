/*-------------------------------------------------------------------------
 *
 * pubsub_test.go
 *    Tests for publish-subscribe.
 *
 *-------------------------------------------------------------------------
 */

package collaboration

import (
	"context"
	"testing"
)

func TestNewPubSub(t *testing.T) {
	ps := NewPubSub()
	if ps == nil {
		t.Fatal("NewPubSub returned nil")
	}
}

func TestPubSub_SubscribePublish(t *testing.T) {
	ps := NewPubSub()
	ch := ps.Subscribe("test")
	defer ps.Unsubscribe("test", ch)
	ctx := context.Background()
	ps.Publish(ctx, "test", "hello")
	got := <-ch
	if got != "hello" {
		t.Errorf("received %v", got)
	}
}

func TestPubSub_Unsubscribe(t *testing.T) {
	ps := NewPubSub()
	ch := ps.Subscribe("test")
	ps.Unsubscribe("test", ch)
	_, ok := <-ch
	if ok {
		t.Error("channel should be closed after Unsubscribe")
	}
}
