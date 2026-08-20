package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"myapp/internal/config"
	"myapp/internal/domain"

	"github.com/shopspring/decimal"
)

type CreateProductInput struct {
	Name         string
	Description  string
	Price        decimal.Decimal
	CategorySlug string
	SellerID     int
	File         multipart.File
	Header       *multipart.FileHeader
}

type ProductRepository interface {
	ListBySellerID(ctx context.Context, sellerID, page int) ([]domain.Product, error)
	ListByCategorySlug(ctx context.Context, categorySlug string, page int) ([]domain.Product, error)
	List(ctx context.Context, page int) ([]domain.Product, error)
	Create(ctx context.Context, p domain.Product) (int, error)
	GetImageByID(ctx context.Context, productID, sellerID int) (string, error)
	Delete(ctx context.Context, productID, sellerID int) error
}

type FileManager interface {
	SaveFile(file multipart.File, header *multipart.FileHeader, saveDirectory string) (string, error)
	DeleteFile(path, filename string) error
}

type CategoryIDGetter interface {
	GetIDBySlug(ctx context.Context, categorySlug string) (int, error)
}

type ProductService struct {
	unitOfWork        UnitOfWork
	productRepository ProductRepository
	categoryIDGetter  CategoryIDGetter
	fileManager       FileManager
}

func NewProductService(unitOfWork UnitOfWork, productRepository ProductRepository, categoryIDGetter CategoryIDGetter, fileManager FileManager) ProductService {
	return ProductService{
		unitOfWork:        unitOfWork,
		productRepository: productRepository,
		categoryIDGetter:  categoryIDGetter,
		fileManager:       fileManager,
	}
}

func (s ProductService) Create(ctx context.Context, input CreateProductInput) (int, error) {
	categoryID, err := s.categoryIDGetter.GetIDBySlug(ctx, input.CategorySlug)
	if err != nil {
		return 0, err
	}

	filename, err := s.fileManager.SaveFile(input.File, input.Header, "uploads")
	if err != nil {
		return 0, err
	}

	product, err := domain.NewProduct(input.Name, input.Description, filename, input.Price, categoryID, input.SellerID)
	if err != nil {
		if err := s.fileManager.DeleteFile(config.UploadsPath, filename); err != nil {
			return 0, err
		}
		return 0, err
	}

	productID, err := s.productRepository.Create(ctx, product)
	if err != nil {
		return 0, err
	}

	return productID, nil
}

func (s ProductService) List(ctx context.Context, categorySlug string, page int) ([]domain.Product, error) {
	if page <= 0 {
		return nil, fmt.Errorf("get products by category slug: %w", domain.ErrValidation)
	}

	var products []domain.Product
	var err error
	if categorySlug == "" {
		products, err = s.productRepository.List(ctx, page)
	} else {
		products, err = s.productRepository.ListByCategorySlug(ctx, categorySlug, page)
	}

	if err != nil {
		return []domain.Product{}, err
	}

	return products, nil
}

func (s ProductService) ListBySellerID(ctx context.Context, sellerID, page int) ([]domain.Product, error) {
	if page <= 0 {
		return nil, fmt.Errorf("get products by sellerID: %w", domain.ErrValidation)
	}

	products, err := s.productRepository.ListBySellerID(ctx, sellerID, page)
	if err != nil {
		return []domain.Product{}, err
	}

	return products, nil
}

func (s ProductService) Delete(ctx context.Context, productID, sellerID int) error {
	filename, err := s.productRepository.GetImageByID(ctx, productID, sellerID)
	if err != nil {
		return err
	}

	if err := s.productRepository.Delete(ctx, productID, sellerID); err != nil {
		return err
	}

	if err := s.fileManager.DeleteFile(config.UploadsPath, filename); err != nil {
		return err
	}

	return nil
}
