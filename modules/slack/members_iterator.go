package slack

import (
	"context"

	"github.com/risor-io/risor/object"
	"github.com/slack-go/slack"
)

// NewConversationMembersIterator creates a new iterator for members of a Slack conversation
func NewConversationMembersIterator(ctx context.Context, client *slack.Client, params *slack.GetUsersInConversationParameters) *GenericIterator {
	// Store state in closures
	var members []string
	var cursor string
	var currentIndex int
	var totalReturned int
	var limit int

	// Initialize with empty batch and more results available
	hasMore := true

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
		if currentIndex >= len(members) {
			if !hasMore {
				return nil, false, nil
			}
			params.Cursor = cursor

			// Fetch the next batch
			newMembers, nextCursor, err := client.GetUsersInConversationContext(ctx, params)
			if err != nil {
				return nil, false, err
			}

			members = newMembers
			cursor = nextCursor
			currentIndex = 0
			hasMore = nextCursor != ""

			// If no members were returned, we're done
			if len(members) == 0 {
				return nil, false, nil
			}
		}

		// Get the current member and convert to a User object
		memberID := members[currentIndex]

		// Try to fetch user info
		user, err := client.GetUserInfoContext(ctx, memberID)
		if err != nil {
			// Fall back to returning just the ID if we can't get user info
			memberObj := object.NewString(memberID)
			currentIndex++
			totalReturned++
			return memberObj, true, nil
		}

		// Return the user object
		memberObj := &User{
			client: client,
			value:  user,
		}
		currentIndex++
		totalReturned++

		return memberObj, true, nil
	}

	// Create and return the generic iterator
	return NewGenericIterator("slack.conversation_members_iterator", nextFn)
}
