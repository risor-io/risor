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
			// If there are no more conversations to fetch, we're done
			if !hasMore || cursor == "" {
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
		channelMap := convertChannelToMap(conversation)
		currentIndex++
		totalReturned++

		return channelMap, true, nil
	}

	// Create and return the generic iterator
	return NewGenericIterator("slack.conversation_iterator", nextFn)
}

// convertChannelToMap converts a slack.Channel to a Risor object.Map
func convertChannelToMap(channel slack.Channel) *object.Map {
	// Create a map for Topic
	topicMap := map[string]object.Object{
		"value":    object.NewString(channel.Topic.Value),
		"creator":  object.NewString(channel.Topic.Creator),
		"last_set": getTime(channel.Topic.LastSet),
	}

	// Create a map for Purpose
	purposeMap := map[string]object.Object{
		"value":    object.NewString(channel.Purpose.Value),
		"creator":  object.NewString(channel.Purpose.Creator),
		"last_set": getTime(channel.Purpose.LastSet),
	}

	// Create a map for members
	membersObjs := make([]object.Object, len(channel.Members))
	for i, member := range channel.Members {
		membersObjs[i] = object.NewString(member)
	}

	// Create properties map if properties exist
	var propertiesMap object.Object = object.Nil
	if channel.Properties != nil {
		propertiesMap = object.NewMap(map[string]object.Object{
			"canvas": object.NewMap(map[string]object.Object{
				"file_id":        object.NewString(channel.Properties.Canvas.FileId),
				"is_empty":       object.NewBool(channel.Properties.Canvas.IsEmpty),
				"quip_thread_id": object.NewString(channel.Properties.Canvas.QuipThreadId),
			}),
		})
	}

	// Combine all fields into a channel map
	channelMap := map[string]object.Object{
		// Conversation fields
		"id":                    object.NewString(channel.ID),
		"created":               getTime(channel.Created),
		"is_open":               object.NewBool(channel.IsOpen),
		"last_read":             object.NewString(channel.LastRead),
		"is_group":              object.NewBool(channel.IsGroup),
		"is_shared":             object.NewBool(channel.IsShared),
		"is_im":                 object.NewBool(channel.IsIM),
		"is_ext_shared":         object.NewBool(channel.IsExtShared),
		"is_org_shared":         object.NewBool(channel.IsOrgShared),
		"is_global_shared":      object.NewBool(channel.IsGlobalShared),
		"is_pending_ext_shared": object.NewBool(channel.IsPendingExtShared),
		"is_private":            object.NewBool(channel.IsPrivate),
		"is_read_only":          object.NewBool(channel.IsReadOnly),
		"is_mpim":               object.NewBool(channel.IsMpIM),
		"unlinked":              object.NewInt(int64(channel.Unlinked)),
		"name_normalized":       object.NewString(channel.NameNormalized),
		"num_members":           object.NewInt(int64(channel.NumMembers)),
		"priority":              object.NewFloat(channel.Priority),
		"user":                  object.NewString(channel.User),

		// GroupConversation fields
		"name":        object.NewString(channel.Name),
		"creator":     object.NewString(channel.Creator),
		"is_archived": object.NewBool(channel.IsArchived),
		"members":     object.NewList(membersObjs),
		"topic":       object.NewMap(topicMap),
		"purpose":     object.NewMap(purposeMap),

		// Channel specific fields
		"is_channel": object.NewBool(channel.IsChannel),
		"is_general": object.NewBool(channel.IsGeneral),
		"is_member":  object.NewBool(channel.IsMember),
		"locale":     object.NewString(channel.Locale),
		"properties": propertiesMap,
	}

	// Add connected team IDs if present
	if len(channel.ConnectedTeamIDs) > 0 {
		teamIDs := make([]object.Object, len(channel.ConnectedTeamIDs))
		for i, id := range channel.ConnectedTeamIDs {
			teamIDs[i] = object.NewString(id)
		}
		channelMap["connected_team_ids"] = object.NewList(teamIDs)
	}

	// Add shared team IDs if present
	if len(channel.SharedTeamIDs) > 0 {
		teamIDs := make([]object.Object, len(channel.SharedTeamIDs))
		for i, id := range channel.SharedTeamIDs {
			teamIDs[i] = object.NewString(id)
		}
		channelMap["shared_team_ids"] = object.NewList(teamIDs)
	}

	// Add internal team IDs if present
	if len(channel.InternalTeamIDs) > 0 {
		teamIDs := make([]object.Object, len(channel.InternalTeamIDs))
		for i, id := range channel.InternalTeamIDs {
			teamIDs[i] = object.NewString(id)
		}
		channelMap["internal_team_ids"] = object.NewList(teamIDs)
	}

	// Add context team ID and conversation host ID if present
	if channel.ContextTeamID != "" {
		channelMap["context_team_id"] = object.NewString(channel.ContextTeamID)
	}
	if channel.ConversationHostID != "" {
		channelMap["conversation_host_id"] = object.NewString(channel.ConversationHostID)
	}

	return object.NewMap(channelMap)
}
