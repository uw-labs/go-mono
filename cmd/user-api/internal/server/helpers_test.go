package server

import (
	"context"
	"encoding/base64"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestAuthenticate(t *testing.T) {
	s := &Server{}

	tests := []struct {
		Name string
		Run  func(*testing.T)
	}{
		{
			Name: "It accepts a correctly formatted header",
			Run: func(t *testing.T) {
				u, pw := "user", "password"
				creds := u + ":" + pw
				value := "Basic " + base64.StdEncoding.EncodeToString([]byte(creds))
				ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("Authorization", value))
				s := &Server{
					Admin: &User{
						Username: u,
						Password: pw,
					},
				}
				ok, err := s.authenticate(ctx)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if !ok {
					t.Error("expected authentication to be successful")
				}
			},
		},
		{
			Name: "It accepts a header regardless of casing",
			Run: func(t *testing.T) {
				u, pw := "user", "password"
				creds := u + ":" + pw
				value := "Basic " + base64.StdEncoding.EncodeToString([]byte(creds))
				ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", value))
				s := &Server{
					Admin: &User{
						Username: u,
						Password: pw,
					},
				}
				ok, err := s.authenticate(ctx)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if !ok {
					t.Error("expected authentication to be successful")
				}
			},
		},
		{
			Name: "It rejects an empty context",
			Run: func(t *testing.T) {
				ok, err := s.authenticate(context.Background())
				if ok {
					t.Error("expected authentication to be unsuccessful")
				}
				st := status.Convert(err)
				if st.Code() != codes.Unauthenticated {
					t.Errorf("expected code %T, got %T", codes.Unauthenticated, st.Code())
				}
				if st.Message() != `missing "Authorization" header` {
					t.Errorf("expected message %q, got %q", `missing "Authorization" header`, st.Message())
				}
			},
		},
		{
			Name: "It rejects a context with the wrong header",
			Run: func(t *testing.T) {
				ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("key", "value"))
				ok, err := s.authenticate(ctx)
				if ok {
					t.Error("expected authentication to be unsuccessful")
				}
				st := status.Convert(err)
				if st.Code() != codes.Unauthenticated {
					t.Errorf("expected code %T, got %T", codes.Unauthenticated, st.Code())
				}
				if st.Message() != `missing "Authorization" header` {
					t.Errorf("expected message %q, got %q", `missing "Authorization" header`, st.Message())
				}
			},
		},
		{
			Name: "It rejects a header with incorrect formatting",
			Run: func(t *testing.T) {
				ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "nope"))
				ok, err := s.authenticate(ctx)
				if ok {
					t.Error("expected authentication to be unsuccessful")
				}
				st := status.Convert(err)
				if st.Code() != codes.Unauthenticated {
					t.Errorf("expected code %T, got %T", codes.Unauthenticated, st.Code())
				}
				if st.Message() != `missing "Basic " prefix in "Authorization" header` {
					t.Errorf("expected message %q, got %q", `missing "Basic " prefix in "Authorization" header`, st.Message())
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, test.Run)
	}
}
