package repository

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"os"
	"time"

	pb "github.com/Andrewalifb/alpha-pos-system-sales-service/api/proto"
	"github.com/Andrewalifb/alpha-pos-system-sales-service/dto"
	"github.com/Andrewalifb/alpha-pos-system-sales-service/entity"

	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PosReturnRepository interface {
	CreatePosReturn(posReturn *entity.PosReturn) error
	ReadPosReturn(returnID string) (*pb.PosReturn, error)
	UpdatePosReturn(posReturn *entity.PosReturn) (*pb.PosReturn, error)
	DeletePosReturn(returnID string) error
	ReadAllPosReturns(pagination dto.Pagination, roleName string, jwtPayload *pb.JWTPayload) (*dto.PaginationResult, error)
}

type posReturnRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewPosReturnRepository(db *gorm.DB, redis *redis.Client) PosReturnRepository {
	return &posReturnRepository{
		db:    db,
		redis: redis,
	}
}

func (r *posReturnRepository) CreatePosReturn(posReturn *entity.PosReturn) error {
	result := r.db.Create(posReturn)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *posReturnRepository) ReadAllPosReturns(pagination dto.Pagination, roleName string, jwtPayload *pb.JWTPayload) (*dto.PaginationResult, error) {
	var posReturns []entity.PosReturn
	var totalRecords int64

	query := r.db.Model(&entity.PosReturn{})

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

	query.Find(&posReturns)
	query.Count(&totalRecords)

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pagination.Limit)))

	return &dto.PaginationResult{
		TotalRecords: totalRecords,
		Records:      posReturns,
		CurrentPage:  pagination.Page,
		TotalPages:   totalPages,
	}, nil
}

func (r *posReturnRepository) ReadPosReturn(returnID string) (*pb.PosReturn, error) {
	// Try to get the return from Redis first
	returnData, err := r.redis.Get(context.Background(), returnID).Result()
	if err == redis.Nil {
		// Return not found in Redis, get from PostgreSQL
		var posReturnEntity entity.PosReturn
		if err := r.db.Where("return_id = ?", returnID).First(&posReturnEntity).Error; err != nil {
			return nil, err
		}

		// Convert entity.PosReturn to pb.PosReturn
		posReturn := &pb.PosReturn{
			ReturnId:   posReturnEntity.ReturnID.String(),
			ReceiptId:  posReturnEntity.ReceiptID,
			ProductId:  posReturnEntity.ProductID.String(),
			Quantity:   int32(posReturnEntity.Quantity),
			Price:      float64(posReturnEntity.Price),
			Amount:     float64(posReturnEntity.Amount),
			ReturnDate: timestamppb.New(posReturnEntity.ReturnDate),
			Reason:     posReturnEntity.Reason,
			StoreId:    posReturnEntity.StoreID.String(),
			BranchId:   posReturnEntity.BranchID.String(),
			CompanyId:  posReturnEntity.CompanyID.String(),
			CreatedAt:  timestamppb.New(posReturnEntity.CreatedAt),
			CreatedBy:  posReturnEntity.CreatedBy.String(),
			UpdatedAt:  timestamppb.New(posReturnEntity.UpdatedAt),
			UpdatedBy:  posReturnEntity.UpdatedBy.String(),
		}

		// Store the return in Redis for future queries
		returnData, err := json.Marshal(posReturnEntity)
		if err != nil {
			return nil, err
		}
		err = r.redis.Set(context.Background(), returnID, returnData, 7*24*time.Hour).Err()
		if err != nil {
			return nil, err
		}

		return posReturn, nil
	} else if err != nil {
		return nil, err
	}

	// Return found in Redis, unmarshal the data
	var posReturnEntity entity.PosReturn
	err = json.Unmarshal([]byte(returnData), &posReturnEntity)
	if err != nil {
		return nil, err
	}

	// Convert entity.PosReturn to pb.PosReturn
	posReturn := &pb.PosReturn{
		ReturnId:   posReturnEntity.ReturnID.String(),
		ReceiptId:  posReturnEntity.ReceiptID,
		ProductId:  posReturnEntity.ProductID.String(),
		Quantity:   int32(posReturnEntity.Quantity),
		Price:      float64(posReturnEntity.Price),
		Amount:     float64(posReturnEntity.Amount),
		ReturnDate: timestamppb.New(posReturnEntity.ReturnDate),
		Reason:     posReturnEntity.Reason,
		StoreId:    posReturnEntity.StoreID.String(),
		BranchId:   posReturnEntity.BranchID.String(),
		CompanyId:  posReturnEntity.CompanyID.String(),
		CreatedAt:  timestamppb.New(posReturnEntity.CreatedAt),
		CreatedBy:  posReturnEntity.CreatedBy.String(),
		UpdatedAt:  timestamppb.New(posReturnEntity.UpdatedAt),
		UpdatedBy:  posReturnEntity.UpdatedBy.String(),
	}

	return posReturn, nil
}

func (r *posReturnRepository) UpdatePosReturn(posReturn *entity.PosReturn) (*pb.PosReturn, error) {
	if err := r.db.Save(posReturn).Error; err != nil {
		return nil, err
	}

	// Convert updated entity.PosReturn back to pb.PosReturn
	updatedPosReturn := &pb.PosReturn{
		ReturnId:   posReturn.ReturnID.String(),
		ReceiptId:  posReturn.ReceiptID,
		ProductId:  posReturn.ProductID.String(),
		Quantity:   int32(posReturn.Quantity),
		Price:      float64(posReturn.Price),
		Amount:     float64(posReturn.Amount),
		ReturnDate: timestamppb.New(posReturn.ReturnDate),
		Reason:     posReturn.Reason,
		StoreId:    posReturn.StoreID.String(),
		BranchId:   posReturn.BranchID.String(),
		CompanyId:  posReturn.CompanyID.String(),
		CreatedAt:  timestamppb.New(posReturn.CreatedAt),
		CreatedBy:  posReturn.CreatedBy.String(),
		UpdatedAt:  timestamppb.New(posReturn.UpdatedAt),
		UpdatedBy:  posReturn.UpdatedBy.String(),
	}

	// Update the return in Redis
	returnData, err := json.Marshal(posReturn)
	if err != nil {
		return nil, err
	}
	err = r.redis.Set(context.Background(), updatedPosReturn.ReturnId, returnData, 7*24*time.Hour).Err()
	if err != nil {
		return nil, err
	}

	return updatedPosReturn, nil
}

func (r *posReturnRepository) DeletePosReturn(returnID string) error {
	if err := r.db.Where("return_id = ?", returnID).Delete(&entity.PosReturn{}).Error; err != nil {
		return err
	}

	// Delete the return from Redis
	err := r.redis.Del(context.Background(), returnID).Err()
	if err != nil {
		return err
	}

	return nil
}
