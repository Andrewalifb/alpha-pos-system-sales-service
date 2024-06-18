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

type PosCashDrawerRepository interface {
	CreatePosCashDrawer(posCashDrawer *entity.PosCashDrawer) error
	ReadPosCashDrawer(drawerID string) (*pb.PosCashDrawer, error)
	UpdatePosCashDrawer(posCashDrawer *entity.PosCashDrawer) (*pb.PosCashDrawer, error)
	DeletePosCashDrawer(drawerID string) error
	ReadAllPosCashDrawers(pagination dto.Pagination, roleName string, jwtPayload *pb.JWTPayload) (*dto.PaginationResult, error)
}

type posCashDrawerRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewPosCashDrawerRepository(db *gorm.DB, redis *redis.Client) PosCashDrawerRepository {
	return &posCashDrawerRepository{
		db:    db,
		redis: redis,
	}
}

func (r *posCashDrawerRepository) CreatePosCashDrawer(posCashDrawer *entity.PosCashDrawer) error {
	result := r.db.Create(posCashDrawer)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *posCashDrawerRepository) ReadAllPosCashDrawers(pagination dto.Pagination, roleName string, jwtPayload *pb.JWTPayload) (*dto.PaginationResult, error) {
	var posCashDrawers []entity.PosCashDrawer
	var totalRecords int64

	query := r.db.Model(&entity.PosCashDrawer{})

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

	query.Find(&posCashDrawers)
	query.Count(&totalRecords)

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pagination.Limit)))

	return &dto.PaginationResult{
		TotalRecords: totalRecords,
		Records:      posCashDrawers,
		CurrentPage:  pagination.Page,
		TotalPages:   totalPages,
	}, nil
}

func (r *posCashDrawerRepository) ReadPosCashDrawer(drawerID string) (*pb.PosCashDrawer, error) {
	// Try to get the cash drawer from Redis first
	cashDrawerData, err := r.redis.Get(context.Background(), drawerID).Result()
	if err == redis.Nil {
		// Cash drawer not found in Redis, get from PostgreSQL
		var posCashDrawerEntity entity.PosCashDrawer
		if err := r.db.Where("drawer_id = ?", drawerID).First(&posCashDrawerEntity).Error; err != nil {
			return nil, err
		}

		// Convert entity.PosCashDrawer to pb.PosCashDrawer
		posCashDrawer := &pb.PosCashDrawer{
			DrawerId:        posCashDrawerEntity.DrawerID.String(),
			StoreId:         posCashDrawerEntity.StoreID.String(),
			EmployeeId:      posCashDrawerEntity.EmployeeID.String(),
			ReceiptId:       posCashDrawerEntity.ReceiptID,
			CashIn:          posCashDrawerEntity.CashIn,
			Amount:          posCashDrawerEntity.Amount,
			CashOut:         posCashDrawerEntity.CashOut,
			TransactionTime: timestamppb.New(posCashDrawerEntity.TransactionTime),
			RoleId:          posCashDrawerEntity.RoleID.String(),
			BranchId:        posCashDrawerEntity.BranchID.String(),
			CompanyId:       posCashDrawerEntity.CompanyID.String(),
			Description:     posCashDrawerEntity.Description,
			CreatedAt:       timestamppb.New(posCashDrawerEntity.CreatedAt),
			CreatedBy:       posCashDrawerEntity.CreatedBy.String(),
			UpdatedAt:       timestamppb.New(posCashDrawerEntity.UpdatedAt),
			UpdatedBy:       posCashDrawerEntity.UpdatedBy.String(),
		}

		// Store the cash drawer in Redis for future queries
		cashDrawerData, err := json.Marshal(posCashDrawerEntity)
		if err != nil {
			return nil, err
		}
		err = r.redis.Set(context.Background(), drawerID, cashDrawerData, 7*24*time.Hour).Err()
		if err != nil {
			return nil, err
		}

		return posCashDrawer, nil
	} else if err != nil {
		return nil, err
	}

	// Cash drawer found in Redis, unmarshal the data
	var posCashDrawerEntity entity.PosCashDrawer
	err = json.Unmarshal([]byte(cashDrawerData), &posCashDrawerEntity)
	if err != nil {
		return nil, err
	}

	// Convert entity.PosCashDrawer to pb.PosCashDrawer
	posCashDrawer := &pb.PosCashDrawer{
		DrawerId:        posCashDrawerEntity.DrawerID.String(),
		StoreId:         posCashDrawerEntity.StoreID.String(),
		EmployeeId:      posCashDrawerEntity.EmployeeID.String(),
		ReceiptId:       posCashDrawerEntity.ReceiptID,
		CashIn:          posCashDrawerEntity.CashIn,
		Amount:          posCashDrawerEntity.Amount,
		CashOut:         posCashDrawerEntity.CashOut,
		TransactionTime: timestamppb.New(posCashDrawerEntity.TransactionTime),
		RoleId:          posCashDrawerEntity.RoleID.String(),
		BranchId:        posCashDrawerEntity.BranchID.String(),
		CompanyId:       posCashDrawerEntity.CompanyID.String(),
		Description:     posCashDrawerEntity.Description,
		CreatedAt:       timestamppb.New(posCashDrawerEntity.CreatedAt),
		CreatedBy:       posCashDrawerEntity.CreatedBy.String(),
		UpdatedAt:       timestamppb.New(posCashDrawerEntity.UpdatedAt),
		UpdatedBy:       posCashDrawerEntity.UpdatedBy.String(),
	}

	return posCashDrawer, nil
}

func (r *posCashDrawerRepository) UpdatePosCashDrawer(posCashDrawer *entity.PosCashDrawer) (*pb.PosCashDrawer, error) {
	if err := r.db.Save(posCashDrawer).Error; err != nil {
		return nil, err
	}

	// Convert updated entity.PosCashDrawer back to pb.PosCashDrawer
	updatedPosCashDrawer := &pb.PosCashDrawer{
		DrawerId:        posCashDrawer.DrawerID.String(),
		StoreId:         posCashDrawer.StoreID.String(),
		EmployeeId:      posCashDrawer.EmployeeID.String(),
		ReceiptId:       posCashDrawer.ReceiptID,
		CashIn:          posCashDrawer.CashIn,
		Amount:          posCashDrawer.Amount,
		CashOut:         posCashDrawer.CashOut,
		TransactionTime: timestamppb.New(posCashDrawer.TransactionTime),
		RoleId:          posCashDrawer.RoleID.String(),
		BranchId:        posCashDrawer.BranchID.String(),
		CompanyId:       posCashDrawer.CompanyID.String(),
		Description:     posCashDrawer.Description,
		CreatedAt:       timestamppb.New(posCashDrawer.CreatedAt),
		CreatedBy:       posCashDrawer.CreatedBy.String(),
		UpdatedAt:       timestamppb.New(posCashDrawer.UpdatedAt),
		UpdatedBy:       posCashDrawer.UpdatedBy.String(),
	}

	// Update the cash drawer in Redis
	cashDrawerData, err := json.Marshal(posCashDrawer)
	if err != nil {
		return nil, err
	}
	err = r.redis.Set(context.Background(), updatedPosCashDrawer.DrawerId, cashDrawerData, 7*24*time.Hour).Err()
	if err != nil {
		return nil, err
	}

	return updatedPosCashDrawer, nil
}

func (r *posCashDrawerRepository) DeletePosCashDrawer(drawerID string) error {
	if err := r.db.Where("drawer_id = ?", drawerID).Delete(&pb.PosCashDrawer{}).Error; err != nil {
		return err
	}

	// Delete the cash drawer from Redis
	err := r.redis.Del(context.Background(), drawerID).Err()
	if err != nil {
		return err
	}

	return nil
}
