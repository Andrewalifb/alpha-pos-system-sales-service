package service

import (
	"context"
	"errors"
	"os"
	"time"

	pb "github.com/Andrewalifb/alpha-pos-system-sales-service/api/proto"
	"github.com/Andrewalifb/alpha-pos-system-sales-service/dto"
	"github.com/Andrewalifb/alpha-pos-system-sales-service/entity"
	"github.com/Andrewalifb/alpha-pos-system-sales-service/pkg/repository"
	"github.com/Andrewalifb/alpha-pos-system-sales-service/utils"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PosInvoiceService interface {
	CreatePosInvoice(ctx context.Context, req *pb.CreatePosInvoiceRequest) (*pb.CreatePosInvoiceResponse, error)
	ReadPosInvoice(ctx context.Context, req *pb.ReadPosInvoiceRequest) (*pb.ReadPosInvoiceResponse, error)
	UpdatePosInvoice(ctx context.Context, req *pb.UpdatePosInvoiceRequest) (*pb.UpdatePosInvoiceResponse, error)
	DeletePosInvoice(ctx context.Context, req *pb.DeletePosInvoiceRequest) (*pb.DeletePosInvoiceResponse, error)
	ReadAllPosInvoices(ctx context.Context, req *pb.ReadAllPosInvoicesRequest) (*pb.ReadAllPosInvoicesResponse, error)
}

type posInvoiceService struct {
	pb.UnimplementedPosInvoiceServiceServer
	invoiceRepo repository.PosInvoiceRepository
}

func NewPosInvoiceService(invoiceRepo repository.PosInvoiceRepository) *posInvoiceService {
	return &posInvoiceService{
		invoiceRepo: invoiceRepo,
	}
}

func (s *posInvoiceService) CreatePosInvoice(ctx context.Context, req *pb.CreatePosInvoiceRequest) (*pb.CreatePosInvoiceResponse, error) {
	// Extract role ID from JWT payload
	jwtRoleID := req.JwtPayload.Role
	token := req.JwtToken

	// Get user login role name
	loginRole, err := utils.GetPosRole(jwtRoleID, token)
	if err != nil {
		return nil, err
	}

	if !utils.IsBranchOrStoreUser(loginRole.Data.RoleName) {
		return nil, errors.New("users cant create invoice")
	}

	req.PosInvoice.InvoiceId = uuid.New().String() // Generate a new UUID for the invoice_id

	now := timestamppb.New(time.Now())
	req.PosInvoice.CreatedAt = now
	req.PosInvoice.UpdatedAt = now

	// Convert pb.PosInvoice to entity.PosInvoice
	gormInvoice := &entity.PosInvoice{
		InvoiceID: uuid.MustParse(req.PosInvoice.InvoiceId),
		ReceiptID: req.PosInvoice.ReceiptId,
		Date:      req.PosInvoice.Date.AsTime(),
		Amount:    req.PosInvoice.Amount,
		Discounts: req.PosInvoice.Discounts,
		Taxes:     req.PosInvoice.Taxes,
		BranchID:  uuid.MustParse(req.JwtPayload.BranchId),
		CompanyID: uuid.MustParse(req.JwtPayload.CompanyId),
		CreatedAt: req.PosInvoice.CreatedAt.AsTime(),
		CreatedBy: uuid.MustParse(req.JwtPayload.UserId),
		UpdatedAt: req.PosInvoice.UpdatedAt.AsTime(),
		UpdatedBy: uuid.MustParse(req.JwtPayload.UserId),
	}

	err = s.invoiceRepo.CreatePosInvoice(gormInvoice)
	if err != nil {
		return nil, err
	}

	return &pb.CreatePosInvoiceResponse{
		PosInvoice: req.PosInvoice,
	}, nil
}

func (s *posInvoiceService) ReadAllPosInvoices(ctx context.Context, req *pb.ReadAllPosInvoicesRequest) (*pb.ReadAllPosInvoicesResponse, error) {
	pagination := dto.Pagination{
		Limit: int(req.Limit),
		Page:  int(req.Page),
	}

	// Extract role ID from JWT payload
	jwtRoleID := req.JwtPayload.Role
	token := req.JwtToken

	// Get user login role name
	loginRole, err := utils.GetPosRole(jwtRoleID, token)
	if err != nil {
		return nil, err
	}

	if !utils.IsCompanyOrBranchUser(loginRole.Data.RoleName) {
		return nil, errors.New("users cant read all invoice")
	}

	paginationResult, err := s.invoiceRepo.ReadAllPosInvoices(pagination, loginRole.Data.RoleName, req.JwtPayload)
	if err != nil {
		return nil, err
	}

	posInvoices := paginationResult.Records.([]entity.PosInvoice)
	pbPosInvoices := make([]*pb.PosInvoice, len(posInvoices))

	for i, posInvoice := range posInvoices {
		pbPosInvoices[i] = &pb.PosInvoice{
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
	}

	return &pb.ReadAllPosInvoicesResponse{
		PosInvoices: pbPosInvoices,
		Limit:       int32(pagination.Limit),
		Page:        int32(pagination.Page),
		MaxPage:     int32(paginationResult.TotalPages),
		Count:       paginationResult.TotalRecords,
	}, nil
}

func (s *posInvoiceService) ReadPosInvoice(ctx context.Context, req *pb.ReadPosInvoiceRequest) (*pb.ReadPosInvoiceResponse, error) {
	// Extract role ID from JWT payload
	jwtRoleID := req.JwtPayload.Role
	token := req.JwtToken

	// Get user login role name
	loginRole, err := utils.GetPosRole(jwtRoleID, token)
	if err != nil {
		return nil, err
	}

	if !utils.IsCompanyOrBranchOrStoreUser(loginRole.Data.RoleName) {
		return nil, errors.New("users cant read invoice")
	}

	posInvoice, err := s.invoiceRepo.ReadPosInvoice(req.InvoiceId)
	if err != nil {
		return nil, err
	}

	companyRole := os.Getenv("COMPANY_USER_ROLE")
	branchRole := os.Getenv("BRANCH_USER_ROLE")
	storeRole := os.Getenv("STORE_USER_ROLE")

	if loginRole.Data.RoleName == companyRole {
		if !utils.VerifyCompanyUserAccess(loginRole.Data.RoleName, posInvoice.CompanyId, req.JwtPayload.CompanyId) {
			return nil, errors.New("company users can only retrieve invoice data within their company")
		}
	}

	if loginRole.Data.RoleName == branchRole {
		if !utils.VerifyBranchUserAccess(loginRole.Data.RoleName, posInvoice.BranchId, req.JwtPayload.BranchId) {
			return nil, errors.New("branch users can only retrieve invoice data within their branch")
		}
	}

	if loginRole.Data.RoleName == storeRole {
		if !utils.VerifyStoreUserAccess(loginRole.Data.RoleName, posInvoice.BranchId, req.JwtPayload.BranchId) {
			return nil, errors.New("store users can only retrieve invoice data within their branch")
		}
	}

	return &pb.ReadPosInvoiceResponse{
		PosInvoice: posInvoice,
	}, nil
}

func (s *posInvoiceService) UpdatePosInvoice(ctx context.Context, req *pb.UpdatePosInvoiceRequest) (*pb.UpdatePosInvoiceResponse, error) {
	// Get the role name from the role ID in the JWT payload
	// Extract role ID from JWT payload
	jwtRoleID := req.JwtPayload.Role
	token := req.JwtToken

	// Get user login role name
	loginRole, err := utils.GetPosRole(jwtRoleID, token)
	if err != nil {
		return nil, err
	}

	if !utils.IsBranchUser(loginRole.Data.RoleName) {
		return nil, errors.New("users cant update invoice data")
	}

	// Get the invoice to be updated
	posInvoice, err := s.invoiceRepo.ReadPosInvoice(req.PosInvoice.InvoiceId)
	if err != nil {
		return nil, err
	}

	branchRole := os.Getenv("BRANCH_USER_ROLE")

	if loginRole.Data.RoleName == branchRole {
		if !utils.VerifyBranchUserAccess(loginRole.Data.RoleName, posInvoice.BranchId, req.JwtPayload.BranchId) {
			return nil, errors.New("branch users can only update invoice data within their branch")
		}
	}

	now := timestamppb.New(time.Now())
	req.PosInvoice.UpdatedAt = now

	newInvoiceData := &entity.PosInvoice{
		InvoiceID: uuid.MustParse(posInvoice.InvoiceId),
		ReceiptID: req.PosInvoice.ReceiptId,
		Date:      posInvoice.Date.AsTime(),
		Amount:    req.PosInvoice.Amount,
		Discounts: req.PosInvoice.Discounts,
		Taxes:     req.PosInvoice.Taxes,
		BranchID:  uuid.MustParse(posInvoice.BranchId),
		CompanyID: uuid.MustParse(posInvoice.CompanyId),
		CreatedAt: posInvoice.CreatedAt.AsTime(),
		CreatedBy: uuid.MustParse(posInvoice.CreatedBy),
		UpdatedAt: req.PosInvoice.UpdatedAt.AsTime(),
		UpdatedBy: uuid.MustParse(req.JwtPayload.UserId),
	}

	// Update the invoice
	posInvoice, err = s.invoiceRepo.UpdatePosInvoice(newInvoiceData)
	if err != nil {
		return nil, err
	}

	return &pb.UpdatePosInvoiceResponse{
		PosInvoice: posInvoice,
	}, nil
}

func (s *posInvoiceService) DeletePosInvoice(ctx context.Context, req *pb.DeletePosInvoiceRequest) (*pb.DeletePosInvoiceResponse, error) {
	// Extract role ID from JWT payload
	jwtRoleID := req.JwtPayload.Role
	token := req.JwtToken

	// Get user login role name
	loginRole, err := utils.GetPosRole(jwtRoleID, token)
	if err != nil {
		return nil, err
	}

	if !utils.IsBranchUser(loginRole.Data.RoleName) {
		return nil, errors.New("users cant update invoice data")
	}

	// Get the invoice to be updated
	posInvoice, err := s.invoiceRepo.ReadPosInvoice(req.InvoiceId)
	if err != nil {
		return nil, err
	}

	branchRole := os.Getenv("BRANCH_USER_ROLE")

	if loginRole.Data.RoleName == branchRole {
		if !utils.VerifyBranchUserAccess(loginRole.Data.RoleName, posInvoice.BranchId, req.JwtPayload.BranchId) {
			return nil, errors.New("branch users can only delete invoice data within their branch")
		}
	}

	// Delete the invoice
	err = s.invoiceRepo.DeletePosInvoice(req.InvoiceId)
	if err != nil {
		return nil, err
	}

	return &pb.DeletePosInvoiceResponse{
		Success: true,
	}, nil
}
