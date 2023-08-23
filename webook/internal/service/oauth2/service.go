package oauth2

import "golang.org/x/net/context"

type Service interface {
	AuthURL(ctx context.Context) (string, error)
}
