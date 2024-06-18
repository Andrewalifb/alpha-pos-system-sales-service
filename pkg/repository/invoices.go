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

type PosInvoiceRepository interface {
	CreatePosInvoice(posInvoice *entity.PosInvoice) error
	ReadPosInvoice(invoiceID string) (*pb.PosInvoice, error)
	UpdatePosInvoice(posInvoice *entity.PosInvoice) (*pb.PosInvoice, error)
	DeletePosInvoice(invoiceID string) error
	ReadAllPosInvoices(pagination dto.Pagination, roleName string, jwtPayload *pb.JWTPayload) (*dto.PaginationResult, error)
}

type posInvoiceRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewPosInvoiceRepository(db *gorm.DB, redis *redis.Client) PosInvoiceRepository {
	return &posInvoiceRepository{
		db:    db,
		redis: redis,
	}
}

func (r *posInvoiceRepository) CreatePosInvoice(posInvoice *entity.PosInvoice) error {
	result := r.db.Create(posInvoice)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *posInvoiceRepository) ReadAllPosInvoices(pagination dto.Pagination, roleName string, jwtPayload *pb.JWTPayload) (*dto.PaginationResult, error) {
	var posInvoices []entity.PosInvoice
	var totalRecords int64

	query := r.db.Model(&entity.PosInvoice{})

	companyRole := os.Getenv("COMPANY_USER_ROLE")
	branchRole := os.Getenv("BRANCH_USER_ROLE")
	storeRole := os.Getenv("STORE_USER_ROLE")

	switch roleName {
	case companyRole:
		query = query.Where("company_id = ?", jwtPayload.CompanyId)
	case branchRole:
		query = query.Where("branch_id = ?", jwtPayload.BranchId)
	case storeRole:
		return nil, errors.New("store users are not allowed to retrieve invoices")
	default:
		return nil, errors.New("invalid role")
	}

	if pagination.Limit > 0 && pagination.Page > 0 {
		offset := (pagination.Page - 1) * pagination.Limit
		query = query.Offset(offset).Limit(pagination.Limit)
	}

	query.Find(&posInvoices)
	query.Count(&totalRecords)

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pagination.Limit)))

	return &dto.PaginationResult{
		TotalRecords: totalRecords,
		Records:      posInvoices,
		CurrentPage:  pagination.Page,
		TotalPages:   totalPages,
	}, nil
}

func (r *posInvoiceRepository) ReadPosInvoice(invoiceID string) (*pb.PosInvoice, error) {
	// Try to get the invoice from Redis first
	invoiceData, err := r.redis.Get(context.Background(), invoiceID).Result()
	if err == redis.Nil {
		// Invoice not found in Redis, get from PostgreSQL
		var posInvoiceEntity entity.PosInvoice
		if err := r.db.Where("invoice_id = ?", invoiceID).First(&posInvoiceEntity).Error; err != nil {
			return nil, err
		}

		// Convert entity.PosInvoice to pb.PosInvoice
		posInvoice := &pb.PosInvoice{
			InvoiceId: posInvoiceEntity.InvoiceID.String(),
			ReceiptId: posInvoiceEntity.ReceiptID,
			Date:      timestamppb.New(posInvoiceEntity.Date),
			Amount:    posInvoiceEntity.Amount,
			Discounts: posInvoiceEntity.Discounts,
			Taxes:     posInvoiceEntity.Taxes,
			BranchId:  posInvoiceEntity.BranchID.String(),
			CompanyId: posInvoiceEntity.CompanyID.String(),
			CreatedAt: timestamppb.New(posInvoiceEntity.CreatedAt),
			CreatedBy: posInvoiceEntity.CreatedBy.String(),
			UpdatedAt: timestamppb.New(posInvoiceEntity.UpdatedAt),
			UpdatedBy: posInvoiceEntity.UpdatedBy.String(),
		}

		// Store the invoice in Redis for future queries
		invoiceData, err := json.Marshal(posInvoiceEntity)
		if err != nil {
			return nil, err
		}
		err = r.redis.Set(context.Background(), invoiceID, invoiceData, 7*24*time.Hour).Err()
		if err != nil {
			return nil, err
		}

		return posInvoice, nil
	} else if err != nil {
		return nil, err
	}

	// Invoice found in Redis, unmarshal the data
	var posInvoiceEntity entity.PosInvoice
	err = json.Unmarshal([]byte(invoiceData), &posInvoiceEntity)
	if err != nil {
		return nil, err
	}

	// Convert entity.PosInvoice to pb.PosInvoice
	posInvoice := &pb.PosInvoice{
		InvoiceId: posInvoiceEntity.InvoiceID.String(),
		ReceiptId: posInvoiceEntity.ReceiptID,
		Date:      timestamppb.New(posInvoiceEntity.Date),
		Amount:    posInvoiceEntity.Amount,
		Discounts: posInvoiceEntity.Discounts,
		Taxes:     posInvoiceEntity.Taxes,
		BranchId:  posInvoiceEntity.BranchID.String(),
		CompanyId: posInvoiceEntity.CompanyID.String(),
		CreatedAt: timestamppb.New(posInvoiceEntity.CreatedAt),
		CreatedBy: posInvoiceEntity.CreatedBy.String(),
		UpdatedAt: timestamppb.New(posInvoiceEntity.UpdatedAt),
		UpdatedBy: posInvoiceEntity.UpdatedBy.String(),
	}

	return posInvoice, nil
}

func (r *posInvoiceRepository) UpdatePosInvoice(posInvoice *entity.PosInvoice) (*pb.PosInvoice, error) {
	if err := r.db.Save(posInvoice).Error; err != nil {
		return nil, err
	}

	// Convert updated entity.PosInvoice back to pb.PosInvoice
	updatedPosInvoice := &pb.PosInvoice{
		InvoiceId: posInvoice.InvoiceID.String(),
		ReceiptId: posInvoice.ReceiptID,
		Date:      timestamppb.New(posInvoice.Date),
		Amount:    posInvoice.Amount,
		Discounts: posInvoice.Discounts,
		Taxes:     posInvoice.Taxes,
		BranchId:  posInvoice.BranchID.String(),
		CompanyId: posInvoice.CompanyID.String(),
		CreatedAt: timestamppb.New(posInvoice.CreatedAt),
		CreatedBy: posInvoice.CreatedBy.String(),
		UpdatedAt: timestamppb.New(posInvoice.UpdatedAt),
		UpdatedBy: posInvoice.UpdatedBy.String(),
	}

	// Update the invoice in Redis
	invoiceData, err := json.Marshal(posInvoice)
	if err != nil {
		return nil, err
	}
	err = r.redis.Set(context.Background(), updatedPosInvoice.InvoiceId, invoiceData, 7*24*time.Hour).Err()
	if err != nil {
		return nil, err
	}

	return updatedPosInvoice, nil
}

func (r *posInvoiceRepository) DeletePosInvoice(invoiceID string) error {
	if err := r.db.Where("invoice_id = ?", invoiceID).Delete(&pb.PosInvoice{}).Error; err != nil {
		return err
	}

	// Delete the invoice from Redis
	err := r.redis.Del(context.Background(), invoiceID).Err()
	if err != nil {
		return err
	}

	return nil
}
