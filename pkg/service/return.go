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
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PosReturnService interface {
	CreatePosReturn(ctx context.Context, req *pb.CreatePosReturnRequest) (*pb.CreatePosReturnResponse, error)
	ReadPosReturn(ctx context.Context, req *pb.ReadPosReturnRequest) (*pb.ReadPosReturnResponse, error)
	UpdatePosReturn(ctx context.Context, req *pb.UpdatePosReturnRequest) (*pb.UpdatePosReturnResponse, error)
	DeletePosReturn(ctx context.Context, req *pb.DeletePosReturnRequest) (*pb.DeletePosReturnResponse, error)
	ReadAllPosReturns(ctx context.Context, req *pb.ReadAllPosReturnsRequest) (*pb.ReadAllPosReturnsResponse, error)
}

type posReturnService struct {
	pb.UnimplementedPosReturnServiceServer
	returnRepo         repository.PosReturnRepository
	CompanyServiceConn *grpc.ClientConn
}

func NewPosReturnService(returnRepo repository.PosReturnRepository, companyServiceConn *grpc.ClientConn) *posReturnService {
	return &posReturnService{
		returnRepo:         returnRepo,
		CompanyServiceConn: companyServiceConn,
	}
}

func (s *posReturnService) CreatePosReturn(ctx context.Context, req *pb.CreatePosReturnRequest) (*pb.CreatePosReturnResponse, error) {
	// Extract role ID from JWT payload
	jwtRoleID := req.JwtPayload.Role
	// token := req.JwtToken

	// Get user login role name
	loginRole, err := utils.GetPosRoleById(s.CompanyServiceConn, jwtRoleID, req.JwtPayload)
	if err != nil {
		return nil, err
	}

	if !utils.IsBranchOrStoreUser(loginRole.PosRole.RoleName) {
		return nil, errors.New("users are not allowed to create new return")
	}

	req.PosReturn.ReturnId = uuid.New().String() // Generate a new UUID for the return_id

	now := timestamppb.New(time.Now())
	req.PosReturn.CreatedAt = now
	req.PosReturn.UpdatedAt = now

	// Convert pb.PosReturn to entity.PosReturn
	gormReturn := &entity.PosReturn{
		ReturnID:   uuid.MustParse(req.PosReturn.ReturnId),
		ReceiptID:  req.PosReturn.ReceiptId,
		ProductID:  uuid.MustParse(req.PosReturn.ProductId),
		Quantity:   int(req.PosReturn.Quantity),
		Price:      req.PosReturn.Price,
		Amount:     req.PosReturn.Amount,
		ReturnDate: req.PosReturn.ReturnDate.AsTime(),
		Reason:     req.PosReturn.Reason,
		StoreID:    uuid.MustParse(req.PosReturn.StoreId),
		BranchID:   uuid.MustParse(req.JwtPayload.BranchId),
		CompanyID:  uuid.MustParse(req.JwtPayload.CompanyId),
		CreatedAt:  req.PosReturn.CreatedAt.AsTime(),
		CreatedBy:  uuid.MustParse(req.JwtPayload.UserId),
		UpdatedAt:  req.PosReturn.UpdatedAt.AsTime(),
		UpdatedBy:  uuid.MustParse(req.JwtPayload.UserId),
	}

	err = s.returnRepo.CreatePosReturn(gormReturn)
	if err != nil {
		return nil, err
	}

	return &pb.CreatePosReturnResponse{
		PosReturn: req.PosReturn,
	}, nil
}

func (s *posReturnService) ReadAllPosReturns(ctx context.Context, req *pb.ReadAllPosReturnsRequest) (*pb.ReadAllPosReturnsResponse, error) {
	pagination := dto.Pagination{
		Limit: int(req.Limit),
		Page:  int(req.Page),
	}

	// Extract role ID from JWT payload
	jwtRoleID := req.JwtPayload.Role
	// token := req.JwtToken

	// Get user login role name
	loginRole, err := utils.GetPosRoleById(s.CompanyServiceConn, jwtRoleID, req.JwtPayload)
	if err != nil {
		return nil, err
	}

	if !utils.IsCompanyOrBranchOrStoreUser(loginRole.PosRole.RoleName) {
		return nil, errors.New("users are not allowed to read all return")
	}

	paginationResult, err := s.returnRepo.ReadAllPosReturns(pagination, loginRole.PosRole.RoleName, req.JwtPayload)
	if err != nil {
		return nil, err
	}

	posReturns := paginationResult.Records.([]entity.PosReturn)
	pbPosReturns := make([]*pb.PosReturn, len(posReturns))

	for i, posReturn := range posReturns {
		pbPosReturns[i] = &pb.PosReturn{
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

	}

	return &pb.ReadAllPosReturnsResponse{
		PosReturns: pbPosReturns,
		Limit:      int32(pagination.Limit),
		Page:       int32(pagination.Page),
		MaxPage:    int32(paginationResult.TotalPages),
		Count:      paginationResult.TotalRecords,
	}, nil
}

func (s *posReturnService) ReadPosReturn(ctx context.Context, req *pb.ReadPosReturnRequest) (*pb.ReadPosReturnResponse, error) {
	// Extract role ID from JWT payload
	jwtRoleID := req.JwtPayload.Role
	// token := req.JwtToken

	// Get user login role name
	loginRole, err := utils.GetPosRoleById(s.CompanyServiceConn, jwtRoleID, req.JwtPayload)
	if err != nil {
		return nil, err
	}

	if !utils.IsCompanyOrBranchOrStoreUser(loginRole.PosRole.RoleName) {
		return nil, errors.New("users are not allowed to read return")
	}

	posReturn, err := s.returnRepo.ReadPosReturn(req.ReturnId)
	if err != nil {
		return nil, err
	}

	companyRole := os.Getenv("COMPANY_USER_ROLE")
	branchRole := os.Getenv("BRANCH_USER_ROLE")
	storeRole := os.Getenv("STORE_USER_ROLE")

	if loginRole.PosRole.RoleName == companyRole {
		if !utils.VerifyCompanyUserAccess(loginRole.PosRole.RoleName, posReturn.CompanyId, req.JwtPayload.CompanyId) {
			return nil, errors.New("company users can only retrieve return data within their company")
		}
	}

	if loginRole.PosRole.RoleName == branchRole {
		if !utils.VerifyBranchUserAccess(loginRole.PosRole.RoleName, posReturn.BranchId, req.JwtPayload.BranchId) {
			return nil, errors.New("branch users can only retrieve return data within their branch")
		}
	}

	if loginRole.PosRole.RoleName == storeRole {
		if !utils.VerifyStoreUserAccess(loginRole.PosRole.RoleName, posReturn.StoreId, req.JwtPayload.StoreId) {
			return nil, errors.New("store users can only retrieve return data within their branch")
		}
	}

	return &pb.ReadPosReturnResponse{
		PosReturn: posReturn,
	}, nil
}

func (s *posReturnService) UpdatePosReturn(ctx context.Context, req *pb.UpdatePosReturnRequest) (*pb.UpdatePosReturnResponse, error) {
	// Extract role ID from JWT payload
	jwtRoleID := req.JwtPayload.Role
	// token := req.JwtToken

	// Get user login role name
	loginRole, err := utils.GetPosRoleById(s.CompanyServiceConn, jwtRoleID, req.JwtPayload)
	if err != nil {
		return nil, err
	}

	if !utils.IsBranchUser(loginRole.PosRole.RoleName) {
		return nil, errors.New("users cant update return data")
	}

	// Get the return to be updated
	posReturn, err := s.returnRepo.ReadPosReturn(req.PosReturn.ReturnId)
	if err != nil {
		return nil, err
	}

	branchRole := os.Getenv("BRANCH_USER_ROLE")

	if loginRole.PosRole.RoleName == branchRole {
		if !utils.VerifyBranchUserAccess(loginRole.PosRole.RoleName, posReturn.BranchId, req.JwtPayload.BranchId) {
			return nil, errors.New("branch users can only update return data within their branch")
		}
	}

	now := timestamppb.New(time.Now())
	req.PosReturn.UpdatedAt = now

	newReturnData := &entity.PosReturn{
		ReturnID:   uuid.MustParse(posReturn.ReturnId),
		ReceiptID:  posReturn.ReceiptId,
		ProductID:  uuid.MustParse(req.PosReturn.ProductId),
		Quantity:   int(req.PosReturn.Quantity),
		Price:      posReturn.Price,
		Amount:     posReturn.Amount,
		ReturnDate: posReturn.ReturnDate.AsTime(),
		Reason:     req.PosReturn.Reason,
		StoreID:    uuid.MustParse(posReturn.StoreId),
		BranchID:   uuid.MustParse(req.JwtPayload.BranchId),
		CompanyID:  uuid.MustParse(req.JwtPayload.CompanyId),
		CreatedAt:  posReturn.CreatedAt.AsTime(),
		CreatedBy:  uuid.MustParse(posReturn.CreatedBy),
		UpdatedAt:  req.PosReturn.UpdatedAt.AsTime(),
		UpdatedBy:  uuid.MustParse(req.JwtPayload.UserId),
	}

	// Update the return
	posReturn, err = s.returnRepo.UpdatePosReturn(newReturnData)
	if err != nil {
		return nil, err
	}

	return &pb.UpdatePosReturnResponse{
		PosReturn: posReturn,
	}, nil
}

func (s *posReturnService) DeletePosReturn(ctx context.Context, req *pb.DeletePosReturnRequest) (*pb.DeletePosReturnResponse, error) {
	// Extract role ID from JWT payload
	jwtRoleID := req.JwtPayload.Role
	// token := req.JwtToken

	// Get user login role name
	loginRole, err := utils.GetPosRoleById(s.CompanyServiceConn, jwtRoleID, req.JwtPayload)
	if err != nil {
		return nil, err
	}

	if !utils.IsBranchUser(loginRole.PosRole.RoleName) {
		return nil, errors.New("users cant delete return data")
	}

	// Get the return to be updated
	posReturn, err := s.returnRepo.ReadPosReturn(req.ReturnId)
	if err != nil {
		return nil, err
	}

	branchRole := os.Getenv("BRANCH_USER_ROLE")

	if loginRole.PosRole.RoleName == branchRole {
		if !utils.VerifyBranchUserAccess(loginRole.PosRole.RoleName, posReturn.BranchId, req.JwtPayload.BranchId) {
			return nil, errors.New("branch users can only delete return data within their branch")
		}
	}

	// Delete the return
	err = s.returnRepo.DeletePosReturn(req.ReturnId)
	if err != nil {
		return nil, err
	}

	return &pb.DeletePosReturnResponse{
		Success: true,
	}, nil
}
