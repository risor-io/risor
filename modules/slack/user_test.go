package slack

import (
	"testing"

	"github.com/risor-io/risor/object"
	"github.com/slack-go/slack"
)

func TestUserObject(t *testing.T) {
	// Create a mock user
	slackUser := &slack.User{
		ID:             "U12345",
		Name:           "testuser",
		TeamID:         "T12345",
		RealName:       "Test User",
		IsAdmin:        true,
		IsBot:          false,
		IsAppUser:      false,
		IsPrimaryOwner: false,
		Profile: slack.UserProfile{
			RealName:      "Test User",
			DisplayName:   "testuser",
			Email:         "test@example.com",
			Image24:       "https://example.com/img24.jpg",
			Image32:       "https://example.com/img32.jpg",
			Image48:       "https://example.com/img48.jpg",
			Image72:       "https://example.com/img72.jpg",
			Image192:      "https://example.com/img192.jpg",
			Image512:      "https://example.com/img512.jpg",
			ImageOriginal: "https://example.com/img.jpg",
		},
	}

	// Create a user object
	user := NewUser(nil, slackUser)

	// Test basic properties
	if user.Type() != USER {
		t.Errorf("Expected Type to be %s, got %s", USER, user.Type())
	}

	if !user.IsTruthy() {
		t.Error("Expected IsTruthy to be true")
	}

	// Test GetAttr
	testCases := []struct {
		attr     string
		expected object.Object
	}{
		{"id", object.NewString("U12345")},
		{"name", object.NewString("testuser")},
		{"team_id", object.NewString("T12345")},
		{"real_name", object.NewString("Test User")},
		{"is_admin", object.True},
		{"is_bot", object.False},
		{"is_app_user", object.False},
		{"is_primary_owner", object.False},
	}

	for _, tc := range testCases {
		val, ok := user.GetAttr(tc.attr)
		if !ok {
			t.Errorf("Expected GetAttr(%q) to return a value", tc.attr)
			continue
		}

		if val.Type() != tc.expected.Type() {
			t.Errorf("Expected GetAttr(%q) to return type %s, got %s", tc.attr, tc.expected.Type(), val.Type())
			continue
		}

		if val.Inspect() != tc.expected.Inspect() {
			t.Errorf("Expected GetAttr(%q) to return %s, got %s", tc.attr, tc.expected.Inspect(), val.Inspect())
		}
	}

	// Test profile attribute
	profile, ok := user.GetAttr("profile")
	if !ok {
		t.Error("Expected GetAttr('profile') to return a value")
	} else {
		profileMap, ok := profile.(*object.Map)
		if !ok {
			t.Errorf("Expected profile to be a map, got %T", profile)
		} else {
			email := profileMap.Get("email")
			if email.Inspect() != `"test@example.com"` {
				t.Errorf("Expected profile.email to be 'test@example.com', got %s", email.Inspect())
			}

			realName := profileMap.Get("real_name")
			if realName.Inspect() != `"Test User"` {
				t.Errorf("Expected profile.real_name to be 'Test User', got %s", realName.Inspect())
			}
		}
	}
}
