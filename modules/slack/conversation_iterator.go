package slack

import (
	"context"
	"fmt"

	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/op"
	"github.com/slack-go/slack"
)

// ConversationIterator implements a Risor iterator that handles paginated conversation results
type ConversationIterator struct {
	client          *slack.Client
	ctx             context.Context
	conversations   []slack.Channel
	currentIndex    int
	currentEntry    *object.Entry
	cursor          string
	hasMore         bool
	excludeArchived bool
	limit           int
	types           []string
}

func (i *ConversationIterator) Type() object.Type {
	return "slack.conversation_iterator"
}

func (i *ConversationIterator) Inspect() string {
	return "slack.conversation_iterator()"
}

func (i *ConversationIterator) Interface() interface{} {
	return i
}

func (i *ConversationIterator) Equals(other object.Object) object.Object {
	if i == other {
		return object.True
	}
	return object.False
}

func (i *ConversationIterator) GetAttr(name string) (object.Object, bool) {
	return nil, false
}

func (i *ConversationIterator) SetAttr(name string, value object.Object) error {
	return fmt.Errorf("cannot set attribute on conversation iterator")
}

func (i *ConversationIterator) IsTruthy() bool {
	return true
}

func (i *ConversationIterator) RunOperation(opType op.BinaryOpType, right object.Object) object.Object {
	return object.Errorf("operation not supported on conversation iterator")
}

func (i *ConversationIterator) Cost() int {
	return 1
}

// Convert a slack.Channel to a Risor map
func convertChannelToMap(channel slack.Channel) *object.Map {
	return object.NewMap(map[string]object.Object{
		"id":              object.NewString(channel.ID),
		"name":            object.NewString(channel.Name),
		"is_archived":     object.NewBool(channel.IsArchived),
		"is_channel":      object.NewBool(channel.IsChannel),
		"is_private":      object.NewBool(channel.IsPrivate),
		"is_group":        object.NewBool(channel.IsGroup),
		"is_im":           object.NewBool(channel.IsIM),
		"is_mpim":         object.NewBool(channel.IsMpIM),
		"is_general":      object.NewBool(channel.IsGeneral),
		"num_members":     object.NewInt(int64(channel.NumMembers)),
		"creator":         object.NewString(channel.Creator),
		"name_normalized": object.NewString(channel.NameNormalized),
	})
}

func (i *ConversationIterator) Next(ctx context.Context) (object.Object, bool) {
	// If we've reached the end of the current batch
	if i.currentIndex >= len(i.conversations) {
		// Check if there are more conversations to fetch
		if i.hasMore && i.cursor != "" {
			// Fetch the next batch
			params := &slack.GetConversationsParameters{
				Cursor:          i.cursor,
				ExcludeArchived: i.excludeArchived,
				Limit:           i.limit,
				Types:           i.types,
			}

			conversations, nextCursor, err := i.client.GetConversationsContext(i.ctx, params)
			if err != nil {
				return object.NewError(err), false
			}

			i.conversations = conversations
			i.cursor = nextCursor
			i.currentIndex = 0
			i.hasMore = nextCursor != ""

			// If no conversations were returned, we're done
			if len(conversations) == 0 {
				return object.Nil, false
			}
		} else {
			// No more conversations to fetch
			return object.Nil, false
		}
	}

	// Get the current conversation and create an entry
	conversation := i.conversations[i.currentIndex]
	i.currentEntry = object.NewEntry(object.NewInt(int64(i.currentIndex)), convertChannelToMap(conversation))
	i.currentIndex++

	return i.currentEntry.Value(), true
}

func (i *ConversationIterator) Entry() (object.IteratorEntry, bool) {
	if i.currentEntry == nil {
		return nil, false
	}
	return i.currentEntry, true
}
