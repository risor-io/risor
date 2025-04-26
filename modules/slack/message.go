package slack

import (
	"context"
	"fmt"

	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/op"
	"github.com/slack-go/slack"
)

// Message represents a Slack message
type Message struct {
	value  *slack.Message
	client *Client
}

func (m *Message) Type() object.Type {
	return "slack.message"
}

func (m *Message) Inspect() string {
	return "slack.message()"
}

func (m *Message) Interface() interface{} {
	return m.value
}

func (m *Message) Equals(other object.Object) object.Object {
	if m == other {
		return object.True
	}
	return object.False
}

func (m *Message) IsTruthy() bool {
	return true
}

func (m *Message) GetAttr(name string) (object.Object, bool) {
	switch name {
	case "text":
		return object.NewString(m.value.Msg.Text), true
	case "channel":
		return object.NewString(m.value.Msg.Channel), true
	case "timestamp":
		return object.NewString(m.value.Msg.Timestamp), true
	case "thread_ts":
		return object.NewString(m.value.Msg.ThreadTimestamp), true
	case "user":
		return object.NewString(m.value.Msg.User), true
	case "reply":
		return object.NewBuiltin("reply", m.Reply), true
	}
	return nil, false
}

func (m *Message) SetAttr(name string, value object.Object) error {
	return fmt.Errorf("type error: cannot set %q on slack.message object", name)
}

func (m *Message) Cost() int {
	return 0
}

func (m *Message) RunOperation(opType op.BinaryOpType, right object.Object) object.Object {
	return object.Errorf("type error: unsupported operation for slack.message")
}

// Reply sends a reply to the message in the thread
func (m *Message) Reply(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError(fmt.Errorf("wrong number of arguments: got=%d, want=1", len(args)))
	}
	text, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	options := []slack.MsgOption{
		slack.MsgOptionText(text, false),
		slack.MsgOptionTS(m.value.Msg.Timestamp),
	}
	channelID, timestamp, _, sendErr := m.client.value.SendMessage(m.value.Msg.Channel, options...)
	if sendErr != nil {
		return object.NewError(sendErr)
	}
	return object.NewMap(map[string]object.Object{
		"channel":   object.NewString(channelID),
		"timestamp": object.NewString(timestamp),
	})
}

func NewMessage(msg *slack.Message, client *Client) *Message {
	return &Message{
		client: client,
		value:  msg,
	}
}
