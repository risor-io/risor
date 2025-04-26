package slack

import (
	"context"
	"fmt"
	"strings"

	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/op"
	"github.com/slack-go/slack"
)

const CLIENT object.Type = "slack.client"

type Client struct {
	value *slack.Client
}

func (c *Client) Type() object.Type {
	return CLIENT
}

func (c *Client) Inspect() string {
	return "slack.client()"
}

func (c *Client) Interface() interface{} {
	return c.value
}

func (c *Client) Equals(other object.Object) object.Object {
	if c == other {
		return object.True
	}
	return object.False
}

func (c *Client) GetAttr(name string) (object.Object, bool) {
	switch name {
	case "get_user_groups":
		return object.NewBuiltin("get_user_groups", c.GetUserGroups), true
	case "get_user_info":
		return object.NewBuiltin("get_user_info", c.GetUserInfo), true
	case "get_billable_info":
		return object.NewBuiltin("get_billable_info", c.GetBillableInfo), true
	case "post_message":
		return object.NewBuiltin("post_message", c.PostMessage), true
	case "post_ephemeral_message":
		return object.NewBuiltin("post_ephemeral_message", c.PostEphemeralMessage), true
	case "update_message":
		return object.NewBuiltin("update_message", c.UpdateMessage), true
	case "delete_message":
		return object.NewBuiltin("delete_message", c.DeleteMessage), true
	case "add_reaction":
		return object.NewBuiltin("add_reaction", c.AddReaction), true
	case "remove_reaction":
		return object.NewBuiltin("remove_reaction", c.RemoveReaction), true
	case "upload_file":
		return object.NewBuiltin("upload_file", c.UploadFile), true
	case "get_channels":
		return object.NewBuiltin("get_conversations", c.GetConversations), true
	case "message_builder":
		return object.NewBuiltin("message_builder", c.MessageBuilder), true
	}
	return nil, false
}

func (c *Client) SetAttr(name string, value object.Object) error {
	return fmt.Errorf("type error: cannot set %q on slack.client object", name)
}

func (c *Client) IsTruthy() bool {
	return true
}

func (c *Client) RunOperation(opType op.BinaryOpType, right object.Object) object.Object {
	return object.Errorf("type error: unsupported operation for slack.client")
}

func (c *Client) Cost() int {
	return 0
}

func (c *Client) MessageBuilder(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 0 {
		return object.NewError(fmt.Errorf("wrong number of arguments: got=%d, want=0", len(args)))
	}
	return NewMessageBuilder(c)
}

// GetUserGroups gets all user groups for the team
func (c *Client) GetUserGroups(ctx context.Context, args ...object.Object) object.Object {
	var options []slack.GetUserGroupsOption

	if len(args) > 0 {
		optsMap, ok := args[0].(*object.Map)
		if !ok {
			return object.NewError(fmt.Errorf("options must be a map"))
		}
		includeUsers := optsMap.Get("include_users")
		if includeUsers != object.Nil {
			includeBool, err := object.AsBool(includeUsers)
			if err != nil {
				return err
			}
			options = append(options, slack.GetUserGroupsOptionIncludeUsers(bool(includeBool)))
		}
		includeCount := optsMap.Get("include_count")
		if includeCount != object.Nil {
			countBool, err := object.AsBool(includeCount)
			if err != nil {
				return err
			}
			options = append(options, slack.GetUserGroupsOptionIncludeCount(bool(countBool)))
		}
		includeDisabled := optsMap.Get("include_disabled")
		if includeDisabled != object.Nil {
			disabledBool, err := object.AsBool(includeDisabled)
			if err != nil {
				return err
			}
			options = append(options, slack.GetUserGroupsOptionIncludeDisabled(bool(disabledBool)))
		}
		teamID := optsMap.Get("team_id")
		if teamID != object.Nil {
			teamIDStr, err := object.AsString(teamID)
			if err != nil {
				return err
			}
			options = append(options, slack.GetUserGroupsOptionWithTeamID(teamIDStr))
		}
	}

	groups, err := c.value.GetUserGroupsContext(ctx, options...)
	if err != nil {
		return object.NewError(err)
	}
	result := make([]object.Object, len(groups))
	for i, group := range groups {
		var groupUsers []object.Object
		if len(group.Users) > 0 {
			for _, user := range group.Users {
				groupUsers = append(groupUsers, object.NewString(user))
			}
		}
		groupMap := map[string]object.Object{
			"id":           object.NewString(group.ID),
			"team_id":      object.NewString(group.TeamID),
			"is_usergroup": object.NewBool(group.IsUserGroup),
			"name":         object.NewString(group.Name),
			"description":  object.NewString(group.Description),
			"handle":       object.NewString(group.Handle),
			"is_external":  object.NewBool(group.IsExternal),
			"date_create":  object.NewString(group.DateCreate.String()),
			"date_update":  object.NewString(group.DateUpdate.String()),
			"date_delete":  object.NewString(group.DateDelete.String()),
			"auto_type":    object.NewString(group.AutoType),
			"created_by":   object.NewString(group.CreatedBy),
			"updated_by":   object.NewString(group.UpdatedBy),
			"deleted_by":   object.NewString(group.DeletedBy),
			"user_count":   object.NewInt(int64(group.UserCount)),
			"users":        object.NewList(groupUsers),
		}
		result[i] = object.NewMap(groupMap)
	}
	return object.NewList(result)
}

// GetUserInfo gets information about a user
func (c *Client) GetUserInfo(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError(fmt.Errorf("wrong number of arguments: got=%d, want=1", len(args)))
	}
	userID, objErr := object.AsString(args[0])
	if objErr != nil {
		return objErr
	}
	user, err := c.value.GetUserInfo(userID)
	if err != nil {
		return object.NewError(err)
	}
	profileMap := map[string]object.Object{
		"real_name":               object.NewString(user.Profile.RealName),
		"real_name_normalized":    object.NewString(user.Profile.RealNameNormalized),
		"display_name":            object.NewString(user.Profile.DisplayName),
		"display_name_normalized": object.NewString(user.Profile.DisplayNameNormalized),
		"email":                   object.NewString(user.Profile.Email),
		"first_name":              object.NewString(user.Profile.FirstName),
		"last_name":               object.NewString(user.Profile.LastName),
		"phone":                   object.NewString(user.Profile.Phone),
		"skype":                   object.NewString(user.Profile.Skype),
		"title":                   object.NewString(user.Profile.Title),
		"team":                    object.NewString(user.Profile.Team),
		"status_text":             object.NewString(user.Profile.StatusText),
		"status_emoji":            object.NewString(user.Profile.StatusEmoji),
		"bot_id":                  object.NewString(user.Profile.BotID),
		"image_24":                object.NewString(user.Profile.Image24),
		"image_32":                object.NewString(user.Profile.Image32),
		"image_48":                object.NewString(user.Profile.Image48),
		"image_72":                object.NewString(user.Profile.Image72),
		"image_192":               object.NewString(user.Profile.Image192),
		"image_512":               object.NewString(user.Profile.Image512),
		"image_original":          object.NewString(user.Profile.ImageOriginal),
	}
	userMap := map[string]object.Object{
		"id":                  object.NewString(user.ID),
		"team_id":             object.NewString(user.TeamID),
		"name":                object.NewString(user.Name),
		"real_name":           object.NewString(user.RealName),
		"deleted":             object.NewBool(user.Deleted),
		"color":               object.NewString(user.Color),
		"tz":                  object.NewString(user.TZ),
		"tz_label":            object.NewString(user.TZLabel),
		"tz_offset":           object.NewInt(int64(user.TZOffset)),
		"is_bot":              object.NewBool(user.IsBot),
		"is_admin":            object.NewBool(user.IsAdmin),
		"is_owner":            object.NewBool(user.IsOwner),
		"is_primary_owner":    object.NewBool(user.IsPrimaryOwner),
		"is_restricted":       object.NewBool(user.IsRestricted),
		"is_ultra_restricted": object.NewBool(user.IsUltraRestricted),
		"is_app_user":         object.NewBool(user.IsAppUser),
		"is_stranger":         object.NewBool(user.IsStranger),
		"is_invited_user":     object.NewBool(user.IsInvitedUser),
		"has_2fa":             object.NewBool(user.Has2FA),
		"has_files":           object.NewBool(user.HasFiles),
		"locale":              object.NewString(user.Locale),
		"presence":            object.NewString(user.Presence),
		"profile":             object.NewMap(profileMap),
	}
	if user.TwoFactorType != nil {
		userMap["two_factor_type"] = object.NewString(*user.TwoFactorType)
	}
	return object.NewMap(userMap)
}

// GetBillableInfo gets the billable info for the team
func (c *Client) GetBillableInfo(ctx context.Context, args ...object.Object) object.Object {
	params := slack.GetBillableInfoParams{}

	if len(args) > 0 {
		optsMap, ok := args[0].(*object.Map)
		if !ok {
			return object.NewError(fmt.Errorf("options must be a map"))
		}
		userID := optsMap.Get("user")
		if userID != object.Nil {
			userIDStr, err := object.AsString(userID)
			if err != nil {
				return err
			}
			params.User = userIDStr
		}
		teamID := optsMap.Get("team_id")
		if teamID != object.Nil {
			teamIDStr, err := object.AsString(teamID)
			if err != nil {
				return err
			}
			params.TeamID = teamIDStr
		}
	}

	billableInfo, err := c.value.GetBillableInfo(params)
	if err != nil {
		return object.NewError(err)
	}

	result := map[string]object.Object{}
	for userId, info := range billableInfo {
		result[userId] = object.NewMap(map[string]object.Object{
			"billing_active": object.NewBool(info.BillingActive),
		})
	}
	return object.NewMap(result)
}

// PostMessage sends a message to a channel
func (c *Client) PostMessage(ctx context.Context, args ...object.Object) object.Object {
	if len(args) < 1 {
		return object.NewError(fmt.Errorf("wrong number of arguments: got=%d, want at least 1", len(args)))
	}

	channelID, err := object.AsString(args[0])
	if err != nil {
		return err
	}

	options := []slack.MsgOption{}

	// Handle second argument - can be either text string or options map
	if len(args) > 1 {
		switch arg := args[1].(type) {
		case *object.String:
			// Traditional usage with text as second parameter
			text, err := object.AsString(arg)
			if err != nil {
				return err
			}
			options = append(options, slack.MsgOptionText(text, false))

		case *object.Map:
			// New usage with map as second parameter
			// Extract text from map if provided
			textObj := arg.Get("text")
			if textObj != object.Nil {
				text, err := object.AsString(textObj)
				if err != nil {
					return err
				}
				options = append(options, slack.MsgOptionText(text, false))
			}

			// Process other options from the map
			c.processMessageOptions(arg, &options)

		default:
			return object.NewError(fmt.Errorf("second argument must be a string or a map"))
		}
	}

	if len(args) > 2 {
		optsMap, ok := args[2].(*object.Map)
		if !ok {
			return object.NewError(fmt.Errorf("options must be a map"))
		}
		c.processMessageOptions(optsMap, &options)
	}

	channelID, timestamp, _, sendErr := c.value.SendMessage(channelID, options...)
	if sendErr != nil {
		return object.NewError(sendErr)
	}
	return object.NewMap(map[string]object.Object{
		"channel":   object.NewString(channelID),
		"timestamp": object.NewString(timestamp),
	})
}

// processMessageOptions handles common message options processing
func (c *Client) processMessageOptions(optsMap *object.Map, options *[]slack.MsgOption) {
	// Check for thread_ts
	threadTs := optsMap.Get("thread_ts")
	if threadTs != object.Nil {
		threadTsStr, err := object.AsString(threadTs)
		if err == nil {
			*options = append(*options, slack.MsgOptionTS(threadTsStr))
		}
	}

	// Check for reply_broadcast
	replyBroadcast := optsMap.Get("reply_broadcast")
	if replyBroadcast != object.Nil {
		broadcast, err := object.AsBool(replyBroadcast)
		if err == nil && bool(broadcast) {
			*options = append(*options, slack.MsgOptionBroadcast())
		}
	}

	// Check for attachments
	attachments := optsMap.Get("attachments")
	if attachments != object.Nil {
		attachmentsArray, ok := attachments.(*object.List)
		if ok {
			slackAttachments := []slack.Attachment{}
			for _, attachObj := range attachmentsArray.Value() {
				attachMap, ok := attachObj.(*object.Map)
				if !ok {
					continue
				}

				attachment := slack.Attachment{}

				title := attachMap.Get("title")
				if title != object.Nil {
					titleStr, err := object.AsString(title)
					if err == nil {
						attachment.Title = titleStr
					}
				}

				text := attachMap.Get("text")
				if text != object.Nil {
					textStr, err := object.AsString(text)
					if err == nil {
						attachment.Text = textStr
					}
				}

				color := attachMap.Get("color")
				if color != object.Nil {
					colorStr, err := object.AsString(color)
					if err == nil {
						attachment.Color = colorStr
					}
				}

				slackAttachments = append(slackAttachments, attachment)
			}

			*options = append(*options, slack.MsgOptionAttachments(slackAttachments...))
		}
	}

	// Check for blocks
	blocks := optsMap.Get("blocks")
	if blocks != object.Nil {
		blocksArray, ok := blocks.(*object.List)
		if ok {
			slackBlocks := []slack.Block{}
			for _, blockObj := range blocksArray.Value() {
				blockMap, ok := blockObj.(*object.Map)
				if !ok {
					continue
				}

				blockType := blockMap.Get("type")
				if blockType == object.Nil {
					continue
				}

				blockTypeStr, err := object.AsString(blockType)
				if err != nil {
					continue
				}

				// Handle different block types
				switch blockTypeStr {
				case "section":
					// Handle text in section
					textObj := blockMap.Get("text")
					if textObj != object.Nil {
						textMap, ok := textObj.(*object.Map)
						if !ok {
							continue
						}

						textType := textMap.Get("type")
						textValue := textMap.Get("text")

						if textType != object.Nil && textValue != object.Nil {
							typeStr, err := object.AsString(textType)
							if err != nil {
								continue
							}

							valueStr, err := object.AsString(textValue)
							if err != nil {
								continue
							}

							var textBlock *slack.TextBlockObject
							if typeStr == "mrkdwn" {
								textBlock = slack.NewTextBlockObject("mrkdwn", valueStr, false, false)
							} else if typeStr == "plain_text" {
								textBlock = slack.NewTextBlockObject("plain_text", valueStr, false, false)
							}

							section := slack.NewSectionBlock(textBlock, nil, nil)
							slackBlocks = append(slackBlocks, section)
						}
					} else {
						section := slack.NewSectionBlock(nil, nil, nil)
						slackBlocks = append(slackBlocks, section)
					}

				case "divider":
					slackBlocks = append(slackBlocks, slack.NewDividerBlock())

				case "header":
					textObj := blockMap.Get("text")
					if textObj != object.Nil {
						textMap, ok := textObj.(*object.Map)
						if ok {
							textValue := textMap.Get("text")
							if textValue != object.Nil {
								valueStr, err := object.AsString(textValue)
								if err == nil {
									headerText := slack.NewTextBlockObject("plain_text", valueStr, false, false)
									header := slack.NewHeaderBlock(headerText)
									slackBlocks = append(slackBlocks, header)
								}
							}
						}
					}
				}
			}

			*options = append(*options, slack.MsgOptionBlocks(slackBlocks...))
		}
	}
}

// UploadFile uploads a file to Slack
func (c *Client) UploadFile(ctx context.Context, args ...object.Object) object.Object {
	if len(args) < 2 {
		return object.NewError(fmt.Errorf("wrong number of arguments: got=%d, want at least 2", len(args)))
	}
	content, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	channel, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	params := slack.UploadFileV2Parameters{
		Content: content,
		Channel: channel,
	}

	if len(args) > 2 {
		optsMap, ok := args[2].(*object.Map)
		if !ok {
			return object.NewError(fmt.Errorf("options must be a map"))
		}
		filename := optsMap.Get("filename")
		if filename != object.Nil {
			filenameStr, err := object.AsString(filename)
			if err != nil {
				return err
			}
			params.Filename = filenameStr
		}
		title := optsMap.Get("title")
		if title != object.Nil {
			titleStr, err := object.AsString(title)
			if err != nil {
				return err
			}
			params.Title = titleStr
		}
	}

	file, uploadErr := c.value.UploadFileV2(params)
	if uploadErr != nil {
		return object.NewError(uploadErr)
	}
	return object.NewMap(map[string]object.Object{
		"id":    object.NewString(file.ID),
		"title": object.NewString(file.Title),
	})
}

// GetConversations gets all conversations for a user
func (c *Client) GetConversations(ctx context.Context, args ...object.Object) object.Object {
	params := &slack.GetConversationsParameters{
		ExcludeArchived: true, // Default to excluding archived
		Limit:           100,  // Default limit to 100
	}

	// Process options if provided
	if len(args) > 0 {
		optsMap, ok := args[0].(*object.Map)
		if !ok {
			return object.NewError(fmt.Errorf("options must be a map"))
		}

		types := optsMap.Get("types")
		if types != object.Nil {
			typesArray, ok := types.(*object.List)
			if !ok {
				return object.NewError(fmt.Errorf("types must be an array"))
			}

			for _, typeObj := range typesArray.Value() {
				typeStr, err := object.AsString(typeObj)
				if err != nil {
					return err
				}
				params.Types = append(params.Types, typeStr)
			}
		}

		excludeArchived := optsMap.Get("exclude_archived")
		if excludeArchived != object.Nil {
			excludeArchivedBool, err := object.AsBool(excludeArchived)
			if err != nil {
				return err
			}
			params.ExcludeArchived = bool(excludeArchivedBool)
		}

		limit := optsMap.Get("limit")
		if limit != object.Nil {
			limitInt, err := object.AsInt(limit)
			if err != nil {
				return err
			}
			params.Limit = int(limitInt)
		}

		cursor := optsMap.Get("cursor")
		if cursor != object.Nil {
			cursorStr, err := object.AsString(cursor)
			if err != nil {
				return err
			}
			params.Cursor = cursorStr
		}
	}

	channels, nextCursor, err := c.value.GetConversations(params)
	if err != nil {
		return object.NewError(err)
	}

	channelsArray := []object.Object{}
	for _, channel := range channels {
		channelMap := map[string]object.Object{
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
		}

		channelsArray = append(channelsArray, object.NewMap(channelMap))
	}
	return object.NewMap(map[string]object.Object{
		"channels":    object.NewList(channelsArray),
		"next_cursor": object.NewString(nextCursor),
	})
}

// PostEphemeralMessage sends a message to a channel that is only visible to a single user
func (c *Client) PostEphemeralMessage(ctx context.Context, args ...object.Object) object.Object {
	if len(args) < 2 {
		return object.NewError(fmt.Errorf("wrong number of arguments: got=%d, want at least 2", len(args)))
	}

	channelID, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	userID, err := object.AsString(args[1])
	if err != nil {
		return err
	}

	options := []slack.MsgOption{}

	// Handle third argument - can be either text string or options map
	if len(args) > 2 {
		switch arg := args[2].(type) {
		case *object.String:
			// Text as third parameter
			text, err := object.AsString(arg)
			if err != nil {
				return err
			}
			options = append(options, slack.MsgOptionText(text, false))

		case *object.Map:
			// Map as third parameter
			// Extract text from map if provided
			textObj := arg.Get("text")
			if textObj != object.Nil {
				text, err := object.AsString(textObj)
				if err != nil {
					return err
				}
				options = append(options, slack.MsgOptionText(text, false))
			}

			// Process other options from the map
			c.processMessageOptions(arg, &options)

		default:
			return object.NewError(fmt.Errorf("third argument must be a string or a map"))
		}
	} else {
		// If no message content is provided
		return object.NewError(fmt.Errorf("message text or options are required"))
	}

	if len(args) > 3 {
		optsMap, ok := args[3].(*object.Map)
		if !ok {
			return object.NewError(fmt.Errorf("options must be a map"))
		}
		c.processMessageOptions(optsMap, &options)
	}

	ts, sendErr := c.value.PostEphemeralContext(ctx, channelID, userID, options...)
	if sendErr != nil {
		return object.NewError(sendErr)
	}
	return object.NewMap(map[string]object.Object{
		"channel":   object.NewString(channelID),
		"timestamp": object.NewString(ts),
		"user":      object.NewString(userID),
	})
}

// UpdateMessage updates a message in a channel
func (c *Client) UpdateMessage(ctx context.Context, args ...object.Object) object.Object {
	if len(args) < 3 {
		return object.NewError(fmt.Errorf("wrong number of arguments: got=%d, want at least 3", len(args)))
	}

	channelID, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	timestamp, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	options := []slack.MsgOption{}

	// Handle third argument - can be either text string or options map
	switch arg := args[2].(type) {
	case *object.String:
		// Text as third parameter
		text, err := object.AsString(arg)
		if err != nil {
			return err
		}
		options = append(options, slack.MsgOptionText(text, false))

	case *object.Map:
		// Map as third parameter
		// Extract text from map if provided
		textObj := arg.Get("text")
		if textObj != object.Nil {
			text, err := object.AsString(textObj)
			if err != nil {
				return err
			}
			options = append(options, slack.MsgOptionText(text, false))
		}

		// Process other options from the map
		c.processMessageOptions(arg, &options)

	default:
		return object.NewError(fmt.Errorf("third argument must be a string or a map"))
	}

	if len(args) > 3 {
		optsMap, ok := args[3].(*object.Map)
		if !ok {
			return object.NewError(fmt.Errorf("options must be a map"))
		}
		c.processMessageOptions(optsMap, &options)
	}
	channelID, newTimestamp, _, updateErr := c.value.UpdateMessageContext(ctx, channelID, timestamp, options...)
	if updateErr != nil {
		return object.NewError(updateErr)
	}
	return object.NewMap(map[string]object.Object{
		"channel":   object.NewString(channelID),
		"timestamp": object.NewString(newTimestamp),
	})
}

// DeleteMessage deletes a message from a channel
func (c *Client) DeleteMessage(ctx context.Context, args ...object.Object) object.Object {
	if len(args) < 2 {
		return object.NewError(fmt.Errorf("wrong number of arguments: got=%d, want=2", len(args)))
	}
	channelID, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	timestamp, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	_, _, deleteErr := c.value.DeleteMessageContext(ctx, channelID, timestamp)
	if deleteErr != nil {
		return object.NewError(deleteErr)
	}
	return object.Nil
}

// AddReaction adds an emoji reaction to a message
func (c *Client) AddReaction(ctx context.Context, args ...object.Object) object.Object {
	if len(args) < 2 {
		return object.NewError(fmt.Errorf("wrong number of arguments: got=%d, want at least 2", len(args)))
	}

	emojiName, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	emojiName = strings.Trim(emojiName, ":")

	var itemRef slack.ItemRef

	if itemMap, ok := args[1].(*object.Map); ok {
		// Process map to extract ItemRef fields
		channel := itemMap.Get("channel")
		if channel != object.Nil {
			channelStr, err := object.AsString(channel)
			if err != nil {
				return err
			}
			itemRef.Channel = channelStr
		}

		timestamp := itemMap.Get("timestamp")
		if timestamp != object.Nil {
			timestampStr, err := object.AsString(timestamp)
			if err != nil {
				return err
			}
			itemRef.Timestamp = timestampStr
		}

		file := itemMap.Get("file")
		if file != object.Nil {
			fileStr, err := object.AsString(file)
			if err != nil {
				return err
			}
			itemRef.File = fileStr
		}

		comment := itemMap.Get("file_comment")
		if comment != object.Nil {
			commentStr, err := object.AsString(comment)
			if err != nil {
				return err
			}
			itemRef.Comment = commentStr
		}
	} else {
		return object.NewError(fmt.Errorf("second argument must be an item reference map"))
	}

	// Validate we have enough information to identify an item
	if itemRef.Channel == "" && itemRef.File == "" {
		return object.NewError(fmt.Errorf("item reference must include either channel or file"))
	}

	addErr := c.value.AddReactionContext(ctx, emojiName, itemRef)
	if addErr != nil {
		return object.NewError(addErr)
	}

	resultMap := map[string]object.Object{
		"emoji": object.NewString(emojiName),
		"added": object.True,
	}
	if itemRef.Channel != "" {
		resultMap["channel"] = object.NewString(itemRef.Channel)
	}
	if itemRef.Timestamp != "" {
		resultMap["timestamp"] = object.NewString(itemRef.Timestamp)
	}
	if itemRef.File != "" {
		resultMap["file"] = object.NewString(itemRef.File)
	}
	if itemRef.Comment != "" {
		resultMap["file_comment"] = object.NewString(itemRef.Comment)
	}
	return object.NewMap(resultMap)
}

// RemoveReaction removes an emoji reaction from a message
func (c *Client) RemoveReaction(ctx context.Context, args ...object.Object) object.Object {
	if len(args) < 2 {
		return object.NewError(fmt.Errorf("wrong number of arguments: got=%d, want at least 2", len(args)))
	}

	emojiName, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	emojiName = strings.Trim(emojiName, ":")

	var itemRef slack.ItemRef

	if itemMap, ok := args[1].(*object.Map); ok {
		// Process map to extract ItemRef fields
		channel := itemMap.Get("channel")
		if channel != object.Nil {
			channelStr, err := object.AsString(channel)
			if err != nil {
				return err
			}
			itemRef.Channel = channelStr
		}

		timestamp := itemMap.Get("timestamp")
		if timestamp != object.Nil {
			timestampStr, err := object.AsString(timestamp)
			if err != nil {
				return err
			}
			itemRef.Timestamp = timestampStr
		}

		file := itemMap.Get("file")
		if file != object.Nil {
			fileStr, err := object.AsString(file)
			if err != nil {
				return err
			}
			itemRef.File = fileStr
		}

		comment := itemMap.Get("file_comment")
		if comment != object.Nil {
			commentStr, err := object.AsString(comment)
			if err != nil {
				return err
			}
			itemRef.Comment = commentStr
		}
	} else {
		return object.NewError(fmt.Errorf("second argument must be an item reference map"))
	}

	// Validate we have enough information to identify an item
	if itemRef.Channel == "" && itemRef.File == "" {
		return object.NewError(fmt.Errorf("item reference must include either channel or file"))
	}

	removeErr := c.value.RemoveReactionContext(ctx, emojiName, itemRef)
	if removeErr != nil {
		return object.NewError(removeErr)
	}
	return object.Nil
}

func New(client *slack.Client) *Client {
	return &Client{value: client}
}
