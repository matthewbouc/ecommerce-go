package rest

import (
	"ecommerce/config"
	"ecommerce/internal/helper"
	"ecommerce/internal/service"
	"ecommerce/pkg/notification/sms"
	catalogpb "ecommerce/proto/catalog"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type RestHandler struct {
	App           *fiber.App
	DB            *gorm.DB
	Auth          helper.Auth
	Config        config.AppConfig
	SmsClient     sms.SmsClient
	CatalogSvc    service.CatalogService
	CatalogClient catalogpb.CatalogServiceClient
}
