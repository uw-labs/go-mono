package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"

	"github.com/Masterminds/squirrel"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	bindata "github.com/golang-migrate/migrate/v4/source/go_bindata"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/logrusadapter"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/luna-duclos/instrumentedsql"
	"github.com/luna-duclos/instrumentedsql/opentracing"
	"github.com/sirupsen/logrus"

	"github.com/uw-labs/go-mono/cmd/user-api/internal/repo/migrations"
)

//go:generate go-bindata -pkg migrations -prefix migrations -nometadata -ignore bindata -o ./migrations/bindata.go ./migrations

// Repository is used to store and fetch users from a database.
type Repository struct {
	db *sql.DB
	sb squirrel.StatementBuilderType
}

// NewRepository creates a new, database-backed repository.
func NewRepository(dbURL string, logger *logrus.Logger) (*Repository, error) {
	parsed, err := url.Parse(dbURL)
	if err != nil {
		return nil, fmt.Errorf("invalid database URL: %w", err)
	}

	// Work around PGX not parsing URLs without a host correctly.
	// https://github.com/jackc/pgconn/issues/19
	if parsed.Hostname() == "" {
		parsed.Host = "localhost" + parsed.Host
	}

	connConfig, err := pgx.ParseConfig(parsed.String())
	if err != nil {
		return nil, fmt.Errorf("couldn't parse DB DSN: %w", err)
	}

	connConfig.Logger = logrusadapter.NewLogger(logger)

	connStr := stdlib.RegisterConnConfig(connConfig)
	drv := instrumentedsql.WrapDriver(
		stdlib.GetDefaultDriver(),
		instrumentedsql.WithTracer(opentracing.NewTracer(false)),
		instrumentedsql.WithOmitArgs(),
		instrumentedsql.WithOpsExcluded(instrumentedsql.OpSQLRowsNext),
	)
	cnctr, err := drv.OpenConnector(connStr)
	if err != nil {
		return nil, fmt.Errorf("open driver: %w", err)
	}

	sqlDB := sql.OpenDB(cnctr)
	if err = sqlDB.Ping(); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("connect to the database: %w", err)
	}

	source, err := bindata.WithInstance(bindata.Resource(migrations.AssetNames(), migrations.Asset))
	if err != nil {
		return nil, fmt.Errorf("creating bindata migration: %w", err)
	}

	target, err := postgres.WithInstance(sqlDB, new(postgres.Config))
	if err != nil {
		return nil, fmt.Errorf("creating postgres migration: %w", err)
	}

	m, err := migrate.NewWithInstance("bindata", source, "postgres", target)
	if err != nil {
		return nil, fmt.Errorf("creating migrater: %w", err)
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf("migrating up: %w", err)
	}

	return &Repository{
		db: sqlDB,
		sb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).RunWith(sqlDB),
	}, nil
}

// Close releases the resources of the Repository.
func (r Repository) Close() error {
	return r.db.Close()
}

// CreateUser creates a new user in the repository.
func (r Repository) CreateUser(ctx context.Context, name string) (User, error) {
	iq := r.sb.Insert(
		"users",
	).SetMap(map[string]interface{}{
		"name": name,
	}).Suffix(
		"RETURNING id, create_time",
	)

	user := User{
		Name: name,
	}
	err := iq.QueryRowContext(ctx).Scan(&user.ID, &user.CreateTime)
	if err != nil {
		return User{}, fmt.Errorf("create new user: %w", err)
	}

	return user, nil
}

// GetUser retrieves any information held for a given user ID.
func (r Repository) GetUser(ctx context.Context, id string) (User, error) {
	var fID pgtype.UUID
	err := fID.Set(id)
	if err != nil {
		return User{}, fmt.Errorf("parse ID as UUID: %w", err)
	}
	q := r.sb.Select(
		"id",
		"name",
		"create_time",
	).From(
		"users",
	).Where(
		squirrel.Eq{"id": fID},
	)

	var user User
	err = q.QueryRowContext(ctx).Scan(
		&user.ID,
		&user.Name,
		&user.CreateTime,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, ErrUserNotFound
		}
		return User{}, fmt.Errorf("retrieve user: %w", err)
	}

	return user, nil
}

// ListUsers will return all users matching the filters.
func (r Repository) ListUsers(ctx context.Context, names []string, order *SortOrder) (_ []User, err error) {
	q := r.sb.Select(
		"id",
		"name",
		"create_time",
	).From(
		"users",
	)

	if order != nil {
		q = q.OrderBy(
			order.SQL(),
		)
	}

	if len(names) > 0 {
		q = q.Where(
			squirrel.Eq{"name": names},
		)
	}

	rows, err := q.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("list users query: %w", err)
	}
	defer func() {
		cErr := rows.Close()
		if err == nil && cErr != nil {
			err = fmt.Errorf("closing users rows: %w", cErr)
		}
	}()

	var users []User
	for rows.Next() {
		var user User
		err = rows.Scan(
			&user.ID,
			&user.Name,
			&user.CreateTime,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning user row: %w", err)
		}
		users = append(users, user)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("iterating user rows: %w", err)
	}

	return users, nil
}
