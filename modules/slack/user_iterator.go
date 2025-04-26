package slack

import (
	"context"
	"fmt"

	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/op"
	"github.com/slack-go/slack"
)

// UserIterator implements a Risor iterator that handles paginated user results
type UserIterator struct {
	client          *slack.Client
	ctx             context.Context
	users           []slack.User
	currentIndex    int
	currentEntry    *object.Entry
	cursor          string
	hasMore         bool
	limit           int
	includePresence bool
	teamID          string
}

func (i *UserIterator) Type() object.Type {
	return "slack.user_iterator"
}

func (i *UserIterator) Inspect() string {
	return "slack.user_iterator()"
}

func (i *UserIterator) Interface() interface{} {
	return i
}

func (i *UserIterator) Equals(other object.Object) object.Object {
	if i == other {
		return object.True
	}
	return object.False
}

func (i *UserIterator) GetAttr(name string) (object.Object, bool) {
	return nil, false
}

func (i *UserIterator) SetAttr(name string, value object.Object) error {
	return fmt.Errorf("cannot set attribute on user iterator")
}

func (i *UserIterator) IsTruthy() bool {
	return true
}

func (i *UserIterator) RunOperation(opType op.BinaryOpType, right object.Object) object.Object {
	return object.Errorf("operation not supported on user iterator")
}

func (i *UserIterator) Cost() int {
	return 1
}

func (i *UserIterator) Next(ctx context.Context) (object.Object, bool) {
	// If we've reached the end of the current batch
	if i.currentIndex >= len(i.users) {
		// Check if there are more users to fetch
		if i.hasMore && i.cursor != "" {
			// Fetch the next batch
			options := []slack.GetUsersOption{}

			if i.limit > 0 {
				options = append(options, slack.GetUsersOptionLimit(i.limit))
			}

			if i.includePresence {
				options = append(options, slack.GetUsersOptionPresence(true))
			}

			if i.teamID != "" {
				options = append(options, slack.GetUsersOptionTeamID(i.teamID))
			}

			users, err := i.client.GetUsersContext(i.ctx, options...)
			if err != nil {
				return object.NewError(err), false
			}

			i.users = users
			i.cursor = "" // No cursor in this API, so we can't paginate further
			i.currentIndex = 0
			i.hasMore = false // No pagination in this API

			// If no users were returned, we're done
			if len(users) == 0 {
				return object.Nil, false
			}
		} else {
			// No more users to fetch
			return object.Nil, false
		}
	}

	// Get the current user and create an entry
	user := i.users[i.currentIndex]
	i.currentEntry = object.NewEntry(object.NewInt(int64(i.currentIndex)), convertUserToMap(user))
	i.currentIndex++

	return i.currentEntry.Value(), true
}

func (i *UserIterator) Entry() (object.IteratorEntry, bool) {
	if i.currentEntry == nil {
		return nil, false
	}
	return i.currentEntry, true
}
