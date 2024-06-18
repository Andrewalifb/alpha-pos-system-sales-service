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

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
)

type PosOnlinePaymentRepository interface {
	CreatePosOnlinePayment(posOnlinePayment *entity.PosOnlinePayment) error
	ReadPosOnlinePayment(paymentID string) (*pb.PosOnlinePayment, error)
	UpdatePosOnlinePayment(posOnlinePayment *entity.PosOnlinePayment) (*pb.PosOnlinePayment, error)
	DeletePosOnlinePayment(paymentID string) error
	ReadAllPosOnlinePayments(pagination dto.Pagination, roleName string, jwtPayload *pb.JWTPayload) (*dto.PaginationResult, error)
}

type posOnlinePaymentRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewPosOnlinePaymentRepository(db *gorm.DB, redis *redis.Client) PosOnlinePaymentRepository {
	return &posOnlinePaymentRepository{
		db:    db,
		redis: redis,
	}
}

func (r *posOnlinePaymentRepository) CreatePosOnlinePayment(posOnlinePayment *entity.PosOnlinePayment) error {
	result := r.db.Create(posOnlinePayment)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *posOnlinePaymentRepository) ReadPosOnlinePayment(paymentID string) (*pb.PosOnlinePayment, error) {
	// Try to get the online payment from Redis first
	onlinePaymentData, err := r.redis.Get(context.Background(), paymentID).Result()
	if err == redis.Nil {
		// Online payment not found in Redis, get from PostgreSQL
		var posOnlinePaymentEntity entity.PosOnlinePayment
		if err := r.db.Where("payment_id = ?", paymentID).First(&posOnlinePaymentEntity).Error; err != nil {
			return nil, err
		}

		// Convert entity.PosOnlinePayment to pb.PosOnlinePayment
		posOnlinePayment := &pb.PosOnlinePayment{
			PaymentId:     posOnlinePaymentEntity.PaymentID.String(),
			StoreId:       posOnlinePaymentEntity.StoreID.String(),
			EmployeeId:    posOnlinePaymentEntity.EmployeeID.String(),
			PaymentDate:   timestamppb.New(posOnlinePaymentEntity.PaymentDate),
			ReceiptId:     posOnlinePaymentEntity.ReceiptID,
			Amount:        posOnlinePaymentEntity.Amount,
			PaymentMethod: posOnlinePaymentEntity.PaymentMethod.String(),
			RoleId:        posOnlinePaymentEntity.RoleID.String(),
			BranchId:      posOnlinePaymentEntity.BranchID.String(),
			CompanyId:     posOnlinePaymentEntity.CompanyID.String(),
			CreatedAt:     timestamppb.New(posOnlinePaymentEntity.CreatedAt),
			CreatedBy:     posOnlinePaymentEntity.CreatedBy.String(),
			UpdatedAt:     timestamppb.New(posOnlinePaymentEntity.UpdatedAt),
			UpdatedBy:     posOnlinePaymentEntity.UpdatedBy.String(),
		}

		// Store the online payment in Redis for future queries
		onlinePaymentData, err := json.Marshal(posOnlinePaymentEntity)
		if err != nil {
			return nil, err
		}
		err = r.redis.Set(context.Background(), paymentID, onlinePaymentData, 7*24*time.Hour).Err()
		if err != nil {
			return nil, err
		}

		return posOnlinePayment, nil
	} else if err != nil {
		return nil, err
	}

	// Online payment found in Redis, unmarshal the data
	var posOnlinePaymentEntity entity.PosOnlinePayment
	err = json.Unmarshal([]byte(onlinePaymentData), &posOnlinePaymentEntity)
	if err != nil {
		return nil, err
	}

	// Convert entity.PosOnlinePayment to pb.PosOnlinePayment
	posOnlinePayment := &pb.PosOnlinePayment{
		PaymentId:     posOnlinePaymentEntity.PaymentID.String(),
		StoreId:       posOnlinePaymentEntity.StoreID.String(),
		EmployeeId:    posOnlinePaymentEntity.EmployeeID.String(),
		PaymentDate:   timestamppb.New(posOnlinePaymentEntity.PaymentDate),
		ReceiptId:     posOnlinePaymentEntity.ReceiptID, // Added missing field
		Amount:        posOnlinePaymentEntity.Amount,
		PaymentMethod: posOnlinePaymentEntity.PaymentMethod.String(),
		RoleId:        posOnlinePaymentEntity.RoleID.String(),
		BranchId:      posOnlinePaymentEntity.BranchID.String(),
		CompanyId:     posOnlinePaymentEntity.CompanyID.String(),
		CreatedAt:     timestamppb.New(posOnlinePaymentEntity.CreatedAt),
		CreatedBy:     posOnlinePaymentEntity.CreatedBy.String(),
		UpdatedAt:     timestamppb.New(posOnlinePaymentEntity.UpdatedAt),
		UpdatedBy:     posOnlinePaymentEntity.UpdatedBy.String(),
	}

	return posOnlinePayment, nil
}

func (r *posOnlinePaymentRepository) UpdatePosOnlinePayment(posOnlinePayment *entity.PosOnlinePayment) (*pb.PosOnlinePayment, error) {
	if err := r.db.Save(posOnlinePayment).Error; err != nil {
		return nil, err
	}

	// Convert updated entity.PosOnlinePayment back to pb.PosOnlinePayment
	updatedPosOnlinePayment := &pb.PosOnlinePayment{
		PaymentId:     posOnlinePayment.PaymentID.String(),
		StoreId:       posOnlinePayment.StoreID.String(),
		EmployeeId:    posOnlinePayment.EmployeeID.String(),
		PaymentDate:   timestamppb.New(posOnlinePayment.PaymentDate),
		ReceiptId:     posOnlinePayment.ReceiptID,
		Amount:        posOnlinePayment.Amount,
		PaymentMethod: posOnlinePayment.PaymentMethod.String(),
		RoleId:        posOnlinePayment.RoleID.String(),
		BranchId:      posOnlinePayment.BranchID.String(),
		CompanyId:     posOnlinePayment.CompanyID.String(),
		CreatedAt:     timestamppb.New(posOnlinePayment.CreatedAt),
		CreatedBy:     posOnlinePayment.CreatedBy.String(),
		UpdatedAt:     timestamppb.New(posOnlinePayment.UpdatedAt),
		UpdatedBy:     posOnlinePayment.UpdatedBy.String(),
	}

	// Update the online payment in Redis
	onlinePaymentData, err := json.Marshal(updatedPosOnlinePayment)
	if err != nil {
		return nil, err
	}
	err = r.redis.Set(context.Background(), updatedPosOnlinePayment.PaymentId, onlinePaymentData, 7*24*time.Hour).Err()
	if err != nil {
		return nil, err
	}

	return updatedPosOnlinePayment, nil
}

func (r *posOnlinePaymentRepository) DeletePosOnlinePayment(paymentID string) error {
	if err := r.db.Where("payment_id = ?", paymentID).Delete(&entity.PosOnlinePayment{}).Error; err != nil {
		return err
	}

	// Delete the online payment from Redis
	err := r.redis.Del(context.Background(), paymentID).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *posOnlinePaymentRepository) ReadAllPosOnlinePayments(pagination dto.Pagination, roleName string, jwtPayload *pb.JWTPayload) (*dto.PaginationResult, error) {
	var posOnlinePayments []entity.PosOnlinePayment
	var totalRecords int64

	query := r.db.Model(&entity.PosOnlinePayment{})

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

	query.Find(&posOnlinePayments)
	query.Count(&totalRecords)

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pagination.Limit)))

	return &dto.PaginationResult{
		TotalRecords: totalRecords,
		Records:      posOnlinePayments,
		CurrentPage:  pagination.Page,
		TotalPages:   totalPages,
	}, nil
}
