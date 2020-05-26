package main

import (
	"context"
	"flag"
	"fmt"
	"mime"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/uw-labs/go-mono/cmd/user-api/internal/repo"
	"github.com/uw-labs/go-mono/cmd/user-api/internal/server"
	"github.com/uw-labs/go-mono/cmd/user-api/third_party/swagger"
	pkgctx "github.com/uw-labs/go-mono/pkg/context"
	usersservicepb "github.com/uw-labs/go-mono/proto/gen/go/uwlabs/users/service/v1"
)

//go:generate cp ../../proto/gen/openapiv2/uwlabs/users/service/v1/service.swagger.json ./third_party/swagger/swagger.json
//go:generate go-bindata -pkg swagger -prefix third_party/swagger -nometadata -ignore bindata -o ./third_party/swagger/bindata.go ./third_party/swagger

var (
	postgresURL   = flag.String("postgres-url", "", "The URL of the postgres database to connect to.")
	grpcPort      = flag.Uint("grpc-port", 8080, "The port to serve the gRPC server on.")
	gatewayPort   = flag.Uint("grpc-gateway-port", 8081, "The port to serve the gRPC-Gateway on.")
	adminUser     = flag.String("admin-user", "admin", "The username of the admin user.")
	adminPassword = flag.String("admin-password", "", "The password of the admin user.")
)

func main() {
	flag.Parse()

	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.StampMilli,
	}

	if *postgresURL == "" {
		logger.Fatal("postgres-url must be specified")
	}

	if *adminPassword == "" {
		logger.Fatal("admin-password must be specified")
	}

	err := run(logger, *postgresURL, *adminUser, *adminPassword, *grpcPort, *gatewayPort)
	if err != nil {
		logger.WithError(err).Fatal()
	}
}

func run(logger *logrus.Logger, postgresURL, adminUser, adminPassword string, grpcPort, gatewayPort uint) (err error) {
	ctx := pkgctx.WithSignalHandler(context.Background())

	rp, err := repo.NewRepository(postgresURL, logger)
	if err != nil {
		return fmt.Errorf("create repository: %w", err)
	}
	defer func() {
		cErr := rp.Close()
		if err == nil {
			err = cErr
		}
	}()

	admin := &server.User{
		Username: adminUser,
		Password: adminPassword,
	}

	grpcAddr := ":" + strconv.Itoa(int(grpcPort))
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return fmt.Errorf("starting TCP listener: %w", err)
	}

	srv := grpc.NewServer()

	backend := &server.Server{
		Logger: logger,
		Repo:   rp,
		Admin:  admin,
	}

	usersservicepb.RegisterUserReaderServiceServer(srv, backend)
	usersservicepb.RegisterUserWriterServiceServer(srv, backend)

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return srv.Serve(lis)
	})
	eg.Go(func() error {
		<-ctx.Done()
		tCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		ok := make(chan struct{})
		go func() {
			srv.GracefulStop()
			close(ok)
		}()

		select {
		case <-tCtx.Done():
			srv.Stop()
		case <-ok:
		}
		return nil
	})

	cc, err := grpc.Dial(grpcAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return fmt.Errorf("dialling gRPC server: %w", err)
	}
	defer func() {
		cErr := cc.Close()
		if err == nil {
			err = cErr
		}
	}()

	jsonpb := &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			Indent: "  ",
		},
	}
	gwmux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, jsonpb),
	)
	err = usersservicepb.RegisterUserReaderServiceHandler(ctx, gwmux, cc)
	if err != nil {
		return fmt.Errorf("register gateway: %w", err)
	}
	err = usersservicepb.RegisterUserWriterServiceHandler(ctx, gwmux, cc)
	if err != nil {
		return fmt.Errorf("register gateway: %w", err)
	}

	err = mime.AddExtensionType(".svg", "image/svg+xml")
	if err != nil {
		return fmt.Errorf("register svg mime type: %w", err)
	}

	swaggerHandler := http.FileServer(&assetfs.AssetFS{
		Asset:     swagger.Asset,
		AssetDir:  swagger.AssetDir,
		AssetInfo: swagger.AssetInfo,
	})

	gatewayAddr := ":" + strconv.Itoa(int(gatewayPort))
	gwServer := &http.Server{
		Addr: gatewayAddr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/v1") {
				gwmux.ServeHTTP(w, r)
				return
			}

			swaggerHandler.ServeHTTP(w, r)
		}),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	eg.Go(func() error {
		<-ctx.Done()
		tCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		err := gwServer.Shutdown(tCtx)
		if err != nil {
			return fmt.Errorf("shutdown gRPC gateway: %w", err)
		}
		return nil
	})
	eg.Go(func() error {
		logger.Info("Serving gRPC-Gateway and Swagger Documentation on http://localhost:", gatewayPort)
		err = gwServer.ListenAndServe()
		if err != http.ErrServerClosed {
			return fmt.Errorf("serve gRPC gateway: %w", err)
		}
		return nil
	})

	return eg.Wait()
}
