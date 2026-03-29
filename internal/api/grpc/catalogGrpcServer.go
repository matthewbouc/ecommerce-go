package grpc

import (
	"context"
	"ecommerce/internal/dto"
	"ecommerce/internal/service"
	catalogpb "ecommerce/proto/catalog"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CatalogGrpcServer implements catalogpb.CatalogServiceServer.
// Embedding UnimplementedCatalogServiceServer means any RPC method we haven't
// implemented yet returns codes.Unimplemented instead of panicking at startup.
type CatalogGrpcServer struct {
	catalogpb.UnimplementedCatalogServiceServer
	svc service.CatalogService
}

func NewCatalogGrpcServer(svc service.CatalogService) *CatalogGrpcServer {
	return &CatalogGrpcServer{svc: svc}
}

// GetProduct is the primary RPC used by internal services (e.g. Order Service)
// to fetch authoritative product data — name, price, stock — at checkout.
func (s *CatalogGrpcServer) GetProduct(ctx context.Context, req *catalogpb.GetProductRequest) (*catalogpb.GetProductResponse, error) {
	productUuid, err := uuid.Parse(req.Uuid)
	if err != nil {
		// codes.InvalidArgument is the gRPC equivalent of HTTP 400.
		return nil, status.Errorf(codes.InvalidArgument, "invalid product uuid: %v", err)
	}

	product, err := s.svc.GetProduct(productUuid)
	if err != nil {
		// codes.NotFound → HTTP 404 equivalent.
		return nil, status.Errorf(codes.NotFound, "product not found: %v", err)
	}

	return &catalogpb.GetProductResponse{
		Product: toProtoProduct(product),
	}, nil
}

// GetProducts returns a paginated list, optionally filtered by seller.
// Used by the REST layer internally and by other services needing a product list.
func (s *CatalogGrpcServer) GetProducts(ctx context.Context, req *catalogpb.GetProductsRequest) (*catalogpb.GetProductsResponse, error) {
	filter := dto.GetProductsRequest{
		Page:     int(req.Page),
		PageSize: int(req.PageSize),
	}

	// seller_id is optional — only set the filter if it was provided.
	if req.SellerId != "" {
		sellerUuid, err := uuid.Parse(req.SellerId)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid seller_id: %v", err)
		}
		filter.SellerID = &sellerUuid
	}

	result, err := s.svc.GetProducts(filter)
	if err != nil {
		// codes.Internal → HTTP 500 equivalent. Use for unexpected server-side errors.
		return nil, status.Errorf(codes.Internal, "get products: %v", err)
	}

	protoProducts := make([]*catalogpb.Product, len(result.Products))
	for i, p := range result.Products {
		protoProducts[i] = toProtoProduct(p)
	}

	return &catalogpb.GetProductsResponse{
		Products: protoProducts,
		Total:    int32(result.Total),
	}, nil
}

// CreateProduct is called by the REST handler after it has authenticated
// the seller. The gRPC path exists so internal services can also create
// products programmatically without going through HTTP.
func (s *CatalogGrpcServer) CreateProduct(ctx context.Context, req *catalogpb.CreateProductRequest) (*catalogpb.CreateProductResponse, error) {
	sellerUuid, err := uuid.Parse(req.SellerId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid seller_id: %v", err)
	}

	categoryUuid, err := uuid.Parse(req.CategoryId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid category_id: %v", err)
	}

	product, err := s.svc.CreateProduct(sellerUuid, dto.CreateProductRequest{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       int(req.Stock),
		CategoryID:  categoryUuid,
		ImageUrl:    req.ImageUrl,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create product: %v", err)
	}

	return &catalogpb.CreateProductResponse{
		Product: toProtoProduct(product),
	}, nil
}

func (s *CatalogGrpcServer) UpdateProduct(ctx context.Context, req *catalogpb.UpdateProductRequest) (*catalogpb.UpdateProductResponse, error) {
	// TODO: extract seller identity from gRPC metadata once auth interceptor is added.
	productUuid, err := uuid.Parse(req.Uuid)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid product uuid: %v", err)
	}

	categoryUuid, err := uuid.Parse(req.CategoryId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid category_id: %v", err)
	}

	// Seller UUID is not in UpdateProductRequest — ownership is enforced via
	// the REST handler which has the authenticated user. This gRPC method is
	// internal-only for now; auth can be added via interceptor later.
	product, err := s.svc.UpdateProduct(uuid.Nil, productUuid, dto.UpdateProductRequest{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       int(req.Stock),
		CategoryID:  categoryUuid,
		ImageUrl:    req.ImageUrl,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "update product: %v", err)
	}

	return &catalogpb.UpdateProductResponse{
		Product: toProtoProduct(product),
	}, nil
}

func (s *CatalogGrpcServer) DeleteProduct(ctx context.Context, req *catalogpb.DeleteProductRequest) (*catalogpb.DeleteProductResponse, error) {
	sellerUuid, err := uuid.Parse(req.SellerId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid seller_id: %v", err)
	}

	productUuid, err := uuid.Parse(req.Uuid)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid product uuid: %v", err)
	}

	if err := s.svc.DeleteProduct(sellerUuid, productUuid); err != nil {
		return nil, status.Errorf(codes.Internal, "delete product: %v", err)
	}

	return &catalogpb.DeleteProductResponse{Success: true}, nil
}

func toProtoProduct(p dto.ProductResponse) *catalogpb.Product {
	return &catalogpb.Product{
		Uuid:        p.Uuid.String(),
		SellerId:    p.SellerID.String(),
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       int32(p.Stock),
		CategoryId:  p.Category.Uuid.String(),
		ImageUrl:    p.ImageUrl,
	}
}
