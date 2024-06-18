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

type PosCustomerService interface {
	CreatePosCustomer(ctx context.Context, req *pb.CreatePosCustomerRequest) (*pb.CreatePosCustomerResponse, error)
	ReadPosCustomer(ctx context.Context, req *pb.ReadPosCustomerRequest) (*pb.ReadPosCustomerResponse, error)
	UpdatePosCustomer(ctx context.Context, req *pb.UpdatePosCustomerRequest) (*pb.UpdatePosCustomerResponse, error)
	DeletePosCustomer(ctx context.Context, req *pb.DeletePosCustomerRequest) (*pb.DeletePosCustomerResponse, error)
	ReadAllPosCustomers(ctx context.Context, req *pb.ReadAllPosCustomersRequest) (*pb.ReadAllPosCustomersResponse, error)
}

type posCustomerService struct {
	pb.UnimplementedPosCustomerServiceServer
	customerRepo repository.PosCustomerRepository
}

func NewPosCustomerService(customerRepo repository.PosCustomerRepository) *posCustomerService {
	return &posCustomerService{
		customerRepo: customerRepo,
	}
}

func (s *posCustomerService) CreatePosCustomer(ctx context.Context, req *pb.CreatePosCustomerRequest) (*pb.CreatePosCustomerResponse, error) {
	// Extract role ID from JWT payload
	jwtRoleID := req.JwtPayload.Role
	token := req.JwtToken

	// Get user login role name
	loginRole, err := utils.GetPosRole(jwtRoleID, token)
	if err != nil {
		return nil, err
	}

	if !utils.IsBranchOrStoreUser(loginRole.Data.RoleName) {
		return nil, errors.New("users cant create customer data")
	}

	req.PosCustomer.CustomerId = uuid.New().String() // Generate a new UUID for the customer_id

	dateOfBirth, err := time.Parse("2006-01-02", req.PosCustomer.DateOfBirth)
	if err != nil {
		return nil, err
	}

	now := timestamppb.New(time.Now())
	req.PosCustomer.CreatedAt = now
	req.PosCustomer.UpdatedAt = now
	req.PosCustomer.RegistrationDate = now

	// Convert pb.PosCustomer to entity.PosCustomer
	gormCustomer := &entity.PosCustomer{
		CustomerID:       uuid.MustParse(req.PosCustomer.CustomerId), // auto
		FirstName:        req.PosCustomer.FirstName,
		LastName:         req.PosCustomer.LastName,
		Email:            req.PosCustomer.Email,
		PhoneNumber:      req.PosCustomer.PhoneNumber,
		DateOfBirth:      dateOfBirth,
		RegistrationDate: req.PosCustomer.RegistrationDate.AsTime(), // auto
		Address:          req.PosCustomer.Address,
		City:             req.PosCustomer.City,
		Country:          req.PosCustomer.Country,
		BranchID:         uuid.MustParse(req.JwtPayload.BranchId),  // auto
		CompanyID:        uuid.MustParse(req.JwtPayload.CompanyId), // auto
		CreatedAt:        req.PosCustomer.CreatedAt.AsTime(),       //  auto
		CreatedBy:        uuid.MustParse(req.JwtPayload.UserId),    // auto
		UpdatedAt:        req.PosCustomer.UpdatedAt.AsTime(),       // auto
		UpdatedBy:        uuid.MustParse(req.JwtPayload.UserId),    // auto
	}

	err = s.customerRepo.CreatePosCustomer(gormCustomer)
	if err != nil {
		return nil, err
	}

	return &pb.CreatePosCustomerResponse{
		PosCustomer: req.PosCustomer,
	}, nil
}

func (s *posCustomerService) ReadAllPosCustomers(ctx context.Context, req *pb.ReadAllPosCustomersRequest) (*pb.ReadAllPosCustomersResponse, error) {
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
		return nil, errors.New("users cant read all customer data")
	}

	paginationResult, err := s.customerRepo.ReadAllPosCustomers(pagination, loginRole.Data.RoleName, req.JwtPayload)
	if err != nil {
		return nil, err
	}

	posCustomers := paginationResult.Records.([]entity.PosCustomer)
	pbPosCustomers := make([]*pb.PosCustomer, len(posCustomers))

	for i, posCustomer := range posCustomers {
		pbPosCustomers[i] = &pb.PosCustomer{
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
	}

	return &pb.ReadAllPosCustomersResponse{
		PosCustomers: pbPosCustomers,
		Limit:        int32(pagination.Limit),
		Page:         int32(pagination.Page),
		MaxPage:      int32(paginationResult.TotalPages),
		Count:        paginationResult.TotalRecords,
	}, nil
}

func (s *posCustomerService) ReadPosCustomer(ctx context.Context, req *pb.ReadPosCustomerRequest) (*pb.ReadPosCustomerResponse, error) {
	// Extract role ID from JWT payload
	jwtRoleID := req.JwtPayload.Role
	token := req.JwtToken

	// Get user login role name
	loginRole, err := utils.GetPosRole(jwtRoleID, token)
	if err != nil {
		return nil, err
	}

	if !utils.IsCompanyOrBranchOrStoreUser(loginRole.Data.RoleName) {
		return nil, errors.New("users cant read customer data")
	}

	posCustomer, err := s.customerRepo.ReadPosCustomer(req.CustomerId)
	if err != nil {
		return nil, err
	}

	companyRole := os.Getenv("COMPANY_USER_ROLE")
	branchRole := os.Getenv("BRANCH_USER_ROLE")
	storeRole := os.Getenv("STORE_USER_ROLE")

	if loginRole.Data.RoleName == companyRole {
		if !utils.VerifyCompanyUserAccess(loginRole.Data.RoleName, posCustomer.CompanyId, req.JwtPayload.CompanyId) {
			return nil, errors.New("company users can only retrieve customer data within their company")
		}
	}

	if loginRole.Data.RoleName == branchRole {
		if !utils.VerifyBranchUserAccess(loginRole.Data.RoleName, posCustomer.BranchId, req.JwtPayload.BranchId) {
			return nil, errors.New("branch users can only retrieve customer data within their branch")
		}
	}

	if loginRole.Data.RoleName == storeRole {
		if !utils.VerifyStoreUserAccess(loginRole.Data.RoleName, posCustomer.BranchId, req.JwtPayload.BranchId) {
			return nil, errors.New("store users can only retrieve customer data within their branch")
		}
	}

	return &pb.ReadPosCustomerResponse{
		PosCustomer: posCustomer,
	}, nil
}

func (s *posCustomerService) UpdatePosCustomer(ctx context.Context, req *pb.UpdatePosCustomerRequest) (*pb.UpdatePosCustomerResponse, error) {
	// Get the role name from the role ID in the JWT payload
	// Extract role ID from JWT payload
	jwtRoleID := req.JwtPayload.Role
	token := req.JwtToken

	// Get user login role name
	loginRole, err := utils.GetPosRole(jwtRoleID, token)
	if err != nil {
		return nil, err
	}

	if !utils.IsBranchOrStoreUser(loginRole.Data.RoleName) {
		return nil, errors.New("users cant update customer data")
	}

	// Get the customer to be updated
	posCustomer, err := s.customerRepo.ReadPosCustomer(req.PosCustomer.CustomerId)
	if err != nil {
		return nil, err
	}

	branchRole := os.Getenv("BRANCH_USER_ROLE")
	storeRole := os.Getenv("STORE_USER_ROLE")

	if loginRole.Data.RoleName == branchRole {
		if !utils.VerifyBranchUserAccess(loginRole.Data.RoleName, posCustomer.BranchId, req.JwtPayload.BranchId) {
			return nil, errors.New("branch users can only update customer data within their branch")
		}
	}

	if loginRole.Data.RoleName == storeRole {
		if !utils.VerifyStoreUserAccess(loginRole.Data.RoleName, posCustomer.BranchId, req.JwtPayload.BranchId) {
			return nil, errors.New("store users can only update customer data within their branch")
		}
	}

	dateOfBirth, err := time.Parse("2006-01-02", req.PosCustomer.DateOfBirth)
	if err != nil {
		return nil, err
	}

	now := timestamppb.New(time.Now())
	req.PosCustomer.UpdatedAt = now

	newCustomerData := &entity.PosCustomer{
		CustomerID:       uuid.MustParse(posCustomer.CustomerId),
		FirstName:        req.PosCustomer.FirstName,
		LastName:         req.PosCustomer.LastName,
		Email:            req.PosCustomer.Email,
		PhoneNumber:      req.PosCustomer.PhoneNumber,
		DateOfBirth:      dateOfBirth,
		RegistrationDate: posCustomer.RegistrationDate.AsTime(),
		Address:          req.PosCustomer.Address,
		City:             req.PosCustomer.City,
		Country:          req.PosCustomer.Country,
		BranchID:         uuid.MustParse(posCustomer.BranchId),
		CompanyID:        uuid.MustParse(posCustomer.CompanyId),
		CreatedAt:        posCustomer.CreatedAt.AsTime(),
		CreatedBy:        uuid.MustParse(posCustomer.CreatedBy),
		UpdatedAt:        req.PosCustomer.UpdatedAt.AsTime(),
		UpdatedBy:        uuid.MustParse(req.JwtPayload.UserId),
	}

	// Update the customer
	posCustomer, err = s.customerRepo.UpdatePosCustomer(newCustomerData)
	if err != nil {
		return nil, err
	}

	return &pb.UpdatePosCustomerResponse{
		PosCustomer: posCustomer,
	}, nil
}
func (s *posCustomerService) DeletePosCustomer(ctx context.Context, req *pb.DeletePosCustomerRequest) (*pb.DeletePosCustomerResponse, error) {
	// Extract role ID from JWT payload
	jwtRoleID := req.JwtPayload.Role
	token := req.JwtToken

	// Get user login role name
	loginRole, err := utils.GetPosRole(jwtRoleID, token)
	if err != nil {
		return nil, err
	}

	if !utils.IsCompanyOrBranchUser(loginRole.Data.RoleName) {
		return nil, errors.New("users cant delete sales data")
	}

	// Get the customer to be updated
	posCustomer, err := s.customerRepo.ReadPosCustomer(req.CustomerId)
	if err != nil {
		return nil, err
	}

	companyRole := os.Getenv("COMPANY_USER_ROLE")
	branchRole := os.Getenv("BRANCH_USER_ROLE")

	if loginRole.Data.RoleName == companyRole {
		if !utils.VerifyCompanyUserAccess(loginRole.Data.RoleName, posCustomer.CompanyId, req.JwtPayload.CompanyId) {
			return nil, errors.New("company users can only delete customer data within their company")
		}
	}

	if loginRole.Data.RoleName == branchRole {
		if !utils.VerifyBranchUserAccess(loginRole.Data.RoleName, posCustomer.BranchId, req.JwtPayload.BranchId) {
			return nil, errors.New("branch users can only delete customer data within their branch")
		}
	}

	// Delete the customer
	err = s.customerRepo.DeletePosCustomer(req.CustomerId)
	if err != nil {
		return nil, err
	}

	return &pb.DeletePosCustomerResponse{
		Success: true,
	}, nil
}
