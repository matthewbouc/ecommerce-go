package api

import (
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

	// ── gRPC server ────────────────────────────────────────────────────────
	// CatalogService is instantiated here so the same instance is shared
	// between the gRPC server and the REST handler (via SetupCatalogRoutes).
	catalogSvc := service.NewCatalogService(repository.NewCatalogRepository(db))
	go startGrpcServer(cfg.GrpcPort, catalogSvc)

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

	log.Fatal(app.Listen(cfg.ServerPort))
}

func setupRoutes(rh *rest.RestHandler) {
	handlers.SetupUserRoutes(rh)
	handlers.SetupCatalogRoutes(rh)
}

// startGrpcServer starts the gRPC server on the configured port.
// Called in a goroutine so it doesn't block the REST server startup.
func startGrpcServer(port string, catalogSvc service.CatalogService) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("gRPC listener failed on %s: %v\n", port, err)
	}

	grpcServer := grpc.NewServer()
	catalogpb.RegisterCatalogServiceServer(grpcServer, grpcHandlers.NewCatalogGrpcServer(catalogSvc))

	log.Printf("gRPC server listening on %s\n", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("gRPC server failed: %v\n", err)
	}
}
