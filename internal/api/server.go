package api

import (
	"context"
	"ecommerce/config"
	grpcHandlers "ecommerce/internal/api/grpc"
	"ecommerce/internal/api/rest"
	"ecommerce/internal/api/rest/handlers"
	"ecommerce/internal/domain"
	"ecommerce/internal/helper"
	"ecommerce/internal/repository"
	"ecommerce/internal/service"
	"ecommerce/pkg/notification/sms/provider"
	catalogpb "ecommerce/proto/catalog"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func StartServer(cfg config.AppConfig) {
	db, err := gorm.Open(postgres.Open(cfg.DatabaseConfig), &gorm.Config{})
	if err != nil {
		log.Fatalf("database connection fatal error: %v\n", err)
	}
	log.Println("database connection established")

	// Category must be migrated before Product because Product has a FK to Category.Uuid.
	// GORM's AutoMigrate doesn't resolve FK ordering automatically
	err = db.AutoMigrate(
		&domain.User{},
		&domain.BankAccount{},
		&domain.CartItem{},
		&domain.Order{},
		&domain.OrderItem{},
		&domain.Category{},
		&domain.Product{},
	)
	if err != nil {
		log.Fatalf("database migration fatal error: %v\n", err)
	}

	auth := helper.SetupAuth(cfg.AuthSecret)
	smsClient := provider.NewTwilioSmsClient(cfg)

	// CatalogService is shared between the gRPC server and the REST handler
	catalogSvc := service.NewCatalogService(repository.NewCatalogRepository(db))

	// ── gRPC server ────────────────────────────────────────────────────────
	grpcServer := grpc.NewServer()
	catalogpb.RegisterCatalogServiceServer(grpcServer, grpcHandlers.NewCatalogGrpcServer(catalogSvc))

	go func() {
		lis, err := net.Listen("tcp", cfg.GrpcPort)
		if err != nil {
			log.Fatalf("gRPC listener failed on %s: %v\n", cfg.GrpcPort, err)
		}
		log.Printf("gRPC server listening on %s\n", cfg.GrpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("gRPC server failed: %v\n", err)
		}
	}()

	// ── REST / HTTP server ─────────────────────────────────────────────────
	app := fiber.New()

	restHandler := &rest.RestHandler{
		App:        app,
		DB:         db,
		Auth:       auth,
		Config:     cfg,
		SmsClient:  smsClient,
		CatalogSvc: catalogSvc,
	}

	setupRoutes(restHandler)

	go func() {
		if err := app.Listen(cfg.ServerPort); err != nil {
			log.Fatalf("REST server failed: %v\n", err)
		}
	}()

	// ── Graceful shutdown ──────────────────────────────────────────────────
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	sig := <-quit
	log.Printf("received signal %s — shutting down gracefully\n", sig)

	// 30s timeout for both gRPC and REST servers.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// gRPC: GracefulStop waits for all in-flight RPCs to complete before
	// closing connections; forces an immediate close instead of hanging forever
	grpcDone := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(grpcDone)
	}()

	select {
	case <-grpcDone:
		log.Println("gRPC server stopped")
	case <-ctx.Done():
		grpcServer.Stop()
		log.Println("gRPC server force-stopped (shutdown timeout exceeded)")
	}

	// REST: ShutdownWithContext stops accepting new requests and waits for active HTTP connections to finish
	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Printf("REST server forced shutdown: %v\n", err)
	} else {
		log.Println("REST server stopped")
	}

	log.Println("shutdown complete")
}

func setupRoutes(rh *rest.RestHandler) {
	handlers.SetupUserRoutes(rh)
	handlers.SetupCatalogRoutes(rh)
}
