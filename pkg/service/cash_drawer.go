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

type PosCashDrawerService interface {
	CreatePosCashDrawer(ctx context.Context, req *pb.CreatePosCashDrawerRequest) (*pb.CreatePosCashDrawerResponse, error)
	ReadPosCashDrawer(ctx context.Context, req *pb.ReadPosCashDrawerRequest) (*pb.ReadPosCashDrawerResponse, error)
	UpdatePosCashDrawer(ctx context.Context, req *pb.UpdatePosCashDrawerRequest) (*pb.UpdatePosCashDrawerResponse, error)
	DeletePosCashDrawer(ctx context.Context, req *pb.DeletePosCashDrawerRequest) (*pb.DeletePosCashDrawerResponse, error)
	ReadAllPosCashDrawers(ctx context.Context, req *pb.ReadAllPosCashDrawersRequest) (*pb.ReadAllPosCashDrawersResponse, error)
}

type PosCashDrawerServiceServer struct {
	pb.UnimplementedPosCashDrawerServiceServer
	cashDrawerRepo repository.PosCashDrawerRepository
}

func NewPosCashDrawerServiceServer(cashDrawerRepo repository.PosCashDrawerRepository) *PosCashDrawerServiceServer {
	return &PosCashDrawerServiceServer{
		cashDrawerRepo: cashDrawerRepo,
	}
}

func (s *PosCashDrawerServiceServer) CreatePosCashDrawer(ctx context.Context, req *pb.CreatePosCashDrawerRequest) (*pb.CreatePosCashDrawerResponse, error) {
	// Extract role ID from JWT payload
	jwtRoleID := req.JwtPayload.Role
	token := req.JwtToken

	// Get user login role name
	loginRole, err := utils.GetPosRole(jwtRoleID, token)
	if err != nil {
		return nil, err
	}

	if !utils.IsBranchOrStoreUser(loginRole.Data.RoleName) {
		return nil, errors.New("users cant create cash drawer")
	}

	req.PosCashDrawer.DrawerId = uuid.New().String() // Generate a new UUID for the drawer_id

	now := timestamppb.New(time.Now())
	req.PosCashDrawer.CreatedAt = now
	req.PosCashDrawer.UpdatedAt = now
	req.PosCashDrawer.TransactionTime = now
	// Convert pb.PosCashDrawer to entity.PosCashDrawer
	gormCashDrawer := &entity.PosCashDrawer{
		DrawerID:        uuid.MustParse(req.PosCashDrawer.DrawerId), // auto
		StoreID:         nil,
		EmployeeID:      uuid.MustParse(req.JwtPayload.UserId),
		ReceiptID:       req.PosCashDrawer.ReceiptId,
		CashIn:          req.PosCashDrawer.CashIn,
		Amount:          req.PosCashDrawer.Amount,
		CashOut:         req.PosCashDrawer.CashOut,
		TransactionTime: req.PosCashDrawer.TransactionTime.AsTime(), // auto
		RoleID:          uuid.MustParse(req.PosCashDrawer.RoleId),
		BranchID:        nil,
		CompanyID:       uuid.MustParse(req.JwtPayload.CompanyId), // auto
		Description:     req.PosCashDrawer.Description,
		CreatedAt:       req.PosCashDrawer.CreatedAt.AsTime(),  // auto
		CreatedBy:       uuid.MustParse(req.JwtPayload.UserId), // auto
		UpdatedAt:       req.PosCashDrawer.UpdatedAt.AsTime(),  // auto
		UpdatedBy:       uuid.MustParse(req.JwtPayload.UserId), // auto
	}

	branchRole := os.Getenv("BRANCH_USER_ROLE")
	storeRole := os.Getenv("STORE_USER_ROLE")

	// set Branch ID base in login role
	switch loginRole.Data.RoleName {
	case branchRole:
		gormCashDrawer.BranchID = utils.ParseUUID(req.JwtPayload.BranchId)
		gormCashDrawer.StoreID = utils.ParseUUID(req.PosCashDrawer.StoreId)

		if gormCashDrawer.StoreID == nil {
			return nil, errors.New("error created cash drawer, store id could not be empty")
		}
	case storeRole:
		gormCashDrawer.BranchID = utils.ParseUUID(req.JwtPayload.BranchId)
		gormCashDrawer.StoreID = utils.ParseUUID(req.JwtPayload.StoreId)
	}

	err = s.cashDrawerRepo.CreatePosCashDrawer(gormCashDrawer)
	if err != nil {
		return nil, err
	}

	return &pb.CreatePosCashDrawerResponse{
		PosCashDrawer: req.PosCashDrawer,
	}, nil
}

func (s *PosCashDrawerServiceServer) ReadAllPosCashDrawers(ctx context.Context, req *pb.ReadAllPosCashDrawersRequest) (*pb.ReadAllPosCashDrawersResponse, error) {
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

	if !utils.IsCompanyOrBranchOrStoreUser(loginRole.Data.RoleName) {
		return nil, errors.New("users cant read all cash drawer")
	}

	paginationResult, err := s.cashDrawerRepo.ReadAllPosCashDrawers(pagination, loginRole.Data.RoleName, req.JwtPayload)
	if err != nil {
		return nil, err
	}

	posCashDrawers := paginationResult.Records.([]entity.PosCashDrawer)
	pbPosCashDrawers := make([]*pb.PosCashDrawer, len(posCashDrawers))

	for i, posCashDrawer := range posCashDrawers {
		pbPosCashDrawers[i] = &pb.PosCashDrawer{
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
	}

	return &pb.ReadAllPosCashDrawersResponse{
		PosCashDrawers: pbPosCashDrawers,
		Limit:          int32(pagination.Limit),
		Page:           int32(pagination.Page),
		MaxPage:        int32(paginationResult.TotalPages),
		Count:          paginationResult.TotalRecords,
	}, nil
}

func (s *PosCashDrawerServiceServer) ReadPosCashDrawer(ctx context.Context, req *pb.ReadPosCashDrawerRequest) (*pb.ReadPosCashDrawerResponse, error) {
	// Extract role ID from JWT payload
	jwtRoleID := req.JwtPayload.Role
	token := req.JwtToken

	// Get user login role name
	loginRole, err := utils.GetPosRole(jwtRoleID, token)
	if err != nil {
		return nil, err
	}

	if !utils.IsCompanyOrBranchOrStoreUser(loginRole.Data.RoleName) {
		return nil, errors.New("users cant read cash drawer")
	}

	posCashDrawer, err := s.cashDrawerRepo.ReadPosCashDrawer(req.DrawerId)
	if err != nil {
		return nil, err
	}

	companyRole := os.Getenv("COMPANY_USER_ROLE")
	branchRole := os.Getenv("BRANCH_USER_ROLE")
	storeRole := os.Getenv("STORE_USER_ROLE")

	if loginRole.Data.RoleName == companyRole {
		if !utils.VerifyCompanyUserAccess(loginRole.Data.RoleName, posCashDrawer.CompanyId, req.JwtPayload.CompanyId) {
			return nil, errors.New("company users can only retrieve cash drawer data within their company")
		}
	}

	if loginRole.Data.RoleName == branchRole {
		if !utils.VerifyBranchUserAccess(loginRole.Data.RoleName, posCashDrawer.BranchId, req.JwtPayload.BranchId) {
			return nil, errors.New("branch users can only retrieve cash drawer data within their branch")
		}
	}

	if loginRole.Data.RoleName == storeRole {
		if !utils.VerifyStoreUserAccess(loginRole.Data.RoleName, posCashDrawer.StoreId, req.JwtPayload.StoreId) {
			return nil, errors.New("store users can only retrieve cash drawer data within their store")
		}
	}

	return &pb.ReadPosCashDrawerResponse{
		PosCashDrawer: posCashDrawer,
	}, nil
}

func (s *PosCashDrawerServiceServer) UpdatePosCashDrawer(ctx context.Context, req *pb.UpdatePosCashDrawerRequest) (*pb.UpdatePosCashDrawerResponse, error) {
	// Extract role ID from JWT payload
	jwtRoleID := req.JwtPayload.Role
	token := req.JwtToken

	// Get user login role name
	loginRole, err := utils.GetPosRole(jwtRoleID, token)
	if err != nil {
		return nil, err
	}

	if !utils.IsBranchUser(loginRole.Data.RoleName) {
		return nil, errors.New("users cant update cash drawer")
	}

	// Get the cash drawer to be updated
	posCashDrawer, err := s.cashDrawerRepo.ReadPosCashDrawer(req.PosCashDrawer.DrawerId)
	if err != nil {
		return nil, err
	}

	branchRole := os.Getenv("BRANCH_USER_ROLE")

	if loginRole.Data.RoleName == branchRole {
		if !utils.VerifyBranchUserAccess(loginRole.Data.RoleName, posCashDrawer.BranchId, req.JwtPayload.BranchId) {
			return nil, errors.New("branch users can only update cash drawer data within their branch")
		}
	}

	now := timestamppb.New(time.Now())
	req.PosCashDrawer.UpdatedAt = now
	req.PosCashDrawer.TransactionTime = now

	newCashDrawerData := &entity.PosCashDrawer{
		DrawerID:        uuid.MustParse(posCashDrawer.DrawerId),   // auto
		StoreID:         nil,                                      // auto
		EmployeeID:      uuid.MustParse(posCashDrawer.EmployeeId), // auto
		ReceiptID:       posCashDrawer.ReceiptId,                  // aut0
		CashIn:          req.PosCashDrawer.CashIn,
		Amount:          req.PosCashDrawer.Amount,
		CashOut:         req.PosCashDrawer.CashOut,
		TransactionTime: posCashDrawer.TransactionTime.AsTime(),  // auto
		RoleID:          uuid.MustParse(posCashDrawer.RoleId),    // auto
		BranchID:        nil,                                     // auto
		CompanyID:       uuid.MustParse(posCashDrawer.CompanyId), // auto
		Description:     req.PosCashDrawer.Description,
		CreatedAt:       posCashDrawer.CreatedAt.AsTime(),        // auto
		CreatedBy:       uuid.MustParse(posCashDrawer.CreatedBy), // auto
		UpdatedAt:       req.PosCashDrawer.UpdatedAt.AsTime(),    // auto
		UpdatedBy:       uuid.MustParse(req.JwtPayload.UserId),   // auto
	}

	// set store and branch id
	newCashDrawerData.StoreID = utils.ParseUUID(posCashDrawer.StoreId)
	newCashDrawerData.BranchID = utils.ParseUUID(posCashDrawer.BranchId)

	// Update the cash drawer
	posCashDrawer, err = s.cashDrawerRepo.UpdatePosCashDrawer(newCashDrawerData)
	if err != nil {
		return nil, err
	}

	return &pb.UpdatePosCashDrawerResponse{
		PosCashDrawer: posCashDrawer,
	}, nil
}

func (s *PosCashDrawerServiceServer) DeletePosCashDrawer(ctx context.Context, req *pb.DeletePosCashDrawerRequest) (*pb.DeletePosCashDrawerResponse, error) {
	// Extract role ID from JWT payload
	jwtRoleID := req.JwtPayload.Role
	token := req.JwtToken

	// Get user login role name
	loginRole, err := utils.GetPosRole(jwtRoleID, token)
	if err != nil {
		return nil, err
	}

	if !utils.IsBranchUser(loginRole.Data.RoleName) {
		return nil, errors.New("users cant delete cash drawer")
	}

	// Get the cash drawer to be updated
	posCashDrawer, err := s.cashDrawerRepo.ReadPosCashDrawer(req.DrawerId)
	if err != nil {
		return nil, err
	}

	branchRole := os.Getenv("BRANCH_USER_ROLE")

	if loginRole.Data.RoleName == branchRole {
		if !utils.VerifyBranchUserAccess(loginRole.Data.RoleName, posCashDrawer.BranchId, req.JwtPayload.BranchId) {
			return nil, errors.New("branch users can only delete cash drawer data within their branch")
		}
	}

	// Delete the cash drawer
	err = s.cashDrawerRepo.DeletePosCashDrawer(req.DrawerId)
	if err != nil {
		return nil, err
	}

	return &pb.DeletePosCashDrawerResponse{
		Success: true,
	}, nil
}
