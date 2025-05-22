package tools

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	database "cloud.google.com/go/spanner/admin/database/apiv1"
	"cloud.google.com/go/spanner/admin/database/apiv1/databasepb"
	instance "cloud.google.com/go/spanner/admin/instance/apiv1"
	"cloud.google.com/go/spanner/admin/instance/apiv1/instancepb"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/moby/moby/pkg/namesgenerator"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	_ "github.com/googleapis/go-sql-spanner"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func GetDB(spanner bool) (*sql.DB, func(), error) {
	dsn, terminate, err := SetupSpannerDSN(spanner)
	if err != nil {
		return nil, nil, err
	}

	var db *sql.DB
	if spanner {
		db, err = sql.Open("spanner", dsn)
		if err != nil {
			return nil, nil, err
		}
	} else {
		db, err = sql.Open("pgx/v4", dsn)
		if err != nil {
			return nil, nil, err
		}
	}
	return db, terminate, nil
}

// SetupSpannerDB starts a spanner adapter container and returns the dsn to connect to it.
//
// DBOptions is kept mostly for compatibility with the postgres container setup function.
// The only option that will actually be used is SpannerDialect.
func SetupSpannerDSN(spanner bool) (dsn string, terminate func(), err error) {
	switch spanner {
	case false:
		dsn, terminate, err = setupSpannerForPGAdapter()
	case true:
		dsn, terminate, err = setupSpannerForGoogleSQL()
	}

	return dsn, terminate, err
}

// SetupSpannerDB starts a spanner pg adapter container and returns the dsn to connect to it.
//
// DBOptions is kept for compatibility with the postgres container setup function, but
// it is not used.
func setupSpannerForPGAdapter() (dsn string, terminate func(), err error) {
	ctx := context.Background()
	opts := DBOptions{PostgresUser: "postgres", PostgresPassword: "postgres", DBName: "test-database"}

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "gcr.io/cloud-spanner-pg-adapter/pgadapter-emulator",
			ExposedPorts: []string{"5432"},
			Name:         namesgenerator.GetRandomName(42),
		},
		Started: true,
	}

	testcontainers.WithWaitStrategyAndDeadline(
		30*time.Second,
		wait.ForSQL("5432", "pgx/v4", getDSN(opts)).
			WithPollInterval(1*time.Second).
			WithQuery("SELECT 1"),
	).Customize(&req)

	dbContainer, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		if dbContainer != nil {
			dbContainer.Terminate(ctx)
		}
		return "", func() {}, fmt.Errorf("could not start spanner pg adapter container: %w", err)
	}

	host, err := dbContainer.Host(ctx)
	if err != nil {
		dbContainer.Terminate(ctx)
		return "", func() {}, fmt.Errorf("could not get spanner pg adapter container host: %w", err)
	}

	port, err := dbContainer.MappedPort(ctx, "5432")
	if err != nil {
		dbContainer.Terminate(ctx)
		return "", func() {}, fmt.Errorf("could not get spanner pg adapter container port: %w", err)
	}

	dsn = getDSN(opts)(host, port)
	return dsn, func() { dbContainer.Terminate(ctx) }, nil
}

const (
	testProject  = "test-project"
	testInstance = "test-instance"
	SpannerDB    = "test-database"
	listenPort   = "9010/tcp"
)

func setupSpannerForGoogleSQL() (dsn string, terminate func(), err error) {
	ctx := context.Background()
	dsn = fmt.Sprintf("projects/%s/instances/%s/databases/%s?x-clean-statements=true", testProject, testInstance, SpannerDB)
	terminate = func() {}

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "gcr.io/cloud-spanner-emulator/emulator",
			ExposedPorts: []string{listenPort},
			Name:         namesgenerator.GetRandomName(42),
			WaitingFor:   wait.ForLog("gRPC server listening").WithStartupTimeout(120 * time.Second),
			HostConfigModifier: func(hostConfig *container.HostConfig) {
				hostConfig.AutoRemove = true
			},
		},
		Started: true,
	}

	spannerEmulator, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return dsn, terminate, err
	}

	terminate = func() { _ = spannerEmulator.Terminate(ctx) }

	host, err := spannerEmulator.Host(ctx)
	if err != nil {
		return dsn, terminate, err
	}

	port, err := spannerEmulator.MappedPort(ctx, "9010")
	if err != nil {
		return dsn, terminate, err
	}

	// OS environment needed for setting up instance and database
	emulatorHost := fmt.Sprintf("%s:%d", host, port.Int())
	os.Setenv("SPANNER_EMULATOR_HOST", emulatorHost)

	// Give the emulator a moment to fully initialize before setting up instances
	time.Sleep(2 * time.Second)

	if err := setupInstance(ctx); err != nil {
		return dsn, terminate, err
	}

	if err := createDatabase(ctx); err != nil {
		return dsn, terminate, err
	}

	return dsn, terminate, nil
}

func setupInstance(ctx context.Context) error {
	instanceAdmin, err := instance.NewInstanceAdminClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create instance admin client: %w", err)
	}
	defer instanceAdmin.Close()

	op, err := instanceAdmin.CreateInstance(ctx, &instancepb.CreateInstanceRequest{
		Parent:     fmt.Sprintf("projects/%s", testProject),
		InstanceId: testInstance,
		Instance: &instancepb.Instance{
			Config:      fmt.Sprintf("projects/%s/instanceConfigs/%s", testProject, "emulator-config"),
			DisplayName: testInstance,
			NodeCount:   1,
		},
	})
	if err != nil {
		return fmt.Errorf("could not create instance %s: %w", fmt.Sprintf("projects/%s/instances/%s", testProject, testInstance), err)
	}

	instance, err := op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("waiting for instance creation to finish failed: %w", err)
	}

	// The instance may not be ready to serve yet.
	if instance.State != instancepb.Instance_READY {
		fmt.Printf("instance state is not READY yet. Got state %v\n", instance.State)
	}
	fmt.Printf("Created emulator instance [%s]\n", testInstance)

	return nil
}

func createDatabase(ctx context.Context) error {
	adminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create database admin client: %w", err)
	}
	defer adminClient.Close()

	op, err := adminClient.CreateDatabase(ctx, &databasepb.CreateDatabaseRequest{
		Parent:          fmt.Sprintf("projects/%s/instances/%s", testProject, testInstance),
		CreateStatement: fmt.Sprintf("CREATE DATABASE `%s`", SpannerDB),
	})
	if err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}
	if _, err := op.Wait(ctx); err != nil {
		return fmt.Errorf("waiting for database creation failed: %w", err)
	}
	return nil
}

type spannerDialect int

const (
	// default value is 0 so if not present this will be selected
	SpannerDialectPGAdapter = iota
	SpannerDialectGoogleSQL
)

type DBOptions struct {
	// Will be used if TEST_WITH_SPANNER environment variable is set
	GCPProjectName string
	// Will be used if TEST_WITH_SPANNER environment variable is set
	GCPCredentialsFile string
	// Will be used if TEST_WITH_SPANNER environment variable is set
	GCPSpannerInstance string
	// Will be used if TEST_WITH_SPANNER environment variable is set
	GCPSpannerDatabase string

	// Will be ignored if TEST_WITH_SPANNER environment variable is set
	PostgresUser string
	// Will be ignored if TEST_WITH_SPANNER environment variable is set
	PostgresPassword string

	DBName string

	SpannerDialect spannerDialect
}

func getDSN(opts DBOptions) func(host string, port nat.Port) string {
	return func(host string, port nat.Port) string {
		return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			opts.PostgresUser,
			opts.PostgresPassword,
			host,
			port.Port(),
			opts.DBName,
		)
	}
}
