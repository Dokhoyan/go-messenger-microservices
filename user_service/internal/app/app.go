package app

import (
	"context"
	"flag"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/Dokhoyan/common/pkg/logger"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/closer"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/config"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/interceptor"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/metric"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/tracing"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/pkg/api/access_v1"
	"github.com/Dokhoyan/go-messenger-microservices/user_service/pkg/api/auth_v1"
	desc "github.com/Dokhoyan/go-messenger-microservices/user_service/pkg/api/user_v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rakyll/statik/fs"
	"github.com/rs/cors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"gopkg.in/natefinch/lumberjack.v2"

	_ "github.com/Dokhoyan/go-messenger-microservices/user_service/statik" // инициализация шаблона swagger
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

const APISwaggerPath = "/api.swagger.json"
const (
	loggerMaxSize    = 10
	loggerMaxBackups = 3
	loggerMaxAge     = 3
)

var logLevel = flag.String("level", "info", "log level for logger")

type App struct {
	serviceProvider  *serviceProvider
	grpcServer       *grpc.Server
	httpServer       *http.Server
	swaggerServer    *http.Server
	prometheusServer *http.Server
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
	wg.Add(4)

	go func() {
		defer wg.Done()

		err := a.runGRPCServer()
		if err != nil {
			log.Fatalf("failed to run GRPC server: %v", err)
		}
	}()

	go func() {
		defer wg.Done()

		err := a.runHTTPServer()
		if err != nil {
			log.Fatalf("failed to run HTTP server: %v", err)
		}
	}()

	go func() {
		defer wg.Done()

		err := a.runSwaggerServer()
		if err != nil {
			log.Fatalf("failed to run Swagger server: %v", err)
		}
	}()

	go func() {
		defer wg.Done()

		err := a.runPrometheus()
		if err != nil {
			log.Fatalf("failed to run prometheus server: %v", err)
		}
	}()

	wg.Wait()

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initGRPCServer,
		a.initHTTPServer,
		a.initSwaggerServer,
		a.initLogger,
		a.initPrometheus,
		a.initMetrics,
		a.initTracing,
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
	err := config.Load(".env")
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
	creds, err := credentials.NewServerTLSFromFile("service.pem", "service.key")
    if err != nil {
        return err
    }

    a.grpcServer = grpc.NewServer(
        grpc.Creds(creds),
        grpc.ChainUnaryInterceptor(
			interceptor.ValidateInterceptor,
			interceptor.LoggingInterceptor,
			interceptor.ServerTracingInterceptor,
			interceptor.MetricsInterceptor,
		),
    )

	reflection.Register(a.grpcServer)

	desc.RegisterUserV1Server(a.grpcServer, a.serviceProvider.UserImpl(ctx))
	auth_v1.RegisterAuthV1Server(a.grpcServer, a.serviceProvider.AuthImpl(ctx))
	access_v1.RegisterAccessV1Server(a.grpcServer, a.serviceProvider.AccessImpl(ctx))


	return nil
}

func (a *App) initPrometheus(_ context.Context) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	a.prometheusServer = &http.Server{
		ReadHeaderTimeout: 0,
		Addr:              a.serviceProvider.PrometheusConfig().Address(),
		Handler:           mux,
	}

	return nil
}

func (a *App) initMetrics(ctx context.Context) error {
	return metric.Init(ctx)
}

func (a *App) initHTTPServer(ctx context.Context) error {
	mux := runtime.NewServeMux()
	creds, err := credentials.NewClientTLSFromFile("service.pem", "")
    if err != nil {
        return err
    }

    opts := []grpc.DialOption{
        grpc.WithTransportCredentials(creds),
    }

	err = desc.RegisterUserV1HandlerFromEndpoint(ctx, mux, a.serviceProvider.GRPCConfig().Address(), opts)
	if err != nil {
		return err
	}

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Authorization"},
		AllowCredentials: true,
	})

	a.httpServer = &http.Server{
		Addr:    a.serviceProvider.HTTPConfig().Address(),
		Handler: corsMiddleware.Handler(mux),
	}

	return nil
}

func (a *App) initSwaggerServer(_ context.Context) error {
	statikFs, err := fs.New()
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.StripPrefix("/", http.FileServer(statikFs)))
	mux.HandleFunc(APISwaggerPath, serveSwaggerFile(APISwaggerPath))

	a.swaggerServer = &http.Server{
		Addr:    a.serviceProvider.SwaggerConfig().Address(),
		Handler: mux,
	}

	return nil
}

func (a *App) initLogger(_ context.Context) error {
	logger.Init(a.getCore(a.getLevel()))

	return nil
}

func (a *App) initTracing(_ context.Context) error {
	tracing.Init("auth-service")

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

func (a *App) runHTTPServer() error {
	log.Printf("HTTP server is running on %s", a.serviceProvider.HTTPConfig().Address())

	err := a.httpServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (a *App) runPrometheus() error {
	log.Printf("Prometheus server listening at %s", a.serviceProvider.PrometheusConfig().Address())

	err := a.prometheusServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (a *App) runSwaggerServer() error {
	log.Printf("Swagger server is running on %s", a.serviceProvider.SwaggerConfig().Address())

	err := a.swaggerServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func serveSwaggerFile(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Serving swagger file: %s", path)

		statikFs, err := fs.New()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Open swagger file: %s", path)

		file, err := statikFs.Open(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		log.Printf("Read swagger file: %s", path)

		content, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Write swagger file: %s", path)

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Served swagger file: %s", path)
	}
}

func (a *App) getCore(level zap.AtomicLevel) zapcore.Core {
	stdout := zapcore.AddSync(os.Stdout)

	file := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    loggerMaxSize, //mb
		MaxBackups: loggerMaxBackups,
		MaxAge:     loggerMaxAge,
	})

	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
	fileEncoder := zapcore.NewJSONEncoder(productionCfg)

	return zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, level),
		zapcore.NewCore(fileEncoder, file, level),
	)
}

func (a *App) getLevel() zap.AtomicLevel {
	flag.Parse()

	var level zapcore.Level

	if err := level.Set(*logLevel); err != nil {
		log.Fatalf("failed to set log level")
	}

	return zap.NewAtomicLevelAt(level)
}