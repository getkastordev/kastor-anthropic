package provider

import (
	"context"

	protocol "github.com/weirdGuy/kastor/protocol/v1"
)

type Object = protocol.Object
type Resource = protocol.Resource
type AttrDiff = protocol.AttrDiff
type Status = protocol.Status
type Check = protocol.Check

const (
	StatusOK      = protocol.StatusOK
	StatusFailed  = protocol.StatusFailed
	StatusUnknown = protocol.StatusUnknown
)

type Provider interface {
	Read(context.Context, string) (Object, bool, error)
	Create(context.Context, *Resource) (string, error)
	Update(context.Context, string, *Resource) error
	Delete(context.Context, string) error
	Diff(*Resource, Object) ([]AttrDiff, error)
}

type Checker interface {
	Check(context.Context, *Resource, Object) ([]Check, error)
}
