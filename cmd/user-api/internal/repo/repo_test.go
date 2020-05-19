package repo_test

import (
	"context"
	"database/sql"
	"errors"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	bindata "github.com/golang-migrate/migrate/v4/source/go_bindata"
	"github.com/google/go-cmp/cmp"
	"github.com/sirupsen/logrus"
	"github.com/uw-labs/podrick"
	_ "github.com/uw-labs/podrick/runtimes/docker" // register docker runtime
	podricklogger "logur.dev/adapter/logrus"

	"github.com/uw-labs/go-mono/cmd/user-api/internal/repo"
	"github.com/uw-labs/go-mono/cmd/user-api/internal/repo/migrations"
	pkgctx "github.com/uw-labs/go-mono/pkg/context"
)

var (
	logger *logrus.Logger

	pgURL *url.URL
)

func TestMain(m *testing.M) {
	code := 0
	defer func() {
		os.Exit(code)
	}()
	ctx := pkgctx.WithSignalHandler(context.Background())

	logger = logrus.New()
	logger.Formatter = &logrus.TextFormatter{
		TimestampFormat: time.StampMilli,
		FullTimestamp:   true,
	}

	ctr, err := podrick.StartContainer(ctx, "postgres", "12-alpine", "5432",
		podrick.WithEnv([]string{
			"POSTGRES_HOST_AUTH_METHOD=trust", // https://github.com/docker-library/postgres/issues/681
		}),
		podrick.WithLivenessCheck(func(address string) error {
			dbURL, err := url.Parse("postgresql://postgres@" + address + "/postgres?sslmode=disable")
			if err != nil {
				return err
			}
			db, err := sql.Open("pgx", dbURL.String())
			if err != nil {
				return err
			}
			defer db.Close()
			return db.Ping()
		}),
		podrick.WithLogger(podricklogger.New(logger)),
	)
	if err != nil {
		logger.Println("Failed to start database container", err)
		return
	}
	defer func() {
		err = ctr.Close(context.Background())
		if err != nil {
			logger.Println("Failed to stop database container", err)
			return
		}
	}()

	pgURL, err = url.Parse("postgresql://postgres@" + ctr.Address() + "/postgres?sslmode=disable")
	if err != nil {
		logger.Println("Failed to parse container address", err)
		return
	}

	code = m.Run()
}

func TestRepository(t *testing.T) {
	ctx := pkgctx.WithSignalHandler(context.Background())

	db, err := sql.Open("pgx", pgURL.String())
	if err != nil {
		t.Fatalf("unexpected error opening test database: %v", err)
	}

	tests := []struct {
		Name string
		Run  func(r *repo.Repository) func(t *testing.T)
	}{
		{
			"Can create a new user & get its details",
			func(r *repo.Repository) func(t *testing.T) {
				return func(t *testing.T) {
					user, err := r.CreateUser(ctx, "Alice")
					if err != nil {
						t.Fatalf("unexpected error creating user: %v", err)
					}

					if user.Name != "Alice" {
						t.Errorf("unexpected name; got %q, wanted %q", user.Name, "Alice")
					}

					if time.Until(user.CreateTime) > time.Second {
						t.Errorf("CreateTime was more than 1 second in the past")
					}

					if user.ID == "" {
						t.Errorf("ID was not set")
					}

					user2, err := r.GetUser(ctx, user.ID)
					if err != nil {
						t.Fatalf("unexpected error getting user: %v", err)
					}

					if diff := cmp.Diff(user, user2); diff != "" {
						t.Errorf("get returned different user than create:\n%s", diff)
					}
				}
			},
		},
		{
			"Can list users",
			func(r *repo.Repository) func(t *testing.T) {
				return func(t *testing.T) {
					user1, err := r.CreateUser(ctx, "Alice")
					if err != nil {
						t.Fatalf("unexpected error creating user: %v", err)
					}

					user2, err := r.CreateUser(ctx, "Bob")
					if err != nil {
						t.Fatalf("unexpected error creating user: %v", err)
					}

					users, err := r.ListUsers(ctx, nil, &repo.SortOrder{By: repo.OrderByName})
					if err != nil {
						t.Fatalf("unexpected error listing users: %v", err)
					}

					for i, user := range []repo.User{user1, user2} {
						if diff := cmp.Diff(user, users[i]); diff != "" {
							t.Errorf("list returned different user than create:\n%s", diff)
						}
					}
				}
			},
		},
	}

	source, err := bindata.WithInstance(bindata.Resource(migrations.AssetNames(), migrations.Asset))
	if err != nil {
		t.Fatalf("unexpected error creating bindata migration: %v", err)
	}

	target, err := postgres.WithInstance(db, new(postgres.Config))
	if err != nil {
		t.Fatalf("unexpected error creating postgres migration: %v", err)
	}

	m, err := migrate.NewWithInstance("bindata", source, "postgres", target)
	if err != nil {
		t.Fatalf("unexpected error creating migrater: %v", err)
	}

	for _, tc := range tests {
		// Ensure database is clean
		err = m.Down()
		if err != nil && !errors.Is(err, migrate.ErrNoChange) {
			t.Fatalf("unexpected error clearing database: %v", err)
		}

		r, err := repo.NewRepository(pgURL.String(), logger)
		if err != nil {
			t.Fatalf("unexpected error creating repo: %v", err)
		}

		t.Run(tc.Name, tc.Run(r))

		err = r.Close()
		if err != nil {
			t.Fatalf("unexpected error closing repo: %v", err)
		}
	}
}
