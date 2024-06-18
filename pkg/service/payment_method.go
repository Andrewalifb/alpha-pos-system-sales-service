package service

import (
	"context"
	"errors"
	"os"
	"time"

	pb "github.com/Andrewalifb/alpha-pos-system-sales-service/api/proto"
	"github.com/Andrewalifb/alpha-pos-system-sales-service/pkg/repository"
	"github.com/Andrewalifb/alpha-pos-system-sales-service/utils"

	"github.com/Andrewalifb/alpha-pos-system-sales-service/dto"
	"github.com/Andrewalifb/alpha-pos-system-sales-service/entity"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PosPaymentMethodService interface {
	CreatePosPaymentMethod(ctx context.Context, req *pb.CreatePosPaymentMethodRequest) (*pb.CreatePosPaymentMethodResponse, error)
	ReadPosPaymentMethod(ctx context.Context, req *pb.ReadPosPaymentMethodRequest) (*pb.ReadPosPaymentMethodResponse, error)
	UpdatePosPaymentMethod(ctx context.Context, req *pb.UpdatePosPaymentMethodRequest) (*pb.UpdatePosPaymentMethodResponse, error)
	DeletePosPaymentMethod(ctx context.Context, req *pb.DeletePosPaymentMethodRequest) (*pb.DeletePosPaymentMethodResponse, error)
	ReadAllPosPaymentMethods(ctx context.Context, req *pb.ReadAllPosPaymentMethodsRequest) (*pb.ReadAllPosPaymentMethodsResponse, error)
}

type posPaymentMethodService struct {
	pb.UnimplementedPosPaymentMethodServiceServer
	repo repository.PosPaymentMethodRepository
}

func NewPosPaymentMethodService(repo repository.PosPaymentMethodRepository) *posPaymentMethodService {
	return &posPaymentMethodService{
		repo: repo,
	}
}

func (s *posPaymentMethodService) CreatePosPaymentMethod(ctx context.Context, req *pb.CreatePosPaymentMethodRequest) (*pb.CreatePosPaymentMethodResponse, error) {
	// Extract role ID from JWT payload
	jwtRoleID := req.JwtPayload.Role
	token := req.JwtToken

	// Get user login role name
	loginRole, err := utils.GetPosRole(jwtRoleID, token)
	if err != nil {
		return nil, err
	}

	if !utils.IsCompanyUser(loginRole.Data.RoleName) {
		return nil, errors.New("users are not allowed to create new payment method")
	}

	req.PosPaymentMethod.PaymentMethodId = uuid.New().String() // Generate a new UUID for the payment_method_id

	now := timestamppb.New(time.Now())
	req.PosPaymentMethod.CreatedAt = now
	req.PosPaymentMethod.UpdatedAt = now

	// Convert pb.PosPaymentMethod to entity.PosPaymentMethod
	gormPaymentMethod := &entity.PosPaymentMethod{
		PaymentMethodID: uuid.MustParse(req.PosPaymentMethod.PaymentMethodId), // auto
		MethodName:      req.PosPaymentMethod.MethodName,
		CompanyID:       uuid.MustParse(req.JwtPayload.CompanyId), // auto
		CreatedAt:       req.PosPaymentMethod.CreatedAt.AsTime(),  // auto
		CreatedBy:       uuid.MustParse(req.JwtPayload.UserId),    // auto
		UpdatedAt:       req.PosPaymentMethod.UpdatedAt.AsTime(),  // auto
		UpdatedBy:       uuid.MustParse(req.JwtPayload.UserId),    // auto
	}

	err = s.repo.CreatePosPaymentMethod(gormPaymentMethod)
	if err != nil {
		return nil, err
	}

	return &pb.CreatePosPaymentMethodResponse{
		PosPaymentMethod: req.PosPaymentMethod,
	}, nil
}

func (s *posPaymentMethodService) ReadPosPaymentMethod(ctx context.Context, req *pb.ReadPosPaymentMethodRequest) (*pb.ReadPosPaymentMethodResponse, error) {
	// Extract role ID from JWT payload
	jwtRoleID := req.JwtPayload.Role
	token := req.JwtToken

	// Get user login role name
	loginRole, err := utils.GetPosRole(jwtRoleID, token)
	if err != nil {
		return nil, err
	}

	if !utils.IsCompanyOrBranchOrStoreUser(loginRole.Data.RoleName) {
		return nil, errors.New("users are not allowed to read payment method")
	}

	posPaymentMethod, err := s.repo.ReadPosPaymentMethod(req.PaymentMethodId)
	if err != nil {
		return nil, err
	}

	companyRole := os.Getenv("COMPANY_USER_ROLE")
	branchRole := os.Getenv("BRANCH_USER_ROLE")
	storeRole := os.Getenv("STORE_USER_ROLE")

	if loginRole.Data.RoleName == companyRole {
		if !utils.VerifyCompanyUserAccess(loginRole.Data.RoleName, posPaymentMethod.CompanyId, req.JwtPayload.CompanyId) {
			return nil, errors.New("company users can only retrieve payment method data within their company")
		}
	}

	if loginRole.Data.RoleName == branchRole {
		if !utils.VerifyBranchUserAccess(loginRole.Data.RoleName, posPaymentMethod.CompanyId, req.JwtPayload.CompanyId) {
			return nil, errors.New("branch users can only retrieve payment method data within their company")
		}
	}

	if loginRole.Data.RoleName == storeRole {
		if !utils.VerifyStoreUserAccess(loginRole.Data.RoleName, posPaymentMethod.CompanyId, req.JwtPayload.CompanyId) {
			return nil, errors.New("store users can only retrieve payment method data within their company")
		}
	}

	return &pb.ReadPosPaymentMethodResponse{
		PosPaymentMethod: posPaymentMethod,
	}, nil
}

func (s *posPaymentMethodService) UpdatePosPaymentMethod(ctx context.Context, req *pb.UpdatePosPaymentMethodRequest) (*pb.UpdatePosPaymentMethodResponse, error) {
	// Extract role ID from JWT payload
	jwtRoleID := req.JwtPayload.Role
	token := req.JwtToken

	// Get user login role name
	loginRole, err := utils.GetPosRole(jwtRoleID, token)
	if err != nil {
		return nil, err
	}

	if !utils.IsCompanyUser(loginRole.Data.RoleName) {
		return nil, errors.New("users are not allowed to update payment method")
	}

	// Get the payment method to be updated
	posPaymentMethod, err := s.repo.ReadPosPaymentMethod(req.PosPaymentMethod.PaymentMethodId)
	if err != nil {
		return nil, err
	}

	companyRole := os.Getenv("COMPANY_USER_ROLE")

	if loginRole.Data.RoleName == companyRole {
		if !utils.VerifyCompanyUserAccess(loginRole.Data.RoleName, posPaymentMethod.CompanyId, req.JwtPayload.CompanyId) {
			return nil, errors.New("company users can only update payment method data within their company")
		}
	}

	now := timestamppb.New(time.Now())
	req.PosPaymentMethod.UpdatedAt = now
	req.PosPaymentMethod.UpdatedBy = req.JwtPayload.UserId

	// Convert pb.PosPaymentMethod to entity.PosPaymentMethod
	gormPaymentMethod := &entity.PosPaymentMethod{
		PaymentMethodID: uuid.MustParse(posPaymentMethod.PaymentMethodId),
		MethodName:      req.PosPaymentMethod.MethodName,
		CompanyID:       uuid.MustParse(posPaymentMethod.CompanyId),
		CreatedAt:       posPaymentMethod.CreatedAt.AsTime(),
		CreatedBy:       uuid.MustParse(posPaymentMethod.CreatedBy),
		UpdatedAt:       req.PosPaymentMethod.UpdatedAt.AsTime(),
		UpdatedBy:       uuid.MustParse(req.PosPaymentMethod.UpdatedBy),
	}

	// Update the payment method
	posPaymentMethod, err = s.repo.UpdatePosPaymentMethod(gormPaymentMethod)
	if err != nil {
		return nil, err
	}

	return &pb.UpdatePosPaymentMethodResponse{
		PosPaymentMethod: posPaymentMethod,
	}, nil
}

func (s *posPaymentMethodService) DeletePosPaymentMethod(ctx context.Context, req *pb.DeletePosPaymentMethodRequest) (*pb.DeletePosPaymentMethodResponse, error) {
	// Extract role ID from JWT payload
	jwtRoleID := req.JwtPayload.Role
	token := req.JwtToken

	// Get user login role name
	loginRole, err := utils.GetPosRole(jwtRoleID, token)
	if err != nil {
		return nil, err
	}

	if !utils.IsCompanyUser(loginRole.Data.RoleName) {
		return nil, errors.New("users are not allowed to delete payment method")
	}

	// Get the payment method to be deleted
	posPaymentMethod, err := s.repo.ReadPosPaymentMethod(req.PaymentMethodId)
	if err != nil {
		return nil, err
	}

	companyRole := os.Getenv("COMPANY_USER_ROLE")

	if loginRole.Data.RoleName == companyRole {
		if !utils.VerifyCompanyUserAccess(loginRole.Data.RoleName, posPaymentMethod.CompanyId, req.JwtPayload.CompanyId) {
			return nil, errors.New("company users can only delete payment method data within their company")
		}
	}

	// Delete the payment method
	err = s.repo.DeletePosPaymentMethod(req.PaymentMethodId)
	if err != nil {
		return nil, err
	}

	return &pb.DeletePosPaymentMethodResponse{
		Success: true,
	}, nil
}

func (s *posPaymentMethodService) ReadAllPosPaymentMethods(ctx context.Context, req *pb.ReadAllPosPaymentMethodsRequest) (*pb.ReadAllPosPaymentMethodsResponse, error) {
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
		return nil, errors.New("users are not allowed to read all payment method")
	}

	paginationResult, err := s.repo.ReadAllPosPaymentMethods(pagination, loginRole.Data.RoleName, req.JwtPayload)
	if err != nil {
		return nil, err
	}

	posPaymentMethods := paginationResult.Records.([]entity.PosPaymentMethod)
	pbPosPaymentMethods := make([]*pb.PosPaymentMethod, len(posPaymentMethods))

	for i, posPaymentMethod := range posPaymentMethods {
		pbPosPaymentMethods[i] = &pb.PosPaymentMethod{
			PaymentMethodId: posPaymentMethod.PaymentMethodID.String(),
			MethodName:      posPaymentMethod.MethodName,
			CompanyId:       posPaymentMethod.CompanyID.String(),
			CreatedAt:       timestamppb.New(posPaymentMethod.CreatedAt),
			CreatedBy:       posPaymentMethod.CreatedBy.String(),
			UpdatedAt:       timestamppb.New(posPaymentMethod.UpdatedAt),
			UpdatedBy:       posPaymentMethod.UpdatedBy.String(),
		}
	}

	return &pb.ReadAllPosPaymentMethodsResponse{
		PosPaymentMethods: pbPosPaymentMethods,
		Limit:             int32(pagination.Limit),
		Page:              int32(pagination.Page),
		MaxPage:           int32(paginationResult.TotalPages),
		Count:             paginationResult.TotalRecords,
	}, nil
}
