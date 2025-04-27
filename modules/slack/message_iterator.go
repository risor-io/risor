package slack

import (
	"context"

	"github.com/risor-io/risor/object"
	"github.com/slack-go/slack"
)

// NewMessageIterator creates a new iterator for messages in a channel
func NewMessageIterator(ctx context.Context, client *slack.Client, params *slack.GetConversationHistoryParameters) *GenericIterator {
	// Store state in closures
	var messages []slack.Message
	var hasMore bool
	var currentIndex int
	var oldestTimestamp string
	var totalReturned int
	var limit int

	// Initialize with the provided timestamps
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
		if currentIndex >= len(messages) {
			// If there are no more messages to fetch, we're done
			if !hasMore {
				return nil, false, nil
			}
			if oldestTimestamp != "" {
				params.Oldest = oldestTimestamp
			}
			history, err := client.GetConversationHistoryContext(ctx, params)
			if err != nil {
				return nil, false, err
			}

			messages = history.Messages
			currentIndex = 0
			hasMore = history.HasMore

			// If no messages were returned, we're done
			if len(messages) == 0 {
				return nil, false, nil
			}

			// Update the oldest timestamp for pagination (get older messages)
			if len(messages) > 0 {
				// Get the oldest timestamp from this batch for next pagination
				oldestTimestamp = messages[len(messages)-1].Timestamp

				// Slack requires subtracting a tiny bit to avoid duplicates
				// Convert to float and subtract a small amount
				ts, _ := object.AsFloat(object.NewString(oldestTimestamp))
				oldestTimestamp = object.NewFloat(ts - 0.000001).Inspect()
			}
		}

		// Get the current message
		msg := messages[currentIndex]

		// Create the Message object
		messageObj := NewMessage(client, &msg)
		if messageObj.value.Channel == "" {
			// The Slack API doesn't fill this in for this call
			messageObj.value.Channel = params.ChannelID
		}
		currentIndex++
		totalReturned++

		return messageObj, true, nil
	}

	// Create and return the generic iterator
	return NewGenericIterator("slack.message_iterator", nextFn)
}
