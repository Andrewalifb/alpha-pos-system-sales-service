package repository

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"os"
	"strconv"
	"time"

	pb "github.com/Andrewalifb/alpha-pos-system-sales-service/api/proto"
	"github.com/Andrewalifb/alpha-pos-system-sales-service/dto"
	"github.com/Andrewalifb/alpha-pos-system-sales-service/entity"
	"github.com/Andrewalifb/alpha-pos-system-sales-service/utils"

	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PosSaleRepository interface {
	CreatePosSales(posSale []*entity.PosSale, token string) ([]*entity.PosSale, error)
	ReadPosSale(saleID string) (*pb.PosSale, error)
	UpdatePosSale(posSale *entity.PosSale) (*pb.PosSale, error)
	DeletePosSale(saleID string) error
	ReadAllPosSales(pagination dto.Pagination, roleName string, jwtPayload *pb.JWTPayload) (*dto.PaginationResult, error)
}

type posSaleRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewPosSaleRepository(db *gorm.DB, redis *redis.Client) PosSaleRepository {
	return &posSaleRepository{
		db:    db,
		redis: redis,
	}
}

func (r *posSaleRepository) CreatePosSales(posSales []*entity.PosSale, token string) ([]*entity.PosSale, error) {
	var createdPosSales []*entity.PosSale
	err := r.db.Transaction(func(tx *gorm.DB) error {
		// Get a new ReceiptID from the Store service
		storeID := posSales[0].StoreID.String()
		receiptID, err := utils.GetNextReceiptID(storeID, token)
		if err != nil {
			return err
		}
		for _, posSale := range posSales {

			posSale.ReceiptID = strconv.Itoa(receiptID.Data.ReceiptID)
			result := tx.Create(posSale)
			if result.Error != nil {
				return result.Error
			}
			createdPosSales = append(createdPosSales, posSale)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return createdPosSales, nil
}

func (r *posSaleRepository) ReadAllPosSales(pagination dto.Pagination, roleName string, jwtPayload *pb.JWTPayload) (*dto.PaginationResult, error) {
	var posSales []entity.PosSale
	var totalRecords int64

	query := r.db.Model(&entity.PosSale{})

	companyRole := os.Getenv("COMPANY_USER_ROLE")
	branchRole := os.Getenv("BRANCH_USER_ROLE")
	storeRole := os.Getenv("STORE_USER_ROLE")

	switch roleName {
	case companyRole:
		query = query.Where("company_id = ?", jwtPayload.CompanyId)
	case branchRole:
		query = query.Where("branch_id = ?", jwtPayload.BranchId)
	case storeRole:
		query = query.Where("store_id = ?", jwtPayload.StoreId)
	default:
		return nil, errors.New("invalid role")
	}

	if pagination.Limit > 0 && pagination.Page > 0 {
		offset := (pagination.Page - 1) * pagination.Limit
		query = query.Offset(offset).Limit(pagination.Limit)
	}

	query.Find(&posSales)
	query.Count(&totalRecords)

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pagination.Limit)))

	return &dto.PaginationResult{
		TotalRecords: totalRecords,
		Records:      posSales,
		CurrentPage:  pagination.Page,
		TotalPages:   totalPages,
	}, nil
}

func (r *posSaleRepository) ReadPosSale(saleID string) (*pb.PosSale, error) {
	// Try to get the sale from Redis first
	saleData, err := r.redis.Get(context.Background(), saleID).Result()
	if err == redis.Nil {
		// Sale not found in Redis, get from PostgreSQL
		var posSaleEntity entity.PosSale
		if err := r.db.Where("sale_id = ?", saleID).First(&posSaleEntity).Error; err != nil {
			return nil, err
		}

		// Convert entity.PosSale to pb.PosSale
		posSale := &pb.PosSale{
			SaleId:          posSaleEntity.SaleID.String(),
			ReceiptId:       posSaleEntity.ReceiptID,
			ProductId:       posSaleEntity.ProductID.String(),
			CustomerId:      posSaleEntity.CustomerID.String(),
			Quantity:        int32(posSaleEntity.Quantity),
			Price:           posSaleEntity.Price,
			SaleDate:        timestamppb.New(posSaleEntity.SaleDate),
			TotalPrice:      posSaleEntity.TotalPrice,
			StoreId:         posSaleEntity.StoreID.String(),
			CashierId:       posSaleEntity.CashierID.String(),
			PaymentMethodId: posSaleEntity.PaymentMethodID.String(),
			BranchId:        posSaleEntity.BranchID.String(),
			CompanyId:       posSaleEntity.CompanyID.String(),
			CreatedAt:       timestamppb.New(posSaleEntity.CreatedAt),
			CreatedBy:       posSaleEntity.CreatedBy.String(),
			UpdatedAt:       timestamppb.New(posSaleEntity.UpdatedAt),
			UpdatedBy:       posSaleEntity.UpdatedBy.String(),
		}

		// Store the sale in Redis for future queries
		saleData, err := json.Marshal(posSaleEntity)
		if err != nil {
			return nil, err
		}
		err = r.redis.Set(context.Background(), saleID, saleData, 7*24*time.Hour).Err()
		if err != nil {
			return nil, err
		}

		return posSale, nil
	} else if err != nil {
		return nil, err
	}

	// Sale found in Redis, unmarshal the data
	var posSaleEntity entity.PosSale
	err = json.Unmarshal([]byte(saleData), &posSaleEntity)
	if err != nil {
		return nil, err
	}

	// Convert entity.PosSale to pb.PosSale
	posSale := &pb.PosSale{
		SaleId:          posSaleEntity.SaleID.String(),
		ReceiptId:       posSaleEntity.ReceiptID,
		ProductId:       posSaleEntity.ProductID.String(),
		CustomerId:      posSaleEntity.CustomerID.String(),
		Quantity:        int32(posSaleEntity.Quantity),
		Price:           posSaleEntity.Price,
		SaleDate:        timestamppb.New(posSaleEntity.SaleDate),
		TotalPrice:      posSaleEntity.TotalPrice,
		StoreId:         posSaleEntity.StoreID.String(),
		CashierId:       posSaleEntity.CashierID.String(),
		PaymentMethodId: posSaleEntity.PaymentMethodID.String(),
		BranchId:        posSaleEntity.BranchID.String(),
		CompanyId:       posSaleEntity.CompanyID.String(),
		CreatedAt:       timestamppb.New(posSaleEntity.CreatedAt),
		CreatedBy:       posSaleEntity.CreatedBy.String(),
		UpdatedAt:       timestamppb.New(posSaleEntity.UpdatedAt),
		UpdatedBy:       posSaleEntity.UpdatedBy.String(),
	}

	return posSale, nil
}

func (r *posSaleRepository) UpdatePosSale(posSale *entity.PosSale) (*pb.PosSale, error) {
	if err := r.db.Save(posSale).Error; err != nil {
		return nil, err
	}

	// Convert updated entity.PosSale back to pb.PosSale
	updatedPosSale := &pb.PosSale{
		SaleId:          posSale.SaleID.String(),
		ReceiptId:       posSale.ReceiptID,
		ProductId:       posSale.ProductID.String(),
		CustomerId:      posSale.CustomerID.String(),
		Quantity:        int32(posSale.Quantity),
		Price:           posSale.Price,
		SaleDate:        timestamppb.New(posSale.SaleDate),
		TotalPrice:      posSale.TotalPrice,
		StoreId:         posSale.StoreID.String(),
		CashierId:       posSale.CashierID.String(),
		PaymentMethodId: posSale.PaymentMethodID.String(),
		BranchId:        posSale.BranchID.String(),
		CompanyId:       posSale.CompanyID.String(),
		CreatedAt:       timestamppb.New(posSale.CreatedAt),
		CreatedBy:       posSale.CreatedBy.String(),
		UpdatedAt:       timestamppb.New(posSale.UpdatedAt),
		UpdatedBy:       posSale.UpdatedBy.String(),
	}

	// Update the sale in Redis
	saleData, err := json.Marshal(posSale)
	if err != nil {
		return nil, err
	}
	err = r.redis.Set(context.Background(), updatedPosSale.SaleId, saleData, 7*24*time.Hour).Err()
	if err != nil {
		return nil, err
	}

	return updatedPosSale, nil
}

func (r *posSaleRepository) DeletePosSale(saleID string) error {
	if err := r.db.Where("sale_id = ?", saleID).Delete(&entity.PosSale{}).Error; err != nil {
		return err
	}

	// Delete the sale from Redis
	err := r.redis.Del(context.Background(), saleID).Err()
	if err != nil {
		return err
	}

	return nil
}
