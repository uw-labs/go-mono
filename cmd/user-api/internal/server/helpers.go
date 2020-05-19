package server

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/uw-labs/go-mono/cmd/user-api/internal/repo"
	userservicepb "github.com/uw-labs/go-mono/proto/gen/go/uwlabs/users/service/v1"
)

func orderProtoToInternal(order *userservicepb.SortOrder) (*repo.SortOrder, error) {
	if order == nil {
		return nil, nil
	}

	o := repo.SortOrder{
		Descending: order.GetDescending(),
	}

	switch order.GetBy() {
	case userservicepb.OrderBy_ORDER_BY_NONE:
		return nil, nil
	case userservicepb.OrderBy_ORDER_BY_NAME:
		o.By = repo.OrderByName
	case userservicepb.OrderBy_ORDER_BY_CREATE_TIME:
		o.By = repo.OrderByCreateTime
	default:
		return nil, fmt.Errorf("unknown sort order: %q", order.GetBy().String())
	}

	return &o, nil
}

const (
	prefix     = "Basic "
	authHeader = "Authorization"
)

// authenticate checks that the user is authenticated.
func (s Server) authenticate(ctx context.Context) (bool, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	hdrs := md.Get(authHeader)
	if len(hdrs) == 0 || hdrs[0] == "" {
		return false, status.Error(codes.Unauthenticated, `missing "Authorization" header`)
	}

	// Just pick the first Authorization header.
	auth := hdrs[0]

	if !strings.HasPrefix(auth, prefix) {
		return false, status.Error(codes.Unauthenticated, `missing "Basic " prefix in "Authorization" header`)
	}

	c, err := base64.StdEncoding.DecodeString(auth[len(prefix):])
	if err != nil {
		return false, status.Error(codes.Unauthenticated, `invalid base64 in header`)
	}

	cs := string(c)
	ci := strings.IndexByte(cs, ':')
	if ci < 0 {
		return false, status.Error(codes.Unauthenticated, `invalid basic auth format`)
	}

	user, password := cs[:ci], cs[ci+1:]
	if user != s.Admin.Username || password != s.Admin.Password {
		return false, status.Error(codes.Unauthenticated, "invalid user or password")
	}

	return true, nil
}
