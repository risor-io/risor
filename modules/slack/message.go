package slack

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/risor-io/risor/object"
	"github.com/slack-go/slack"
)

// Ensure Message implements the Object interface
var _ object.Object = (*Message)(nil)

const MESSAGE object.Type = "slack.message"

// Message represents a Slack message
type Message struct {
	base
	value        *slack.Message
	client       *slack.Client
	isBotMessage bool
}

func (m *Message) Inspect() string {
	return fmt.Sprintf("slack.message({channel: %q, timestamp: %q, text: %q})",
		m.value.Msg.Channel, m.value.Msg.Timestamp, m.value.Msg.Text)
}

func (m *Message) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.value)
}

func (m *Message) Value() *slack.Message {
	return m.value
}

func (m *Message) Text() string {
	return m.value.Text
}

func (m *Message) Username() string {
	return m.value.Username
}

func (m *Message) IsBotMessage() bool {
	return m.isBotMessage
}

func (m *Message) Equals(other object.Object) object.Object {
	switch other := other.(type) {
	case *Message:
		return object.NewBool(
			(m.value.Timestamp == other.value.Timestamp) &&
				(m.value.Channel == other.value.Channel),
		)
	default:
		return object.False
	}
}

func (m *Message) GetAttr(name string) (object.Object, bool) {
	switch name {
	case "text":
		return object.NewString(m.value.Text), true
	case "channel":
		return object.NewString(m.value.Channel), true
	case "timestamp":
		return object.NewString(m.value.Timestamp), true
	case "thread_timestamp":
		return object.NewString(m.value.ThreadTimestamp), true
	case "user":
		return object.NewString(m.value.User), true
	case "type":
		return object.NewString(m.value.Type), true
	case "client_msg_id":
		return object.NewString(m.value.ClientMsgID), true
	case "is_starred":
		return object.NewBool(m.value.IsStarred), true
	case "subtype":
		return object.NewString(m.value.SubType), true
	case "team":
		return object.NewString(m.value.Team), true
	case "bot_id":
		return object.NewString(m.value.BotID), true
	case "username":
		return object.NewString(m.value.Username), true
	case "permalink":
		return object.NewString(m.value.Permalink), true
	case "reply_count":
		return object.NewInt(int64(m.value.ReplyCount)), true
	case "latest_reply":
		return object.NewString(m.value.LatestReply), true
	case "is_bot_message":
		return object.NewBool(m.isBotMessage), true
	case "conversation":
		return object.NewBuiltin("slack.message.conversation", m.GetConversation), true
	case "reactions":
		if len(m.value.Reactions) > 0 {
			reactions := make([]object.Object, len(m.value.Reactions))
			for i, reaction := range m.value.Reactions {
				users := make([]object.Object, len(reaction.Users))
				for j, user := range reaction.Users {
					users[j] = object.NewString(user)
				}
				reactions[i] = object.NewMap(map[string]object.Object{
					"name":  object.NewString(reaction.Name),
					"count": object.NewInt(int64(reaction.Count)),
					"users": object.NewList(users),
				})
			}
			return object.NewList(reactions), true
		}
		return object.NewList([]object.Object{}), true
	case "json":
		return asMap(m.value), true
	}
	return nil, false
}

func (m *Message) GetTimestamp(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 0 {
		return object.NewArgsError("slack.message.get_timestamp", 0, len(args))
	}
	// Determine thread timestamp - use thread_ts if available, otherwise use ts
	var ts string
	if m.value.ThreadTimestamp != "" {
		ts = m.value.ThreadTimestamp
	} else {
		ts = m.value.Timestamp
	}
	return object.NewString(ts)
}

// GetConversation retrieves the conversation thread for this message
func (m *Message) GetConversation(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 0 {
		return object.NewArgsError("slack.message.conversation", 0, len(args))
	}
	var ts string
	if m.value.ThreadTimestamp != "" {
		ts = m.value.ThreadTimestamp
	} else {
		ts = m.value.Timestamp
	}
	replies, _, _, err := m.client.GetConversationRepliesContext(ctx,
		&slack.GetConversationRepliesParameters{
			ChannelID: m.value.Channel,
			Timestamp: ts,
			Limit:     100,
		})
	if err != nil {
		return object.NewError(err)
	}
	messages := make([]object.Object, len(replies))
	for i, reply := range replies {
		if reply.Channel == "" {
			reply.Channel = m.value.Channel
		}
		messages[i] = NewMessage(m.client, &reply)
	}
	return object.NewList(messages)
}

func NewMessage(client *slack.Client, msg *slack.Message) *Message {
	return &Message{
		client: client,
		value:  msg,
		base: base{
			typeName:       MESSAGE,
			interfaceValue: msg,
		},
	}
}

func NewMessages(client *slack.Client, msgs []slack.Msg) *object.List {
	items := make([]object.Object, len(msgs))
	for i, msg := range msgs {
		items[i] = NewMessage(client, &slack.Message{Msg: msg})
	}
	return object.NewList(items)
}
