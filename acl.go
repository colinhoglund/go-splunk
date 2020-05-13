package splunk

import (
	"errors"
	"net/url"
)

// ACL represents an ACL for an API resource
type ACL struct {
	Owner      string
	Sharing    string
	PermsRead  string
	PermsWrite string
}

// NewACL returns a new ACL
func NewACL(owner, sharing, read, write string) (*ACL, error) {
	if owner == "" {
		return nil, errors.New("error creating acl: owner is a required field")
	}
	if sharing == "" {
		return nil, errors.New("error creating acl: sharing is a required field")
	}

	return &ACL{owner, sharing, read, write}, nil
}

// Encode returns the ACL as an encoded query string
func (acl *ACL) Encode() string {
	form := url.Values{}
	form.Set("owner", acl.Owner)
	form.Set("sharing", acl.Sharing)
	if acl.PermsRead != "" {
		form.Set("perms.read", acl.PermsRead)
	}
	if acl.PermsWrite != "" {
		form.Set("perms.write", acl.PermsWrite)
	}

	return form.Encode()
}
