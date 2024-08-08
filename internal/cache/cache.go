package cache

import "context"

type Cache interface {
	Add(context.Context, string)
	Remove(context.Context, string)
}
