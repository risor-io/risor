package slack

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/op"
	"github.com/slack-go/slack"
)

// Ensure Message implements the Object interface
var _ object.Object = (*Message)(nil)

// Message represents a Slack message
type Message struct {
	value        *slack.Message
	client       *slack.Client
	isBotMessage bool
}

func (m *Message) Type() object.Type {
	return "slack.message"
}

func (m *Message) Inspect() string {
	return fmt.Sprintf("slack.message({channel: %q, timestamp: %q, text: %q})",
		m.value.Msg.Channel, m.value.Msg.Timestamp, m.value.Msg.Text)
}

func (m *Message) Interface() interface{} {
	return m.value
}

func (m *Message) Value() *slack.Message {
	return m.value
}

func (m *Message) Text() string {
	return m.value.Msg.Text
}

func (m *Message) Username() string {
	return m.value.Msg.Username
}

func (m *Message) IsBotMessage() bool {
	return m.isBotMessage
}

func (m *Message) Equals(other object.Object) object.Object {
	switch other := other.(type) {
	case *Message:
		return object.NewBool(
			(m.value.Msg.Timestamp == other.value.Msg.Timestamp) &&
				(m.value.Msg.Channel == other.value.Msg.Channel),
		)
	default:
		return object.False
	}
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
	case "thread_timestamp":
		return object.NewString(m.value.Msg.ThreadTimestamp), true
	case "user":
		return object.NewString(m.value.Msg.User), true
	case "type":
		return object.NewString(m.value.Msg.Type), true
	case "client_msg_id":
		return object.NewString(m.value.Msg.ClientMsgID), true
	case "is_starred":
		return object.NewBool(m.value.Msg.IsStarred), true
	case "subtype":
		return object.NewString(m.value.Msg.SubType), true
	case "team":
		return object.NewString(m.value.Msg.Team), true
	case "bot_id":
		return object.NewString(m.value.Msg.BotID), true
	case "username":
		return object.NewString(m.value.Msg.Username), true
	case "permalink":
		return object.NewString(m.value.Msg.Permalink), true
	case "reply_count":
		return object.NewInt(int64(m.value.Msg.ReplyCount)), true
	case "latest_reply":
		return object.NewString(m.value.Msg.LatestReply), true
	case "is_bot_message":
		return object.NewBool(m.isBotMessage), true
	case "reply":
		return object.NewBuiltin("slack.message.reply", m.Reply), true
	case "conversation":
		return object.NewBuiltin("slack.message.conversation", m.GetConversation), true
	case "reactions":
		if len(m.value.Msg.Reactions) > 0 {
			reactions := make([]object.Object, len(m.value.Msg.Reactions))
			for i, reaction := range m.value.Msg.Reactions {
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

func (m *Message) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.value)
}

// Reply sends a reply to the message in the thread
func (m *Message) Reply(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewArgsError("slack.message.reply", 1, len(args))
	}
	text, err := object.AsString(args[0])
	if err != nil {
		return err
	}

	// Determine thread timestamp - use thread_ts if available, otherwise use ts
	var ts string
	if m.value.Msg.ThreadTimestamp != "" {
		ts = m.value.Msg.ThreadTimestamp
	} else {
		ts = m.value.Msg.Timestamp
	}

	options := []slack.MsgOption{
		slack.MsgOptionText(text, false),
		slack.MsgOptionTS(ts),
	}

	channelID, timestamp, _, sendErr := m.client.SendMessage(m.value.Msg.Channel, options...)
	if sendErr != nil {
		return object.NewError(sendErr)
	}

	// Create a new message object for the reply
	replyMsg := &slack.Message{
		Msg: slack.Msg{
			Channel:         channelID,
			Timestamp:       timestamp,
			Text:            text,
			ThreadTimestamp: ts,
		},
	}
	return NewMessage(replyMsg, m.client)
}

// GetConversation retrieves the conversation thread for this message
func (m *Message) GetConversation(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 0 {
		return object.NewArgsError("slack.message.conversation", 0, len(args))
	}

	// Determine which timestamp to use
	var ts string
	if m.value.Msg.ThreadTimestamp != "" {
		ts = m.value.Msg.ThreadTimestamp
	} else {
		ts = m.value.Msg.Timestamp
	}

	// Get conversation replies
	replies, _, _, err := m.client.GetConversationRepliesContext(ctx,
		&slack.GetConversationRepliesParameters{
			ChannelID: m.value.Msg.Channel,
			Timestamp: ts,
			Limit:     100,
		})
	if err != nil {
		return object.NewError(err)
	}

	// Convert replies to a list of Message objects
	messages := make([]object.Object, len(replies))
	for i, reply := range replies {
		// Channel may not be set in the replies
		if reply.Channel == "" {
			reply.Channel = m.value.Msg.Channel
		}

		// Pass the reply directly to NewMessage as a pointer
		messages[i] = NewMessage(&reply, m.client)
	}

	return object.NewList(messages)
}

func NewMessage(msg *slack.Message, client *slack.Client) *Message {
	return &Message{
		client: client,
		value:  msg,
	}
}

func NewMessages(msgs []slack.Msg, client *Client) *object.List {
	items := make([]object.Object, len(msgs))
	for i, msg := range msgs {
		items[i] = NewMessage(&slack.Message{Msg: msg}, client.value)
	}
	return object.NewList(items)
}
