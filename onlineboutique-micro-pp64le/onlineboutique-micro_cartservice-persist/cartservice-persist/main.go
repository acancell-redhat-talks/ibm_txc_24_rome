package main

import (
	"context"
	"log"
	"strings"
	"time"

	grpcc "github.com/go-micro/plugins/v4/client/grpc"
	_ "github.com/go-micro/plugins/v4/registry/kubernetes"
	grpcs "github.com/go-micro/plugins/v4/server/grpc"
	"github.com/go-micro/plugins/v4/wrapper/trace/opentelemetry"
	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
	"go-micro.dev/v4/server"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"

	"github.com/go-micro/demo/cartservice/cartstore"
	"github.com/go-micro/demo/cartservice/config"
	"github.com/go-micro/demo/cartservice/handler"
	pb "github.com/go-micro/demo/cartservice/proto"

	"github.com/valkey-io/valkey-go"
	"strconv"
)

var (
	name    = "cartservice"
	version = "1.0.0"
)

func main() {
	// Load conigurations
	if err := config.Load(); err != nil {
		logger.Fatal(err)
	}

	// Create service
	srv := micro.NewService(
		micro.Server(grpcs.NewServer()),
		micro.Client(grpcc.NewClient()),
	)
	opts := []micro.Option{
		micro.Name(name),
		micro.Version(version),
		micro.Address(config.Address()),
	}
	if cfg := config.Tracing(); cfg.Enable {
		tp, err := newTracerProvider(name, srv.Server().Options().Id, cfg.Jaeger.URL)
		if err != nil {
			logger.Fatal(err)
		}
		defer func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()
			if err := tp.Shutdown(ctx); err != nil {
				logger.Fatal(err)
			}
		}()
		otel.SetTracerProvider(tp)
		otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
		traceOpts := []opentelemetry.Option{
			opentelemetry.WithHandleFilter(func(ctx context.Context, r server.Request) bool {
				if e := r.Endpoint(); strings.HasPrefix(e, "Health.") {
					return true
				}
				return false
			}),
		}
		opts = append(opts, micro.WrapHandler(opentelemetry.NewHandlerWrapper(traceOpts...)))
	}
	srv.Init(opts...)

	// Connect to db
	client_db, err_db := valkey.NewClient(valkey.ClientOption{InitAddress: []string{config.Redis().Addr}})
	if err_db != nil {
		logger.Fatal(err_db)
		panic(err_db)
	}
	defer client_db.Close()

	// Load from db to map
	ctx := context.Background()
	initmap := make(map[string]map[string]int32)
	scan, err := client_db.Do(ctx, client_db.B().Scan().Cursor(0).Count(9999).Build()).AsScanEntry()
	if valkey.IsParseErr(err) {
		logger.Fatal(err.Error())
		panic(err.Error())
	}
	initkeys := scan.Elements
	for _, key := range initkeys {
		initmap[key] = map[string]int32{}
	
		m, err := client_db.Do(ctx, client_db.B().Hgetall().Key(key).Build()).AsStrMap()
		if valkey.IsParseErr(err) {
			logger.Fatal(err.Error())
			panic(err.Error())
		}
		for k, v := range m {
			initvalue, _ := strconv.Atoi(v)
			initmap[key][k] = int32(initvalue)
		}
	}

	// Register handler
	if err := pb.RegisterCartServiceHandler(srv.Server(), &handler.CartService{Store: cartstore.NewMemoryCartStore(client_db, initmap)}); err != nil {
		log.Fatal(err)
	}
	if err := pb.RegisterHealthHandler(srv.Server(), new(handler.Health)); err != nil {
		log.Fatal(err)
	}

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
