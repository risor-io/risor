package slack

import (
	"encoding/json"
	"fmt"

	"github.com/risor-io/risor/object"
	"github.com/slack-go/slack"
)

const CHANNEL object.Type = "slack.channel"

type Channel struct {
	base
	client *slack.Client
	value  *slack.Channel
}

func (c *Channel) Inspect() string {
	return fmt.Sprintf("slack.channel({id: %q, name: %q})", c.value.ID, c.value.Name)
}

func (c *Channel) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.value)
}

func (c *Channel) Equals(other object.Object) object.Object {
	if other, ok := other.(*Channel); ok {
		return object.NewBool(c.value.ID == other.value.ID)
	}
	return object.False
}

func (c *Channel) GetAttr(name string) (object.Object, bool) {
	switch name {
	case "id":
		return object.NewString(c.value.ID), true
	case "name":
		return object.NewString(c.value.Name), true
	case "is_channel":
		return object.NewBool(c.value.IsChannel), true
	case "is_group":
		return object.NewBool(c.value.IsGroup), true
	case "is_im":
		return object.NewBool(c.value.IsIM), true
	case "is_mpim":
		return object.NewBool(c.value.IsMpIM), true
	case "created":
		return getTime(c.value.Created), true
	case "creator":
		return object.NewString(c.value.Creator), true
	case "is_archived":
		return object.NewBool(c.value.IsArchived), true
	case "is_general":
		return object.NewBool(c.value.IsGeneral), true
	case "unlinked":
		return object.NewInt(int64(c.value.Unlinked)), true
	case "name_normalized":
		return object.NewString(c.value.NameNormalized), true
	case "is_shared":
		return object.NewBool(c.value.IsShared), true
	case "is_ext_shared":
		return object.NewBool(c.value.IsExtShared), true
	case "is_org_shared":
		return object.NewBool(c.value.IsOrgShared), true
	case "is_pending_ext_shared":
		return object.NewBool(c.value.IsPendingExtShared), true
	case "is_member":
		return object.NewBool(c.value.IsMember), true
	case "is_private":
		return object.NewBool(c.value.IsPrivate), true
	case "num_members":
		return object.NewInt(int64(c.value.NumMembers)), true
	case "is_open":
		return object.NewBool(c.value.IsOpen), true
	case "last_read":
		return object.NewString(c.value.LastRead), true
	case "is_global_shared":
		return object.NewBool(c.value.IsGlobalShared), true
	case "is_read_only":
		return object.NewBool(c.value.IsReadOnly), true
	case "priority":
		return object.NewFloat(c.value.Priority), true
	case "user":
		return object.NewString(c.value.User), true
	case "locale":
		return object.NewString(c.value.Locale), true
	case "topic":
		topicMap := map[string]object.Object{
			"value":    object.NewString(c.value.Topic.Value),
			"creator":  object.NewString(c.value.Topic.Creator),
			"last_set": getTime(c.value.Topic.LastSet),
		}
		return object.NewMap(topicMap), true
	case "purpose":
		purposeMap := map[string]object.Object{
			"value":    object.NewString(c.value.Purpose.Value),
			"creator":  object.NewString(c.value.Purpose.Creator),
			"last_set": getTime(c.value.Purpose.LastSet),
		}
		return object.NewMap(purposeMap), true
	case "members":
		membersObjs := make([]object.Object, len(c.value.Members))
		for i, member := range c.value.Members {
			membersObjs[i] = object.NewString(member)
		}
		return object.NewList(membersObjs), true
	case "properties":
		if c.value.Properties != nil {
			return object.NewMap(map[string]object.Object{
				"canvas": object.NewMap(map[string]object.Object{
					"file_id":        object.NewString(c.value.Properties.Canvas.FileId),
					"is_empty":       object.NewBool(c.value.Properties.Canvas.IsEmpty),
					"quip_thread_id": object.NewString(c.value.Properties.Canvas.QuipThreadId),
				}),
			}), true
		}
		return object.Nil, true
	case "connected_team_ids":
		if len(c.value.ConnectedTeamIDs) > 0 {
			teamIDs := make([]object.Object, len(c.value.ConnectedTeamIDs))
			for i, id := range c.value.ConnectedTeamIDs {
				teamIDs[i] = object.NewString(id)
			}
			return object.NewList(teamIDs), true
		}
		return object.NewList([]object.Object{}), true
	case "shared_team_ids":
		if len(c.value.SharedTeamIDs) > 0 {
			teamIDs := make([]object.Object, len(c.value.SharedTeamIDs))
			for i, id := range c.value.SharedTeamIDs {
				teamIDs[i] = object.NewString(id)
			}
			return object.NewList(teamIDs), true
		}
		return object.NewList([]object.Object{}), true
	case "internal_team_ids":
		if len(c.value.InternalTeamIDs) > 0 {
			teamIDs := make([]object.Object, len(c.value.InternalTeamIDs))
			for i, id := range c.value.InternalTeamIDs {
				teamIDs[i] = object.NewString(id)
			}
			return object.NewList(teamIDs), true
		}
		return object.NewList([]object.Object{}), true
	case "context_team_id":
		return object.NewString(c.value.ContextTeamID), true
	case "conversation_host_id":
		return object.NewString(c.value.ConversationHostID), true
	case "json":
		return asMap(c.value), true
	}
	return nil, false
}

func NewChannel(client *slack.Client, channel *slack.Channel) *Channel {
	return &Channel{
		client: client,
		value:  channel,
		base: base{
			typeName:       CHANNEL,
			interfaceValue: channel,
		},
	}
}
