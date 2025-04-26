package slack

import (
	"context"
	"fmt"

	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/op"
	"github.com/slack-go/slack"
)

// MessageBuilder helps build rich Slack messages
type MessageBuilder struct {
	client *Client
	blocks []slack.Block
}

func (b *MessageBuilder) Type() object.Type {
	return "slack.message_builder"
}

func (b *MessageBuilder) Inspect() string {
	return "slack.message_builder()"
}

func (b *MessageBuilder) Interface() interface{} {
	return b.blocks
}

func (b *MessageBuilder) Equals(other object.Object) object.Object {
	if b == other {
		return object.True
	}
	return object.False
}

func (b *MessageBuilder) IsTruthy() bool {
	return true
}

func (b *MessageBuilder) GetAttr(name string) (object.Object, bool) {
	switch name {
	case "add_section":
		return object.NewBuiltin("add_section", b.AddSection), true
	case "add_divider":
		return object.NewBuiltin("add_divider", b.AddDivider), true
	case "add_header":
		return object.NewBuiltin("add_header", b.AddHeader), true
	case "build":
		return object.NewBuiltin("build", b.Build), true
	case "send":
		return object.NewBuiltin("send", b.Send), true
	}
	return nil, false
}

func (b *MessageBuilder) SetAttr(name string, value object.Object) error {
	return fmt.Errorf("type error: cannot set %q on slack.message_builder object", name)
}

func (b *MessageBuilder) Cost() int {
	return 0
}

func (b *MessageBuilder) RunOperation(opType op.BinaryOpType, right object.Object) object.Object {
	return object.Errorf("type error: unsupported operation for slack.message_builder")
}

// AddSection adds a section block to the message
func (b *MessageBuilder) AddSection(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError(fmt.Errorf("wrong number of arguments: got=%d, want=1", len(args)))
	}
	text, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	textBlock := slack.NewTextBlockObject("mrkdwn", text, false, false)
	b.blocks = append(b.blocks, slack.NewSectionBlock(textBlock, nil, nil))
	return b
}

// AddDivider adds a divider block to the message
func (b *MessageBuilder) AddDivider(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 0 {
		return object.NewError(fmt.Errorf("wrong number of arguments: got=%d, want=0", len(args)))
	}
	b.blocks = append(b.blocks, slack.NewDividerBlock())
	return b
}

// AddHeader adds a header block to the message
func (b *MessageBuilder) AddHeader(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError(fmt.Errorf("wrong number of arguments: got=%d, want=1", len(args)))
	}
	text, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	textBlock := slack.NewTextBlockObject("plain_text", text, false, false)
	b.blocks = append(b.blocks, slack.NewHeaderBlock(textBlock))
	return b
}

// Build creates a map representing the message with blocks
func (b *MessageBuilder) Build(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 0 {
		return object.NewError(fmt.Errorf("wrong number of arguments: got=%d, want=0", len(args)))
	}
	blocksValue := make([]object.Object, 0, len(b.blocks))
	for _, block := range b.blocks {
		var blockMap map[string]object.Object
		switch b := block.(type) {
		case *slack.SectionBlock:
			blockMap = map[string]object.Object{
				"type": object.NewString("section"),
			}
			if b.Text != nil {
				blockMap["text"] = object.NewMap(map[string]object.Object{
					"type": object.NewString(string(b.Text.Type)),
					"text": object.NewString(b.Text.Text),
				})
			}
		case *slack.DividerBlock:
			blockMap = map[string]object.Object{
				"type": object.NewString("divider"),
			}
		case *slack.HeaderBlock:
			blockMap = map[string]object.Object{
				"type": object.NewString("header"),
			}
			if b.Text != nil {
				blockMap["text"] = object.NewMap(map[string]object.Object{
					"type": object.NewString(string(b.Text.Type)),
					"text": object.NewString(b.Text.Text),
				})
			}
		default:
			blockMap = map[string]object.Object{
				"type": object.NewString("unknown"),
			}
		}
		blocksValue = append(blocksValue, object.NewMap(blockMap))
	}
	return object.NewMap(map[string]object.Object{
		"blocks": object.NewList(blocksValue),
	})
}

// Send sends the message to a channel
func (b *MessageBuilder) Send(ctx context.Context, args ...object.Object) object.Object {
	if len(args) < 1 || len(args) > 2 {
		return object.NewError(fmt.Errorf("wrong number of arguments: got=%d, want=1-2", len(args)))
	}
	channelID, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	text := ""
	if len(args) > 1 {
		text, err = object.AsString(args[1])
		if err != nil {
			return err
		}
	}
	options := []slack.MsgOption{
		slack.MsgOptionText(text, false),
		slack.MsgOptionBlocks(b.blocks...),
	}
	channelID, timestamp, _, sendErr := b.client.value.SendMessage(channelID, options...)
	if sendErr != nil {
		return object.NewError(sendErr)
	}
	return object.NewMap(map[string]object.Object{
		"channel":   object.NewString(channelID),
		"timestamp": object.NewString(timestamp),
	})
}

func NewMessageBuilder(client *Client) *MessageBuilder {
	return &MessageBuilder{
		client: client,
		blocks: []slack.Block{},
	}
}
