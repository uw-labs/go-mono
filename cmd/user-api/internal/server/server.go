package server

import (
	"context"
	"errors"

	"github.com/golang/protobuf/ptypes"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/uw-labs/go-mono/cmd/user-api/internal/repo"
	usersservicepb "github.com/uw-labs/go-mono/proto/gen/go/uwlabs/users/service/v1"
	userspb "github.com/uw-labs/go-mono/proto/gen/go/uwlabs/users/v1"
)

type (
	// Repository allows storing and fetching of users and deployments.
	Repository interface {
		GetUser(ctx context.Context, id string) (repo.User, error)
		CreateUser(ctx context.Context, name string) (repo.User, error)
		ListUsers(ctx context.Context, names []string, order *repo.SortOrder) ([]repo.User, error)
	}

	// Server implements the gRPC server interface
	Server struct {
		Repo   Repository
		Admin  *User
		Logger *logrus.Logger
	}

	// User describes an authenticated user
	User struct {
		Username string
		Password string
	}
)

// CreateUser creates a new user in the database.
func (s *Server) CreateUser(ctx context.Context, req *usersservicepb.CreateUserRequest) (*usersservicepb.CreateUserResponse, error) {
	ok, err := s.authenticate(ctx)
	if err != nil {
		// Returns a status error
		return nil, err
	}
	if !ok {
		return nil, status.Error(codes.PermissionDenied, "not permitted to create users")
	}

	user, err := s.Repo.CreateUser(ctx, req.GetName())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	ct, err := ptypes.TimestampProto(user.CreateTime)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to convert create time to proto type")
	}

	return &usersservicepb.CreateUserResponse{
		User: &userspb.User{
			Name:       user.Name,
			Id:         user.ID,
			CreateTime: ct,
		},
	}, nil
}

// GetUser returns the user corresponding to the ID, if found.
func (s *Server) GetUser(ctx context.Context, req *usersservicepb.GetUserRequest) (*usersservicepb.GetUserResponse, error) {
	user, err := s.Repo.GetUser(ctx, req.GetUserId())
	if err != nil {
		if errors.Is(err, repo.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user could not be found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	ct, err := ptypes.TimestampProto(user.CreateTime)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to convert create time to proto type")
	}

	return &usersservicepb.GetUserResponse{
		User: &userspb.User{
			Name:       user.Name,
			Id:         user.ID,
			CreateTime: ct,
		},
	}, nil
}

// ListUsers returns all the users matching the filters.
func (s *Server) ListUsers(ctx context.Context, req *usersservicepb.ListUsersRequest) (res *usersservicepb.ListUsersResponse, err error) {
	order, err := orderProtoToInternal(req.GetSortOrder())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	users, err := s.Repo.ListUsers(
		ctx,
		req.GetNames(),
		order,
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	pbUsers := []*userspb.User{}
	for _, user := range users {
		ct, err := ptypes.TimestampProto(user.CreateTime)
		if err != nil {
			return nil, status.Error(codes.Internal, "failed to convert create time to proto type")
		}

		pbUser := &userspb.User{
			Name:       user.Name,
			Id:         user.ID,
			CreateTime: ct,
		}
		pbUsers = append(pbUsers, pbUser)
	}

	return &usersservicepb.ListUsersResponse{
		Users: pbUsers,
	}, nil
}
