package slack

import (
	"context"
	"fmt"
	"strings"

	"github.com/risor-io/risor/object"
	"github.com/slack-go/slack"
)

const CLIENT object.Type = "slack.client"

type Client struct {
	base
	value *slack.Client
}

func (c *Client) Inspect() string {
	return "slack.client()"
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
	case "get_users":
		return object.NewBuiltin("get_users", c.GetUsers), true
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
	case "get_conversations":
		return object.NewBuiltin("get_conversations", c.GetConversations), true
	case "create_conversation":
		return object.NewBuiltin("create_conversation", c.CreateConversation), true
	case "get_conversation_history":
		return object.NewBuiltin("get_conversation_history", c.GetConversationHistory), true
	case "get_conversation_members":
		return object.NewBuiltin("get_conversation_members", c.GetConversationMembers), true
	case "message_builder":
		return object.NewBuiltin("message_builder", c.MessageBuilder), true
	}
	return nil, false
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
		opts, ok := args[0].(*object.Map)
		if !ok {
			return object.NewError(fmt.Errorf("options must be a map"))
		}
		for key, value := range opts.Value() {
			switch key {
			case "include_users":
				includeBool, err := object.AsBool(value)
				if err != nil {
					return err
				}
				options = append(options, slack.GetUserGroupsOptionIncludeUsers(bool(includeBool)))
			case "include_count":
				countBool, err := object.AsBool(value)
				if err != nil {
					return err
				}
				options = append(options, slack.GetUserGroupsOptionIncludeCount(bool(countBool)))
			case "include_disabled":
				disabledBool, err := object.AsBool(value)
				if err != nil {
					return err
				}
				options = append(options, slack.GetUserGroupsOptionIncludeDisabled(bool(disabledBool)))
			case "team_id":
				teamIDStr, err := object.AsString(value)
				if err != nil {
					return err
				}
				options = append(options, slack.GetUserGroupsOptionWithTeamID(teamIDStr))
			default:
				return object.NewError(fmt.Errorf("unknown option: %s", key))
			}
		}
	}

	groups, err := c.value.GetUserGroupsContext(ctx, options...)
	if err != nil {
		return object.NewError(err)
	}
	result := make([]object.Object, len(groups))
	for i, group := range groups {
		groupUsers := []object.Object{}
		if len(group.Users) > 0 {
			for _, userID := range group.Users {
				// Just use string IDs for now, users can fetch details separately if needed
				groupUsers = append(groupUsers, object.NewString(userID))
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
			"date_create":  getTime(group.DateCreate),
			"date_update":  getTime(group.DateUpdate),
			"date_delete":  getTime(group.DateDelete),
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
	return NewUser(c.value, user)
}

func (c *Client) GetUsers(ctx context.Context, args ...object.Object) object.Object {
	options := []slack.GetUsersOption{}
	if len(args) > 0 {
		opts, ok := args[0].(*object.Map)
		if !ok {
			return object.NewError(fmt.Errorf("options must be a map"))
		}
		for key, value := range opts.Value() {
			switch key {
			case "limit":
				limitInt, err := object.AsInt(value)
				if err != nil {
					return err
				}
				options = append(options, slack.GetUsersOptionLimit(int(limitInt)))
			case "presence":
				presenceBool, err := object.AsBool(value)
				if err != nil {
					return err
				}
				options = append(options, slack.GetUsersOptionPresence(bool(presenceBool)))
			case "team_id":
				teamIDStr, err := object.AsString(value)
				if err != nil {
					return err
				}
				options = append(options, slack.GetUsersOptionTeamID(teamIDStr))
			default:
				return object.NewError(fmt.Errorf("unknown option: %s", key))
			}
		}
	}
	users, err := c.value.GetUsersContext(ctx, options...)
	if err != nil {
		return object.NewError(err)
	}
	results := make([]object.Object, len(users))
	for i, user := range users {
		results[i] = NewUser(c.value, &user)
	}
	return object.NewList(results)
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
			text, err := object.AsString(arg)
			if err != nil {
				return err
			}
			options = append(options, slack.MsgOptionText(text, false))
		case *object.Map:
			textObj := arg.Get("text")
			if textObj != object.Nil {
				text, err := object.AsString(textObj)
				if err != nil {
					return err
				}
				options = append(options, slack.MsgOptionText(text, false))
			}
			c.processMessageOptions(arg, &options)
		default:
			return object.NewError(fmt.Errorf("second argument must be a string or a map"))
		}
	}

	if len(args) > 2 {
		opts, ok := args[2].(*object.Map)
		if !ok {
			return object.NewError(fmt.Errorf("options must be a map"))
		}
		c.processMessageOptions(opts, &options)
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
func (c *Client) processMessageOptions(opts *object.Map, options *[]slack.MsgOption) {
	// Check for thread_ts
	threadTs := opts.Get("thread_ts")
	if threadTs != object.Nil {
		threadTsStr, err := object.AsString(threadTs)
		if err == nil {
			*options = append(*options, slack.MsgOptionTS(threadTsStr))
		}
	}

	// Check for reply_broadcast
	replyBroadcast := opts.Get("reply_broadcast")
	if replyBroadcast != object.Nil {
		broadcast, err := object.AsBool(replyBroadcast)
		if err == nil && bool(broadcast) {
			*options = append(*options, slack.MsgOptionBroadcast())
		}
	}

	// Check for attachments
	attachments := opts.Get("attachments")
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
	blocks := opts.Get("blocks")
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
				default:
					// Skip unsupported block types instead of raising an error
					// This is more forgiving for block types that may be added in future Slack API updates
					continue
				}
			}

			*options = append(*options, slack.MsgOptionBlocks(slackBlocks...))
		}
	}
}

// UploadFile uploads a file to Slack
func (c *Client) UploadFile(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 2 {
		return object.NewError(fmt.Errorf("wrong number of arguments: got=%d, want=2", len(args)))
	}

	channelID, err := object.AsString(args[0])
	if err != nil {
		return err
	}

	paramsMap, ok := args[1].(*object.Map)
	if !ok {
		return object.NewError(fmt.Errorf("second argument must be a map"))
	}
	params := slack.UploadFileV2Parameters{}
	params.Channel = channelID

	content := paramsMap.Get("content")
	if content != object.Nil {
		contentStr, err := object.AsString(content)
		if err != nil {
			return err
		}
		params.Content = contentStr
	}
	file := paramsMap.Get("file")
	if file != object.Nil {
		fileStr, err := object.AsString(file)
		if err != nil {
			return err
		}
		params.File = fileStr
	}
	fileSize := paramsMap.Get("file_size")
	if fileSize != object.Nil {
		fileSizeInt, err := object.AsInt(fileSize)
		if err != nil {
			return err
		}
		params.FileSize = int(fileSizeInt)
	}
	filename := paramsMap.Get("filename")
	if filename != object.Nil {
		filenameStr, err := object.AsString(filename)
		if err != nil {
			return err
		}
		params.Filename = filenameStr
	}
	title := paramsMap.Get("title")
	if title != object.Nil {
		titleStr, err := object.AsString(title)
		if err != nil {
			return err
		}
		params.Title = titleStr
	}
	initialComment := paramsMap.Get("initial_comment")
	if initialComment != object.Nil {
		initialCommentStr, err := object.AsString(initialComment)
		if err != nil {
			return err
		}
		params.InitialComment = initialCommentStr
	}
	threadTs := paramsMap.Get("thread_ts")
	if threadTs != object.Nil {
		threadTsStr, err := object.AsString(threadTs)
		if err != nil {
			return err
		}
		params.ThreadTimestamp = threadTsStr
	}
	altTxt := paramsMap.Get("alt_txt")
	if altTxt != object.Nil {
		altTxtStr, err := object.AsString(altTxt)
		if err != nil {
			return err
		}
		params.AltTxt = altTxtStr
	}
	snippetText := paramsMap.Get("snippet_text")
	if snippetText != object.Nil {
		snippetTextStr, err := object.AsString(snippetText)
		if err != nil {
			return err
		}
		params.SnippetText = snippetTextStr
	}
	if len(params.Content) > 0 && params.FileSize == 0 {
		params.FileSize = len(params.Content)
	}
	// Validate that we have the minimum required parameters
	if params.Content == "" && params.File == "" {
		return object.NewError(fmt.Errorf("either content or file must be provided"))
	}
	fileSummary, uploadErr := c.value.UploadFileV2(params)
	if uploadErr != nil {
		return object.NewError(uploadErr)
	}
	return object.NewMap(map[string]object.Object{
		"id":    object.NewString(fileSummary.ID),
		"title": object.NewString(fileSummary.Title),
	})
}

// GetConversations gets all conversations for a user
func (c *Client) GetConversations(ctx context.Context, args ...object.Object) object.Object {
	params := &slack.GetConversationsParameters{
		ExcludeArchived: true, // Default to excluding archived
		Limit:           100,  // Default limit to 100
	}
	if len(args) > 0 {
		opts, ok := args[0].(*object.Map)
		if !ok {
			return object.NewError(fmt.Errorf("options must be a map"))
		}

		// Process options with switch statement
		for key, value := range opts.Value() {
			switch key {
			case "types":
				typesArray, ok := value.(*object.List)
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
			case "exclude_archived":
				excludeArchivedBool, err := object.AsBool(value)
				if err != nil {
					return err
				}
				params.ExcludeArchived = bool(excludeArchivedBool)
			case "limit":
				limitInt, err := object.AsInt(value)
				if err != nil {
					return err
				}
				params.Limit = int(limitInt)
			case "team_id":
				teamIDStr, err := object.AsString(value)
				if err != nil {
					return err
				}
				params.TeamID = teamIDStr
			default:
				return object.NewError(fmt.Errorf("unknown option: %s", key))
			}
		}
	}
	return NewConversationIterator(ctx, c.value, params)
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
			text, err := object.AsString(arg)
			if err != nil {
				return err
			}
			options = append(options, slack.MsgOptionText(text, false))
		case *object.Map:
			textObj := arg.Get("text")
			if textObj != object.Nil {
				text, err := object.AsString(textObj)
				if err != nil {
					return err
				}
				options = append(options, slack.MsgOptionText(text, false))
			}
			c.processMessageOptions(arg, &options)
		default:
			return object.NewError(fmt.Errorf("third argument must be a string or a map"))
		}
	} else {
		return object.NewError(fmt.Errorf("message text or options are required"))
	}

	if len(args) > 3 {
		opts, ok := args[3].(*object.Map)
		if !ok {
			return object.NewError(fmt.Errorf("options must be a map"))
		}
		c.processMessageOptions(opts, &options)
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

	switch arg := args[2].(type) {
	case *object.String:
		text, err := object.AsString(arg)
		if err != nil {
			return err
		}
		options = append(options, slack.MsgOptionText(text, false))
	case *object.Map:
		textObj := arg.Get("text")
		if textObj != object.Nil {
			text, err := object.AsString(textObj)
			if err != nil {
				return err
			}
			options = append(options, slack.MsgOptionText(text, false))
		}
		c.processMessageOptions(arg, &options)
	default:
		return object.NewError(fmt.Errorf("third argument must be a string or a map"))
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
		for key, value := range itemMap.Value() {
			switch key {
			case "channel":
				channelStr, err := object.AsString(value)
				if err != nil {
					return err
				}
				itemRef.Channel = channelStr
			case "timestamp":
				timestampStr, err := object.AsString(value)
				if err != nil {
					return err
				}
				itemRef.Timestamp = timestampStr
			case "file":
				fileStr, err := object.AsString(value)
				if err != nil {
					return err
				}
				itemRef.File = fileStr
			case "file_comment":
				commentStr, err := object.AsString(value)
				if err != nil {
					return err
				}
				itemRef.Comment = commentStr
			default:
				return object.NewError(fmt.Errorf("unknown option: %s", key))
			}
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
		for key, value := range itemMap.Value() {
			switch key {
			case "channel":
				channelStr, err := object.AsString(value)
				if err != nil {
					return err
				}
				itemRef.Channel = channelStr
			case "timestamp":
				timestampStr, err := object.AsString(value)
				if err != nil {
					return err
				}
				itemRef.Timestamp = timestampStr
			case "file":
				fileStr, err := object.AsString(value)
				if err != nil {
					return err
				}
				itemRef.File = fileStr
			case "file_comment":
				commentStr, err := object.AsString(value)
				if err != nil {
					return err
				}
				itemRef.Comment = commentStr
			default:
				return object.NewError(fmt.Errorf("unknown option: %s", key))
			}
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

// CreateConversation creates a new channel, either public or private
func (c *Client) CreateConversation(ctx context.Context, args ...object.Object) object.Object {
	if len(args) < 1 {
		return object.NewError(fmt.Errorf("wrong number of arguments: got=%d, want at least 1", len(args)))
	}
	name, argErr := object.AsString(args[0])
	if argErr != nil {
		return argErr
	}
	isPrivate := false
	if len(args) > 1 {
		opts, ok := args[1].(*object.Map)
		if !ok {
			return object.NewError(fmt.Errorf("options must be a map"))
		}
		for key, value := range opts.Value() {
			switch key {
			case "is_private":
				isPrivateBool, err := object.AsBool(value)
				if err != nil {
					return err
				}
				isPrivate = bool(isPrivateBool)
			default:
				return object.NewError(fmt.Errorf("unknown option: %s", key))
			}
		}
	}
	channel, err := c.value.CreateConversationContext(ctx,
		slack.CreateConversationParams{
			ChannelName: name,
			IsPrivate:   isPrivate,
		})
	if err != nil {
		return object.NewError(err)
	}
	return NewChannel(c.value, channel)
}

// GetConversationHistory gets the history of a conversation
func (c *Client) GetConversationHistory(ctx context.Context, args ...object.Object) object.Object {
	if len(args) < 1 {
		return object.NewError(fmt.Errorf("wrong number of arguments: got=%d, want at least 1", len(args)))
	}
	channelID, argErr := object.AsString(args[0])
	if argErr != nil {
		return argErr
	}
	params := &slack.GetConversationHistoryParameters{
		ChannelID: channelID,
		Limit:     100,
	}
	if len(args) > 1 {
		opts, ok := args[1].(*object.Map)
		if !ok {
			return object.NewError(fmt.Errorf("options must be a map"))
		}
		for key, value := range opts.Value() {
			switch key {
			case "oldest":
				oldest, err := object.AsString(value)
				if err != nil {
					return err
				}
				params.Oldest = oldest
			case "latest":
				latest, err := object.AsString(value)
				if err != nil {
					return err
				}
				params.Latest = latest
			case "limit":
				limitInt, err := object.AsInt(value)
				if err != nil {
					return err
				}
				params.Limit = int(limitInt)
			case "cursor":
				cursor, err := object.AsString(value)
				if err != nil {
					return err
				}
				params.Cursor = cursor
			case "inclusive":
				inclusive, err := object.AsBool(value)
				if err != nil {
					return err
				}
				params.Inclusive = bool(inclusive)
			case "include_all_metadata":
				metadata, err := object.AsBool(value)
				if err != nil {
					return err
				}
				params.IncludeAllMetadata = bool(metadata)
			default:
				return object.NewError(fmt.Errorf("unknown option: %s", key))
			}
		}
	}
	return NewMessageIterator(ctx, c.value, params)
}

// GetConversationMembers gets members of a conversation
func (c *Client) GetConversationMembers(ctx context.Context, args ...object.Object) object.Object {
	if len(args) < 1 {
		return object.NewError(fmt.Errorf("wrong number of arguments: got=%d, want at least 1", len(args)))
	}
	channelID, argErr := object.AsString(args[0])
	if argErr != nil {
		return argErr
	}
	options := slack.GetUsersInConversationParameters{
		ChannelID: channelID,
		Limit:     100,
	}
	if len(args) > 1 {
		opts, ok := args[1].(*object.Map)
		if !ok {
			return object.NewError(fmt.Errorf("options must be a map"))
		}
		for key, value := range opts.Value() {
			switch key {
			case "cursor":
				cursorStr, err := object.AsString(value)
				if err != nil {
					return err
				}
				options.Cursor = cursorStr
			case "limit":
				limitInt, err := object.AsInt(value)
				if err != nil {
					return err
				}
				options.Limit = int(limitInt)
			default:
				return object.NewError(fmt.Errorf("unknown option: %s", key))
			}
		}
	}
	return NewConversationMembersIterator(ctx, c.value, &options)
}

func New(client *slack.Client) *Client {
	return &Client{
		value: client,
		base: base{
			typeName:       CLIENT,
			interfaceValue: client,
		},
	}
}

func getTime(t slack.JSONTime) object.Object {
	if t == 0 {
		return object.Nil
	}
	return object.NewTime(t.Time())
}
