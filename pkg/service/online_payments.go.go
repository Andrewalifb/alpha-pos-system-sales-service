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

type PosOnlinePaymentService interface {
	CreatePosOnlinePayment(ctx context.Context, req *pb.CreatePosOnlinePaymentRequest) (*pb.CreatePosOnlinePaymentResponse, error)
	ReadPosOnlinePayment(ctx context.Context, req *pb.ReadPosOnlinePaymentRequest) (*pb.ReadPosOnlinePaymentResponse, error)
	UpdatePosOnlinePayment(ctx context.Context, req *pb.UpdatePosOnlinePaymentRequest) (*pb.UpdatePosOnlinePaymentResponse, error)
	DeletePosOnlinePayment(ctx context.Context, req *pb.DeletePosOnlinePaymentRequest) (*pb.DeletePosOnlinePaymentResponse, error)
	ReadAllPosOnlinePayments(ctx context.Context, req *pb.ReadAllPosOnlinePaymentsRequest) (*pb.ReadAllPosOnlinePaymentsResponse, error)
}

type posOnlinePaymentService struct {
	pb.UnimplementedPosOnlinePaymentServiceServer
	onlinePaymentRepo  repository.PosOnlinePaymentRepository
	CompanyServiceConn *grpc.ClientConn
}

func NewPosOnlinePaymentService(onlinePaymentRepo repository.PosOnlinePaymentRepository, companyServiceConn *grpc.ClientConn) *posOnlinePaymentService {
	return &posOnlinePaymentService{
		onlinePaymentRepo:  onlinePaymentRepo,
		CompanyServiceConn: companyServiceConn,
	}
}

func (s *posOnlinePaymentService) CreatePosOnlinePayment(ctx context.Context, req *pb.CreatePosOnlinePaymentRequest) (*pb.CreatePosOnlinePaymentResponse, error) {
	// Extract role ID from JWT payload
	jwtRoleID := req.JwtPayload.Role
	// token := req.JwtToken

	// Get user login role name
	loginRole, err := utils.GetPosRoleById(s.CompanyServiceConn, jwtRoleID, req.JwtPayload)
	if err != nil {
		return nil, err
	}

	if !utils.IsStoreUser(loginRole.PosRole.RoleName) {
		return nil, errors.New("users cant create new online payment")
	}

	req.PosOnlinePayment.PaymentId = uuid.New().String() // Generate a new UUID for the payment_id

	now := timestamppb.New(time.Now())
	req.PosOnlinePayment.CreatedAt = now
	req.PosOnlinePayment.UpdatedAt = now

	// Convert pb.PosOnlinePayment to entity.PosOnlinePayment
	gormOnlinePayment := &entity.PosOnlinePayment{
		PaymentID:     uuid.MustParse(req.PosOnlinePayment.PaymentId),
		StoreID:       uuid.MustParse(req.PosOnlinePayment.StoreId),
		EmployeeID:    uuid.MustParse(req.PosOnlinePayment.EmployeeId),
		PaymentDate:   req.PosOnlinePayment.PaymentDate.AsTime(),
		ReceiptID:     req.PosOnlinePayment.ReceiptId,
		Amount:        req.PosOnlinePayment.Amount,
		PaymentMethod: uuid.MustParse(req.PosOnlinePayment.PaymentMethod),
		RoleID:        uuid.MustParse(req.PosOnlinePayment.RoleId),
		BranchID:      uuid.MustParse(req.PosOnlinePayment.BranchId),
		CompanyID:     uuid.MustParse(req.PosOnlinePayment.CompanyId),
		CreatedAt:     req.PosOnlinePayment.CreatedAt.AsTime(),
		CreatedBy:     uuid.MustParse(req.PosOnlinePayment.CreatedBy),
		UpdatedAt:     req.PosOnlinePayment.UpdatedAt.AsTime(),
		UpdatedBy:     uuid.MustParse(req.PosOnlinePayment.UpdatedBy),
	}

	err = s.onlinePaymentRepo.CreatePosOnlinePayment(gormOnlinePayment)
	if err != nil {
		return nil, err
	}

	return &pb.CreatePosOnlinePaymentResponse{
		PosOnlinePayment: req.PosOnlinePayment,
	}, nil
}

func (s *posOnlinePaymentService) ReadPosOnlinePayment(ctx context.Context, req *pb.ReadPosOnlinePaymentRequest) (*pb.ReadPosOnlinePaymentResponse, error) {
	// Extract role ID from JWT payload
	jwtRoleID := req.JwtPayload.Role
	// token := req.JwtToken

	// Get user login role name
	loginRole, err := utils.GetPosRoleById(s.CompanyServiceConn, jwtRoleID, req.JwtPayload)
	if err != nil {
		return nil, err
	}

	if !utils.IsCompanyOrBranchOrStoreUser(loginRole.PosRole.RoleName) {
		return nil, errors.New("users cant read online payment")
	}

	posOnlinePayment, err := s.onlinePaymentRepo.ReadPosOnlinePayment(req.PaymentId)
	if err != nil {
		return nil, err
	}

	companyRole := os.Getenv("COMPANY_USER_ROLE")
	branchRole := os.Getenv("BRANCH_USER_ROLE")
	storeRole := os.Getenv("STORE_USER_ROLE")

	if loginRole.PosRole.RoleName == companyRole {
		if !utils.VerifyCompanyUserAccess(loginRole.PosRole.RoleName, posOnlinePayment.CompanyId, req.JwtPayload.CompanyId) {
			return nil, errors.New("company users can only retrieve online payment data within their company")
		}
	}

	if loginRole.PosRole.RoleName == branchRole {
		if !utils.VerifyBranchUserAccess(loginRole.PosRole.RoleName, posOnlinePayment.BranchId, req.JwtPayload.BranchId) {
			return nil, errors.New("branch users can only retrieve online payment data within their branch")
		}
	}

	if loginRole.PosRole.RoleName == storeRole {
		if !utils.VerifyStoreUserAccess(loginRole.PosRole.RoleName, posOnlinePayment.StoreId, req.JwtPayload.StoreId) {
			return nil, errors.New("store users can only retrieve online payment data within their branch")
		}
	}

	return &pb.ReadPosOnlinePaymentResponse{
		PosOnlinePayment: posOnlinePayment,
	}, nil
}

func (s *posOnlinePaymentService) UpdatePosOnlinePayment(ctx context.Context, req *pb.UpdatePosOnlinePaymentRequest) (*pb.UpdatePosOnlinePaymentResponse, error) {
	// Get the role name from the role ID in the JWT payload
	// Extract role ID from JWT payload
	jwtRoleID := req.JwtPayload.Role
	// token := req.JwtToken

	// Get user login role name
	loginRole, err := utils.GetPosRoleById(s.CompanyServiceConn, jwtRoleID, req.JwtPayload)
	if err != nil {
		return nil, err
	}

	if !utils.IsBranchUser(loginRole.PosRole.RoleName) {
		return nil, errors.New("users cant update online payment data")
	}

	// Get the online payment to be updated
	posOnlinePayment, err := s.onlinePaymentRepo.ReadPosOnlinePayment(req.PosOnlinePayment.PaymentId)
	if err != nil {
		return nil, err
	}

	branchRole := os.Getenv("BRANCH_USER_ROLE")

	if loginRole.PosRole.RoleName == branchRole {
		if !utils.VerifyBranchUserAccess(loginRole.PosRole.RoleName, posOnlinePayment.BranchId, req.JwtPayload.BranchId) {
			return nil, errors.New("branch users can only update online payment data within their branch")
		}
	}

	now := timestamppb.New(time.Now())
	req.PosOnlinePayment.UpdatedAt = now

	newOnlinePaymentData := &entity.PosOnlinePayment{
		PaymentID:     uuid.MustParse(req.PosOnlinePayment.PaymentId),
		StoreID:       uuid.MustParse(req.PosOnlinePayment.StoreId),
		EmployeeID:    uuid.MustParse(req.PosOnlinePayment.EmployeeId),
		PaymentDate:   req.PosOnlinePayment.PaymentDate.AsTime(),
		ReceiptID:     req.PosOnlinePayment.ReceiptId,
		Amount:        req.PosOnlinePayment.Amount,
		PaymentMethod: uuid.MustParse(req.PosOnlinePayment.PaymentMethod),
		RoleID:        uuid.MustParse(req.PosOnlinePayment.RoleId),
		BranchID:      uuid.MustParse(req.PosOnlinePayment.BranchId),
		CompanyID:     uuid.MustParse(req.PosOnlinePayment.CompanyId),
		CreatedAt:     req.PosOnlinePayment.CreatedAt.AsTime(),
		CreatedBy:     uuid.MustParse(req.PosOnlinePayment.CreatedBy),
		UpdatedAt:     req.PosOnlinePayment.UpdatedAt.AsTime(),
		UpdatedBy:     uuid.MustParse(req.PosOnlinePayment.UpdatedBy),
	}

	// Update the online payment
	posOnlinePayment, err = s.onlinePaymentRepo.UpdatePosOnlinePayment(newOnlinePaymentData)
	if err != nil {
		return nil, err
	}

	return &pb.UpdatePosOnlinePaymentResponse{
		PosOnlinePayment: posOnlinePayment,
	}, nil
}

func (s *posOnlinePaymentService) DeletePosOnlinePayment(ctx context.Context, req *pb.DeletePosOnlinePaymentRequest) (*pb.DeletePosOnlinePaymentResponse, error) {
	// Extract role ID from JWT payload
	jwtRoleID := req.JwtPayload.Role
	// token := req.JwtToken

	// Get user login role name
	loginRole, err := utils.GetPosRoleById(s.CompanyServiceConn, jwtRoleID, req.JwtPayload)
	if err != nil {
		return nil, err
	}

	if !utils.IsBranchUser(loginRole.PosRole.RoleName) {
		return nil, errors.New("users cant delete online payment data")
	}

	// Get the online payment to be updated
	posOnlinePayment, err := s.onlinePaymentRepo.ReadPosOnlinePayment(req.PaymentId)
	if err != nil {
		return nil, err
	}

	branchRole := os.Getenv("BRANCH_USER_ROLE")

	if loginRole.PosRole.RoleName == branchRole {
		if !utils.VerifyBranchUserAccess(loginRole.PosRole.RoleName, posOnlinePayment.BranchId, req.JwtPayload.BranchId) {
			return nil, errors.New("branch users can only delete online payment data within their branch")
		}
	}

	// Delete the online payment
	err = s.onlinePaymentRepo.DeletePosOnlinePayment(req.PaymentId)
	if err != nil {
		return nil, err
	}

	return &pb.DeletePosOnlinePaymentResponse{
		Success: true,
	}, nil
}

func (s *posOnlinePaymentService) ReadAllPosOnlinePayments(ctx context.Context, req *pb.ReadAllPosOnlinePaymentsRequest) (*pb.ReadAllPosOnlinePaymentsResponse, error) {
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
		return nil, errors.New("users cant read all online payment")
	}

	paginationResult, err := s.onlinePaymentRepo.ReadAllPosOnlinePayments(pagination, loginRole.PosRole.RoleName, req.JwtPayload)
	if err != nil {
		return nil, err
	}

	posOnlinePayments := paginationResult.Records.([]entity.PosOnlinePayment)
	pbPosOnlinePayments := make([]*pb.PosOnlinePayment, len(posOnlinePayments))

	for i, posOnlinePayment := range posOnlinePayments {
		pbPosOnlinePayments[i] = &pb.PosOnlinePayment{
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

	}

	return &pb.ReadAllPosOnlinePaymentsResponse{
		PosOnlinePayments: pbPosOnlinePayments,
		Limit:             int32(pagination.Limit),
		Page:              int32(pagination.Page),
		MaxPage:           int32(paginationResult.TotalPages),
		Count:             paginationResult.TotalRecords,
	}, nil
}
