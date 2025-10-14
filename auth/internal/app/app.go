package app

import (
	"context"
	"flag"
	"log"
	"net"
	"sync"

	"github.com/Dokhoyan/common/pkg/closer"
	"github.com/Dokhoyan/go-messenger-microservices/auth/internal/config"
	"github.com/Dokhoyan/go-messenger-microservices/auth/pkg/api/access_v1"
	"github.com/Dokhoyan/go-messenger-microservices/auth/pkg/api/auth_v1"
	"google.golang.org/grpc"

	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

const APISwaggerPath = "/api.swagger.json"
const (
	loggerMaxSize    = 10
	loggerMaxBackups = 3
	loggerMaxAge     = 3
	reqLimite        = 100
	reqSecondtime    = 1
)

var configPath string
var logLevel = flag.String("level", "info", "log level for logger")

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type App struct {
	serviceProvider  *serviceProvider
	grpcServer       *grpc.Server
}

func NewApp(ctx context.Context)(*App, error){
	a:=&App{}

	err:=a.initDeps(ctx)
	if err!=nil{
		return nil, err
	}

	return a, nil
}

func (a *App) Run()(error){
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		err := a.runGRPCServer()
		if err != nil {
			log.Fatalf("failed to run GRPC server: %v", err)
		}
	}()

	wg.Wait()

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	flag.Parse()

	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initGRPCServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	err := config.Load(configPath)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {


    a.grpcServer = grpc.NewServer(
        grpc.Creds(insecure.NewCredentials()),
        grpc.ChainUnaryInterceptor(),
    )

	reflection.Register(a.grpcServer)

	auth_v1.RegisterAuthV1Server(a.grpcServer, a.serviceProvider.AuthImpl(ctx))
	access_v1.RegisterAccessV1Server(a.grpcServer, a.serviceProvider.AccessImpl(ctx))


	return nil
}

func (a *App) runGRPCServer() error {
	log.Printf("GRPC server is running on %s", a.serviceProvider.GRPCConfig().Address())

	list, err := net.Listen("tcp", a.serviceProvider.GRPCConfig().Address())
	if err != nil {
		return err
	}

	err = a.grpcServer.Serve(list)
	if err != nil {
		return err
	}

	return nil
}