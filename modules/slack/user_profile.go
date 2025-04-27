package slack

import (
	"encoding/json"
	"fmt"

	"github.com/risor-io/risor/object"
	"github.com/slack-go/slack"
)

const USER_PROFILE object.Type = "slack.user_profile"

type UserProfile struct {
	base
	value *slack.UserProfile
}

func (p *UserProfile) Inspect() string {
	return fmt.Sprintf("slack.user_profile({real_name: %q, display_name: %q})",
		p.value.RealName, p.value.DisplayName)
}

func (p *UserProfile) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.value)
}

func (p *UserProfile) Equals(other object.Object) object.Object {
	if o, ok := other.(*UserProfile); ok {
		// Compare relevant fields for equality
		return object.NewBool(p.value.RealName == o.value.RealName &&
			p.value.DisplayName == o.value.DisplayName &&
			p.value.Email == o.value.Email)
	}
	return object.False
}

func (p *UserProfile) GetAttr(name string) (object.Object, bool) {
	switch name {
	case "real_name":
		return object.NewString(p.value.RealName), true
	case "real_name_normalized":
		return object.NewString(p.value.RealNameNormalized), true
	case "display_name":
		return object.NewString(p.value.DisplayName), true
	case "display_name_normalized":
		return object.NewString(p.value.DisplayNameNormalized), true
	case "email":
		return object.NewString(p.value.Email), true
	case "first_name":
		return object.NewString(p.value.FirstName), true
	case "last_name":
		return object.NewString(p.value.LastName), true
	case "phone":
		return object.NewString(p.value.Phone), true
	case "skype":
		return object.NewString(p.value.Skype), true
	case "title":
		return object.NewString(p.value.Title), true
	case "team":
		return object.NewString(p.value.Team), true
	case "status_text":
		return object.NewString(p.value.StatusText), true
	case "status_emoji":
		return object.NewString(p.value.StatusEmoji), true
	case "bot_id":
		return object.NewString(p.value.BotID), true
	case "image_24":
		return object.NewString(p.value.Image24), true
	case "image_32":
		return object.NewString(p.value.Image32), true
	case "image_48":
		return object.NewString(p.value.Image48), true
	case "image_72":
		return object.NewString(p.value.Image72), true
	case "image_192":
		return object.NewString(p.value.Image192), true
	case "image_512":
		return object.NewString(p.value.Image512), true
	case "image_original":
		return object.NewString(p.value.ImageOriginal), true
	}
	return nil, false
}

func NewUserProfile(profile *slack.UserProfile) *UserProfile {
	return &UserProfile{
		value: profile,
		base: base{
			typeName:       USER_PROFILE,
			interfaceValue: profile,
		},
	}
}
