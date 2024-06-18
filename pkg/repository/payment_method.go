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

type PosPaymentMethodRepository interface {
	CreatePosPaymentMethod(posPaymentMethod *entity.PosPaymentMethod) error
	ReadPosPaymentMethod(paymentMethodID string) (*pb.PosPaymentMethod, error)
	UpdatePosPaymentMethod(posPaymentMethod *entity.PosPaymentMethod) (*pb.PosPaymentMethod, error)
	DeletePosPaymentMethod(paymentMethodID string) error
	ReadAllPosPaymentMethods(agination dto.Pagination, roleName string, jwtPayload *pb.JWTPayload) (*dto.PaginationResult, error)
}

type posPaymentMethodRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewPosPaymentMethodRepository(db *gorm.DB, redis *redis.Client) PosPaymentMethodRepository {
	return &posPaymentMethodRepository{
		db:    db,
		redis: redis,
	}
}

func (r *posPaymentMethodRepository) CreatePosPaymentMethod(posPaymentMethod *entity.PosPaymentMethod) error {
	result := r.db.Create(posPaymentMethod)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *posPaymentMethodRepository) ReadPosPaymentMethod(paymentMethodID string) (*pb.PosPaymentMethod, error) {
	// Try to get the payment method from Redis first
	paymentMethodData, err := r.redis.Get(context.Background(), paymentMethodID).Result()
	if err == redis.Nil {
		// Payment method not found in Redis, get from PostgreSQL
		var posPaymentMethodEntity entity.PosPaymentMethod
		if err := r.db.Where("payment_method_id = ?", paymentMethodID).First(&posPaymentMethodEntity).Error; err != nil {
			return nil, err
		}

		// Convert entity.PosPaymentMethod to pb.PosPaymentMethod
		posPaymentMethod := &pb.PosPaymentMethod{
			PaymentMethodId: posPaymentMethodEntity.PaymentMethodID.String(),
			MethodName:      posPaymentMethodEntity.MethodName,
			CompanyId:       posPaymentMethodEntity.CompanyID.String(),
			CreatedAt:       timestamppb.New(posPaymentMethodEntity.CreatedAt),
			CreatedBy:       posPaymentMethodEntity.CreatedBy.String(),
			UpdatedAt:       timestamppb.New(posPaymentMethodEntity.UpdatedAt),
			UpdatedBy:       posPaymentMethodEntity.UpdatedBy.String(),
		}

		// Store the payment method in Redis for future queries
		paymentMethodData, err := json.Marshal(posPaymentMethodEntity)
		if err != nil {
			return nil, err
		}
		err = r.redis.Set(context.Background(), paymentMethodID, paymentMethodData, 7*24*time.Hour).Err()
		if err != nil {
			return nil, err
		}

		return posPaymentMethod, nil
	} else if err != nil {
		return nil, err
	}

	// Payment method found in Redis, unmarshal the data
	var posPaymentMethodEntity entity.PosPaymentMethod
	err = json.Unmarshal([]byte(paymentMethodData), &posPaymentMethodEntity)
	if err != nil {
		return nil, err
	}

	// Convert entity.PosPaymentMethod to pb.PosPaymentMethod
	posPaymentMethod := &pb.PosPaymentMethod{
		PaymentMethodId: posPaymentMethodEntity.PaymentMethodID.String(),
		MethodName:      posPaymentMethodEntity.MethodName,
		CompanyId:       posPaymentMethodEntity.CompanyID.String(),
		CreatedAt:       timestamppb.New(posPaymentMethodEntity.CreatedAt),
		CreatedBy:       posPaymentMethodEntity.CreatedBy.String(),
		UpdatedAt:       timestamppb.New(posPaymentMethodEntity.UpdatedAt),
		UpdatedBy:       posPaymentMethodEntity.UpdatedBy.String(),
	}

	return posPaymentMethod, nil
}

func (r *posPaymentMethodRepository) UpdatePosPaymentMethod(posPaymentMethod *entity.PosPaymentMethod) (*pb.PosPaymentMethod, error) {
	if err := r.db.Save(posPaymentMethod).Error; err != nil {
		return nil, err
	}

	// Convert updated entity.PosPaymentMethod back to pb.PosPaymentMethod
	updatedPosPaymentMethod := &pb.PosPaymentMethod{
		PaymentMethodId: posPaymentMethod.PaymentMethodID.String(),
		MethodName:      posPaymentMethod.MethodName,
		CompanyId:       posPaymentMethod.CompanyID.String(),
		CreatedAt:       timestamppb.New(posPaymentMethod.CreatedAt),
		CreatedBy:       posPaymentMethod.CreatedBy.String(),
		UpdatedAt:       timestamppb.New(posPaymentMethod.UpdatedAt),
		UpdatedBy:       posPaymentMethod.UpdatedBy.String(),
	}

	// Update the payment method in Redis
	paymentMethodData, err := json.Marshal(posPaymentMethod)
	if err != nil {
		return nil, err
	}
	err = r.redis.Set(context.Background(), updatedPosPaymentMethod.PaymentMethodId, paymentMethodData, 7*24*time.Hour).Err()
	if err != nil {
		return nil, err
	}

	return updatedPosPaymentMethod, nil
}

func (r *posPaymentMethodRepository) DeletePosPaymentMethod(paymentMethodID string) error {
	if err := r.db.Where("payment_method_id = ?", paymentMethodID).Delete(&entity.PosPaymentMethod{}).Error; err != nil {
		return err
	}

	// Delete the payment method from Redis
	err := r.redis.Del(context.Background(), paymentMethodID).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *posPaymentMethodRepository) ReadAllPosPaymentMethods(pagination dto.Pagination, roleName string, jwtPayload *pb.JWTPayload) (*dto.PaginationResult, error) {
	var posPaymentMethods []entity.PosPaymentMethod
	var totalRecords int64

	query := r.db.Model(&entity.PosPaymentMethod{})

	companyRole := os.Getenv("COMPANY_USER_ROLE")
	branchRole := os.Getenv("BRANCH_USER_ROLE")
	storeRole := os.Getenv("STORE_USER_ROLE")

	switch roleName {
	case companyRole:
		query = query.Where("company_id = ?", jwtPayload.CompanyId)
	case branchRole:
		query = query.Where("company_id = ?", jwtPayload.CompanyId)
	case storeRole:
		query = query.Where("company_id = ?", jwtPayload.CompanyId)
	default:
		return nil, errors.New("invalid role")
	}

	if pagination.Limit > 0 && pagination.Page > 0 {
		offset := (pagination.Page - 1) * pagination.Limit
		query = query.Offset(offset).Limit(pagination.Limit)
	}

	query.Find(&posPaymentMethods)
	query.Count(&totalRecords)

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pagination.Limit)))

	return &dto.PaginationResult{
		TotalRecords: totalRecords,
		Records:      posPaymentMethods,
		CurrentPage:  pagination.Page,
		TotalPages:   totalPages,
	}, nil
}
