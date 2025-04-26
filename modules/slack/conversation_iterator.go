package slack

import (
	"context"

	"github.com/risor-io/risor/object"
	"github.com/slack-go/slack"
)

// NewConversationIterator creates a new iterator for Slack conversations
func NewConversationIterator(ctx context.Context, client *slack.Client, params *slack.GetConversationsParameters) *GenericIterator {
	// Store state in closures
	var conversations []slack.Channel
	var cursor string
	var hasMore bool
	var currentIndex int
	var totalReturned int
	var limit int

	// Initialize with an empty batch
	hasMore = true

	// Store the limit if provided
	if params != nil && params.Limit > 0 {
		limit = params.Limit
	}

	// Create the next function
	nextFn := func(ctx context.Context) (object.Object, bool, error) {

		// If we've reached the limit, we're done
		if limit > 0 && totalReturned >= limit {
			return nil, false, nil
		}

		// If we've reached the end of the current batch
		if currentIndex >= len(conversations) {
			if !hasMore {
				return nil, false, nil
			}
			params.Cursor = cursor

			// Fetch the next batch
			newConversations, nextCursor, err := client.GetConversationsContext(ctx, params)
			if err != nil {
				return nil, false, err
			}

			conversations = newConversations
			cursor = nextCursor
			currentIndex = 0
			hasMore = nextCursor != ""

			// If no conversations were returned, we're done
			if len(conversations) == 0 {
				return nil, false, nil
			}
		}

		// Get the current conversation and convert to a map
		conversation := conversations[currentIndex]
		channelObj := NewChannel(client, &conversation)
		currentIndex++
		totalReturned++

		return channelObj, true, nil
	}

	// Create and return the generic iterator
	return NewGenericIterator("slack.conversation_iterator", nextFn)
}
