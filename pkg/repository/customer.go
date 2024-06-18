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

type PosCustomerRepository interface {
	CreatePosCustomer(posCustomer *entity.PosCustomer) error
	ReadPosCustomer(customerID string) (*pb.PosCustomer, error)
	UpdatePosCustomer(posCustomer *entity.PosCustomer) (*pb.PosCustomer, error)
	DeletePosCustomer(customerID string) error
	ReadAllPosCustomers(pagination dto.Pagination, roleName string, jwtPayload *pb.JWTPayload) (*dto.PaginationResult, error)
}

type posCustomerRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewPosCustomerRepository(db *gorm.DB, redis *redis.Client) PosCustomerRepository {
	return &posCustomerRepository{
		db:    db,
		redis: redis,
	}
}

func (r *posCustomerRepository) CreatePosCustomer(posCustomer *entity.PosCustomer) error {
	result := r.db.Create(posCustomer)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *posCustomerRepository) ReadAllPosCustomers(pagination dto.Pagination, roleName string, jwtPayload *pb.JWTPayload) (*dto.PaginationResult, error) {
	var posCustomers []entity.PosCustomer
	var totalRecords int64

	query := r.db.Model(&entity.PosCustomer{})

	companyRole := os.Getenv("COMPANY_USER_ROLE")
	branchRole := os.Getenv("BRANCH_USER_ROLE")
	storeRole := os.Getenv("STORE_USER_ROLE")

	switch roleName {
	case companyRole:
		query = query.Where("company_id = ?", jwtPayload.CompanyId)
	case branchRole:
		query = query.Where("branch_id = ?", jwtPayload.BranchId)
	case storeRole:
		query = query.Where("branch_id = ?", jwtPayload.BranchId)
	default:
		return nil, errors.New("invalid role")
	}

	if pagination.Limit > 0 && pagination.Page > 0 {
		offset := (pagination.Page - 1) * pagination.Limit
		query = query.Offset(offset).Limit(pagination.Limit)
	}

	query.Find(&posCustomers)
	query.Count(&totalRecords)

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pagination.Limit)))

	return &dto.PaginationResult{
		TotalRecords: totalRecords,
		Records:      posCustomers,
		CurrentPage:  pagination.Page,
		TotalPages:   totalPages,
	}, nil
}

func (r *posCustomerRepository) ReadPosCustomer(customerID string) (*pb.PosCustomer, error) {
	// Try to get the customer from Redis first
	customerData, err := r.redis.Get(context.Background(), customerID).Result()
	if err == redis.Nil {
		// Customer not found in Redis, get from PostgreSQL
		var posCustomerEntity entity.PosCustomer
		if err := r.db.Where("customer_id = ?", customerID).First(&posCustomerEntity).Error; err != nil {
			return nil, err
		}

		// Convert entity.PosCustomer to pb.PosCustomer
		posCustomer := &pb.PosCustomer{
			CustomerId:       posCustomerEntity.CustomerID.String(),
			FirstName:        posCustomerEntity.FirstName,
			LastName:         posCustomerEntity.LastName,
			Email:            posCustomerEntity.Email,
			PhoneNumber:      posCustomerEntity.PhoneNumber,
			DateOfBirth:      posCustomerEntity.DateOfBirth.Format(time.RFC3339),
			RegistrationDate: timestamppb.New(posCustomerEntity.RegistrationDate),
			Address:          posCustomerEntity.Address,
			City:             posCustomerEntity.City,
			Country:          posCustomerEntity.Country,
			BranchId:         posCustomerEntity.BranchID.String(),
			CompanyId:        posCustomerEntity.CompanyID.String(),
			CreatedAt:        timestamppb.New(posCustomerEntity.CreatedAt),
			CreatedBy:        posCustomerEntity.CreatedBy.String(),
			UpdatedAt:        timestamppb.New(posCustomerEntity.UpdatedAt),
			UpdatedBy:        posCustomerEntity.UpdatedBy.String(),
		}

		// Store the customer in Redis for future queries
		customerData, err := json.Marshal(posCustomerEntity)
		if err != nil {
			return nil, err
		}
		err = r.redis.Set(context.Background(), customerID, customerData, 7*24*time.Hour).Err()
		if err != nil {
			return nil, err
		}

		return posCustomer, nil
	} else if err != nil {
		return nil, err
	}

	// Customer found in Redis, unmarshal the data
	var posCustomerEntity entity.PosCustomer
	err = json.Unmarshal([]byte(customerData), &posCustomerEntity)
	if err != nil {
		return nil, err
	}

	// Convert entity.PosCustomer to pb.PosCustomer
	posCustomer := &pb.PosCustomer{
		CustomerId:       posCustomerEntity.CustomerID.String(),
		FirstName:        posCustomerEntity.FirstName,
		LastName:         posCustomerEntity.LastName,
		Email:            posCustomerEntity.Email,
		PhoneNumber:      posCustomerEntity.PhoneNumber,
		DateOfBirth:      posCustomerEntity.DateOfBirth.Format(time.RFC3339),
		RegistrationDate: timestamppb.New(posCustomerEntity.RegistrationDate),
		Address:          posCustomerEntity.Address,
		City:             posCustomerEntity.City,
		Country:          posCustomerEntity.Country,
		BranchId:         posCustomerEntity.BranchID.String(),
		CompanyId:        posCustomerEntity.CompanyID.String(),
		CreatedAt:        timestamppb.New(posCustomerEntity.CreatedAt),
		CreatedBy:        posCustomerEntity.CreatedBy.String(),
		UpdatedAt:        timestamppb.New(posCustomerEntity.UpdatedAt),
		UpdatedBy:        posCustomerEntity.UpdatedBy.String(),
	}

	return posCustomer, nil
}

func (r *posCustomerRepository) UpdatePosCustomer(posCustomer *entity.PosCustomer) (*pb.PosCustomer, error) {
	if err := r.db.Save(posCustomer).Error; err != nil {
		return nil, err
	}

	// Convert updated entity.PosCustomer back to pb.PosCustomer
	updatedPosCustomer := &pb.PosCustomer{
		CustomerId:       posCustomer.CustomerID.String(),
		FirstName:        posCustomer.FirstName,
		LastName:         posCustomer.LastName,
		Email:            posCustomer.Email,
		PhoneNumber:      posCustomer.PhoneNumber,
		DateOfBirth:      posCustomer.DateOfBirth.Format(time.RFC3339),
		RegistrationDate: timestamppb.New(posCustomer.RegistrationDate),
		Address:          posCustomer.Address,
		City:             posCustomer.City,
		Country:          posCustomer.Country,
		BranchId:         posCustomer.BranchID.String(),
		CompanyId:        posCustomer.CompanyID.String(),
		CreatedAt:        timestamppb.New(posCustomer.CreatedAt),
		CreatedBy:        posCustomer.CreatedBy.String(),
		UpdatedAt:        timestamppb.New(posCustomer.UpdatedAt),
		UpdatedBy:        posCustomer.UpdatedBy.String(),
	}

	// Update the customer in Redis
	customerData, err := json.Marshal(posCustomer)
	if err != nil {
		return nil, err
	}
	err = r.redis.Set(context.Background(), updatedPosCustomer.CustomerId, customerData, 7*24*time.Hour).Err()
	if err != nil {
		return nil, err
	}

	return updatedPosCustomer, nil
}

func (r *posCustomerRepository) DeletePosCustomer(customerID string) error {
	if err := r.db.Where("customer_id = ?", customerID).Delete(&pb.PosCustomer{}).Error; err != nil {
		return err
	}

	// Delete the customer from Redis
	err := r.redis.Del(context.Background(), customerID).Err()
	if err != nil {
		return err
	}

	return nil
}
