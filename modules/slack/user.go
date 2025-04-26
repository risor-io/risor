package slack

import (
	"context"
	"fmt"

	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/op"
	"github.com/slack-go/slack"
)

// USER type constant
const USER object.Type = "slack.user"

// User wraps a slack.User and implements the object.Object interface
type User struct {
	value  *slack.User
	client *slack.Client
}

func (u *User) Type() object.Type {
	return USER
}

func (u *User) Inspect() string {
	return fmt.Sprintf("slack.user({id: %q, name: %q})", u.value.ID, u.value.Name)
}

func (u *User) Interface() interface{} {
	return u.value
}

func (u *User) Equals(other object.Object) object.Object {
	if o, ok := other.(*User); ok {
		return object.NewBool(u.value.ID == o.value.ID)
	}
	return object.False
}

func (u *User) GetAttr(name string) (object.Object, bool) {
	switch name {
	case "id":
		return object.NewString(u.value.ID), true
	case "team_id":
		return object.NewString(u.value.TeamID), true
	case "name":
		return object.NewString(u.value.Name), true
	case "deleted":
		return object.NewBool(u.value.Deleted), true
	case "color":
		return object.NewString(u.value.Color), true
	case "real_name":
		return object.NewString(u.value.RealName), true
	case "tz":
		return object.NewString(u.value.TZ), true
	case "tz_label":
		return object.NewString(u.value.TZLabel), true
	case "tz_offset":
		return object.NewInt(int64(u.value.TZOffset)), true
	case "profile":
		return convertProfileToMap(&u.value.Profile), true
	case "is_admin":
		return object.NewBool(u.value.IsAdmin), true
	case "is_owner":
		return object.NewBool(u.value.IsOwner), true
	case "is_primary_owner":
		return object.NewBool(u.value.IsPrimaryOwner), true
	case "is_restricted":
		return object.NewBool(u.value.IsRestricted), true
	case "is_ultra_restricted":
		return object.NewBool(u.value.IsUltraRestricted), true
	case "is_bot":
		return object.NewBool(u.value.IsBot), true
	case "is_app_user":
		return object.NewBool(u.value.IsAppUser), true
	case "updated":
		if !u.value.Updated.Time().IsZero() {
			return object.NewTime(u.value.Updated.Time()), true
		}
		return object.Nil, true
	case "has_2fa":
		return object.NewBool(u.value.Has2FA), true
	case "two_factor_type":
		if u.value.TwoFactorType != nil {
			return object.NewString(*u.value.TwoFactorType), true
		}
		return object.Nil, true
	case "has_files":
		return object.NewBool(u.value.HasFiles), true
	case "presence":
		return object.NewString(u.value.Presence), true
	case "locale":
		return object.NewString(u.value.Locale), true
	case "is_stranger":
		return object.NewBool(u.value.IsStranger), true
	case "is_invited_user":
		return object.NewBool(u.value.IsInvitedUser), true
	case "dm":
		return object.NewBuiltin("dm", u.dm), true
	}
	return nil, false
}

func (u *User) SetAttr(name string, value object.Object) error {
	return fmt.Errorf("type error: cannot set %q on slack.user object", name)
}

func (u *User) IsTruthy() bool {
	return true
}

func (u *User) RunOperation(opType op.BinaryOpType, right object.Object) object.Object {
	return object.Errorf("type error: unsupported operation for slack.user")
}

func (u *User) Cost() int {
	return 0
}

// dm sends a direct message to the user
func (u *User) dm(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError(fmt.Errorf("wrong number of arguments: got=%d, want=1", len(args)))
	}

	text, err := object.AsString(args[0])
	if err != nil {
		return err
	}

	_, _, msgErr := u.client.PostMessageContext(
		ctx,
		u.value.ID,
		slack.MsgOptionText(text, false),
	)

	if msgErr != nil {
		return object.NewError(msgErr)
	}

	return object.Nil
}

// NewUser creates a new User object
func NewUser(client *slack.Client, user *slack.User) *User {
	return &User{
		client: client,
		value:  user,
	}
}

// Helper to convert UserProfile to a map
func convertProfileToMap(profile *slack.UserProfile) *object.Map {
	profileMap := map[string]object.Object{
		"real_name":               object.NewString(profile.RealName),
		"real_name_normalized":    object.NewString(profile.RealNameNormalized),
		"display_name":            object.NewString(profile.DisplayName),
		"display_name_normalized": object.NewString(profile.DisplayNameNormalized),
		"email":                   object.NewString(profile.Email),
		"first_name":              object.NewString(profile.FirstName),
		"last_name":               object.NewString(profile.LastName),
		"phone":                   object.NewString(profile.Phone),
		"skype":                   object.NewString(profile.Skype),
		"title":                   object.NewString(profile.Title),
		"team":                    object.NewString(profile.Team),
		"status_text":             object.NewString(profile.StatusText),
		"status_emoji":            object.NewString(profile.StatusEmoji),
		"bot_id":                  object.NewString(profile.BotID),
		"image_24":                object.NewString(profile.Image24),
		"image_32":                object.NewString(profile.Image32),
		"image_48":                object.NewString(profile.Image48),
		"image_72":                object.NewString(profile.Image72),
		"image_192":               object.NewString(profile.Image192),
		"image_512":               object.NewString(profile.Image512),
		"image_original":          object.NewString(profile.ImageOriginal),
	}

	return object.NewMap(profileMap)
}
