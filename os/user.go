package os

import (
	"os/user"
)

var (
	_ User  = (*UserWrapper)(nil)
	_ Group = (*GroupWrapper)(nil)
)

// UserWrapper wraps the standard library's user.User type to implement the User interface.
type UserWrapper struct {
	*user.User
}

// Uid returns the user ID.
func (u *UserWrapper) Uid() string {
	return u.User.Uid
}

// Gid returns the primary group ID.
func (u *UserWrapper) Gid() string {
	return u.User.Gid
}

// Username returns the username.
func (u *UserWrapper) Username() string {
	return u.User.Username
}

// Name returns the user's name.
func (u *UserWrapper) Name() string {
	return u.User.Name
}

// HomeDir returns the user's home directory.
func (u *UserWrapper) HomeDir() string {
	return u.User.HomeDir
}

// GroupWrapper wraps the standard library's user.Group type to implement the Group interface.
type GroupWrapper struct {
	*user.Group
}

// Gid returns the group ID.
func (g *GroupWrapper) Gid() string {
	return g.Group.Gid
}

// Name returns the group name.
func (g *GroupWrapper) Name() string {
	return g.Group.Name
}

// Current returns the current user.
func Current() (User, error) {
	u, err := user.Current()
	if err != nil {
		return nil, err
	}
	return &UserWrapper{User: u}, nil
}

// LookupUser looks up a user by username.
func LookupUser(username string) (User, error) {
	u, err := user.Lookup(username)
	if err != nil {
		return nil, err
	}
	return &UserWrapper{User: u}, nil
}

// LookupUid looks up a user by user ID.
func LookupUid(uid string) (User, error) {
	u, err := user.LookupId(uid)
	if err != nil {
		return nil, err
	}
	return &UserWrapper{User: u}, nil
}

// LookupGroup looks up a group by name.
func LookupGroup(name string) (Group, error) {
	g, err := user.LookupGroup(name)
	if err != nil {
		return nil, err
	}
	return &GroupWrapper{Group: g}, nil
}

// LookupGid looks up a group by group ID.
func LookupGid(gid string) (Group, error) {
	g, err := user.LookupGroupId(gid)
	if err != nil {
		return nil, err
	}
	return &GroupWrapper{Group: g}, nil
}
