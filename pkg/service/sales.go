package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	pb "github.com/Andrewalifb/alpha-pos-system-sales-service/api/proto"
	"github.com/Andrewalifb/alpha-pos-system-sales-service/dto"
	"github.com/Andrewalifb/alpha-pos-system-sales-service/entity"
	"github.com/Andrewalifb/alpha-pos-system-sales-service/pkg/repository"
	"github.com/Andrewalifb/alpha-pos-system-sales-service/utils"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PosSaleService interface {
	CreatePosSale(ctx context.Context, req *pb.CreatePosSalesRequest) (*pb.CreatePosSalesResponse, error)
	ReadPosSale(ctx context.Context, req *pb.ReadPosSaleRequest) (*pb.ReadPosSaleResponse, error)
	UpdatePosSale(ctx context.Context, req *pb.UpdatePosSaleRequest) (*pb.UpdatePosSaleResponse, error)
	DeletePosSale(ctx context.Context, req *pb.DeletePosSaleRequest) (*pb.DeletePosSaleResponse, error)
	ReadAllPosSales(ctx context.Context, req *pb.ReadAllPosSalesRequest) (*pb.ReadAllPosSalesResponse, error)
}

type posSaleService struct {
	pb.UnimplementedPosSaleServiceServer
	saleRepo          repository.PosSaleRepository
	invoiceRepo       repository.PosInvoiceRepository
	cashDrawerRepo    repository.PosCashDrawerRepository
	onlinePyamentRepo repository.PosOnlinePaymentRepository
	paymentMethod     repository.PosPaymentMethodRepository
	customer          repository.PosCustomerRepository
	RabbitMQConn      *amqp.Connection
}

func NewPosSaleService(saleRepo repository.PosSaleRepository, invoiceRepo repository.PosInvoiceRepository, cashDrawerRepo repository.PosCashDrawerRepository, onlinePyamentRepo repository.PosOnlinePaymentRepository, paymentMethod repository.PosPaymentMethodRepository, customer repository.PosCustomerRepository, rabbitMQConn *amqp.Connection) *posSaleService {
	return &posSaleService{
		saleRepo:          saleRepo,
		invoiceRepo:       invoiceRepo,
		cashDrawerRepo:    cashDrawerRepo,
		onlinePyamentRepo: onlinePyamentRepo,
		paymentMethod:     paymentMethod,
		customer:          customer,
		RabbitMQConn:      rabbitMQConn,
	}
}

func (s *posSaleService) CreatePosSales(ctx context.Context, req *pb.CreatePosSalesRequest) (*pb.CreatePosSalesResponse, error) {
	// Extract role ID from JWT payload
	jwtRoleID := req.JwtPayload.Role
	token := req.JwtToken

	// Get user login role name
	loginRole, err := utils.GetPosRole(jwtRoleID, req.JwtToken)
	if err != nil {
		return nil, err
	}

	if !utils.IsStoreUser(loginRole.Data.RoleName) {
		return nil, errors.New("users cant create sales transactions")
	}

	var gormSales []*entity.PosSale
	var getTotalDiscount float64
	var totalSalesAfterDiscount float64

	var itemList []dto.Items
	now := timestamppb.New(time.Now())
	timeStamp := now

	for _, posSale := range req.PosSales {
		posSale.SaleId = uuid.New().String() // Generate a new UUID for the sale_id

		// set current time stamp
		posSale.CreatedAt = timeStamp
		posSale.UpdatedAt = timeStamp
		posSale.SaleDate = timeStamp

		// Get product data
		productData, err := utils.GetPosProduct(posSale.ProductId, token)
		if err != nil {
			return nil, err
		}

		// check if product id has promotions
		promotionData, err := utils.GetPosPromotionByProductID(productData.Data.ProductID, req.JwtToken)
		if err != nil {
			fmt.Println("Err :", err)
		}

		// Count total all discount
		if promotionData == nil {
			posSale.Price = productData.Data.Price
		} else {
			if promotionData.DiscountRate != 0.0 && promotionData.Active {
				countDiscountAmount := productData.Data.Price * promotionData.DiscountRate
				getTotalDiscount += countDiscountAmount
				posSale.Price = productData.Data.Price - countDiscountAmount
			} else if promotionData.DiscountRate == 0.0 {
				posSale.Price = productData.Data.Price
			}
		}

		// Count sub total
		posSale.TotalPrice = posSale.Price * float64(posSale.Quantity)

		// Convert pb.PosSale to entity.PosSale
		gormSale := &entity.PosSale{
			SaleID:          uuid.MustParse(posSale.SaleId), // auto
			ReceiptID:       posSale.ReceiptId,              // auto
			ProductID:       uuid.MustParse(productData.Data.ProductID),
			CustomerID:      uuid.MustParse(posSale.CustomerId),
			Quantity:        int(posSale.Quantity),
			Price:           posSale.Price,                          // auto
			SaleDate:        posSale.SaleDate.AsTime(),              // auto
			TotalPrice:      posSale.TotalPrice,                     // auto
			StoreID:         uuid.MustParse(req.JwtPayload.StoreId), // auto
			CashierID:       uuid.MustParse(req.JwtPayload.UserId),  // auto
			PaymentMethodID: uuid.MustParse(posSale.PaymentMethodId),
			BranchID:        uuid.MustParse(req.JwtPayload.BranchId),  // auto
			CompanyID:       uuid.MustParse(req.JwtPayload.CompanyId), // auto
			CreatedAt:       posSale.CreatedAt.AsTime(),               // auto
			CreatedBy:       uuid.MustParse(req.JwtPayload.UserId),    // auto
			UpdatedAt:       posSale.UpdatedAt.AsTime(),               // auto
			UpdatedBy:       uuid.MustParse(req.JwtPayload.UserId),    // auto
		}

		item := dto.Items{
			ProductName: productData.Data.ProductName,
			Quantity:    int(posSale.Quantity),
			Price:       posSale.Price,
			TotalPrice:  posSale.TotalPrice,
		}

		itemList = append(itemList, item)
		gormSales = append(gormSales, gormSale)

		// Create record stock out into inventory history
		inventoryHistory := &dto.PosInventoryHistory{
			ProductId: gormSale.ProductID.String(),
			StoreId:   gormSale.StoreID.String(),
			Quantity:  -int32(gormSale.Quantity),
			BranchId:  gormSale.BranchID.String(),
		}

		_, err = utils.CreatePosInventoryHistory(inventoryHistory, req.JwtPayload, token)
		if err != nil {
			return nil, err
		}
	}

	// Count total sales
	subTotalSales := 0.0
	for _, sale := range gormSales {
		subTotalSales += sale.TotalPrice
	}

	// Decrease total sales with discount
	totalSalesAfterDiscount = subTotalSales - getTotalDiscount

	// insert data into sales database
	createdPosSales, err := s.saleRepo.CreatePosSales(gormSales, token)
	if err != nil {
		return nil, err
	}

	receiptID := createdPosSales[0].ReceiptID
	// Get payment method
	paymentMethodData, err := s.paymentMethod.ReadPosPaymentMethod(gormSales[0].PaymentMethodID.String())
	if err != nil {
		return nil, err
	}

	// If payment method cash --> cash drawer
	drawerId := uuid.New().String()
	invoiceId := uuid.New().String()
	onlinePyamentId := uuid.New().String()

	cashMethod := os.Getenv("CASH_METHOD")
	payLaterMethod := os.Getenv("PAY_LATER_METHOD")

	if paymentMethodData.MethodName == cashMethod {
		cashDrawerData := &entity.PosCashDrawer{
			DrawerID:        uuid.MustParse(drawerId),
			StoreID:         nil,
			EmployeeID:      uuid.MustParse(req.JwtPayload.Role),
			ReceiptID:       receiptID,
			CashIn:          subTotalSales,
			Amount:          totalSalesAfterDiscount,
			CashOut:         0,
			TransactionTime: now.AsTime(),
			RoleID:          uuid.MustParse(req.JwtPayload.Role),
			BranchID:        nil,
			CompanyID:       uuid.MustParse(req.JwtPayload.CompanyId),
			Description:     fmt.Sprintf("Sales Receipt ID %s", gormSales[0].ReceiptID),
			CreatedAt:       now.AsTime(),
			CreatedBy:       uuid.MustParse(req.JwtPayload.Role),
			UpdatedAt:       now.AsTime(),
			UpdatedBy:       uuid.MustParse(req.JwtPayload.Role),
		}
		cashDrawerData.StoreID = utils.ParseUUID(req.JwtPayload.StoreId)
		cashDrawerData.BranchID = utils.ParseUUID(req.JwtPayload.BranchId)
		err := s.cashDrawerRepo.CreatePosCashDrawer(cashDrawerData)
		if err != nil {
			return nil, err
		}
		// if payment method pay later
	} else if paymentMethodData.MethodName == payLaterMethod {
		invoiceData := &entity.PosInvoice{
			InvoiceID: uuid.MustParse(invoiceId),
			ReceiptID: gormSales[0].ReceiptID,
			Date:      now.AsTime(),
			Amount:    totalSalesAfterDiscount,
			Discounts: getTotalDiscount,
			Taxes:     0,
			BranchID:  uuid.MustParse(req.JwtPayload.BranchId),
			CompanyID: uuid.MustParse(req.JwtPayload.CompanyId),
			CreatedAt: now.AsTime(),
			CreatedBy: uuid.MustParse(req.JwtPayload.UserId),
			UpdatedAt: now.AsTime(),
			UpdatedBy: uuid.MustParse(req.JwtPayload.UserId),
		}
		// Create new invoice
		err := s.invoiceRepo.CreatePosInvoice(invoiceData)
		if err != nil {
			return nil, err
		}
	} else {
		// if payment method not cash or pay later
		onlinePaymentData := &entity.PosOnlinePayment{
			PaymentID:     uuid.MustParse(onlinePyamentId),
			StoreID:       uuid.MustParse(req.JwtPayload.StoreId),
			EmployeeID:    uuid.MustParse(req.JwtPayload.UserId),
			PaymentDate:   now.AsTime(),
			ReceiptID:     gormSales[0].ReceiptID,
			Amount:        totalSalesAfterDiscount,
			PaymentMethod: uuid.MustParse(paymentMethodData.PaymentMethodId),
			RoleID:        uuid.MustParse(req.JwtPayload.Role),
			BranchID:      uuid.MustParse(req.JwtPayload.BranchId),
			CompanyID:     uuid.MustParse(req.JwtPayload.CompanyId),
			CreatedAt:     now.AsTime(),
			CreatedBy:     uuid.MustParse(req.JwtPayload.UserId),
			UpdatedAt:     now.AsTime(),
			UpdatedBy:     uuid.MustParse(req.JwtPayload.UserId),
		}
		err := s.onlinePyamentRepo.CreatePosOnlinePayment(onlinePaymentData)
		if err != nil {
			return nil, err
		}
	}

	userData, err := utils.GetPosUser(req.JwtPayload.UserId, token)
	if err != nil {
		return nil, err
	}

	storeData, err := utils.GetPosStore(req.JwtPayload.StoreId, token)
	if err != nil {
		return nil, err
	}

	customer, err := s.customer.ReadPosCustomer(gormSales[0].CustomerID.String())
	if err != nil {
		return nil, err
	}

	ch, err := s.RabbitMQConn.Channel()
	if err != nil {
		return nil, err
	}
	defer ch.Close()

	receipt := dto.DigitalReceipt{
		Receiver: dto.EmailReceiver{
			EmailAddress: customer.Email,
		},
		Header: dto.HeaderReceipt{
			StoreName:           storeData.Data.StoreName,
			StoreAddress:        storeData.Data.Location,
			CashierName:         userData.Data.Username,
			ReceiptID:           receiptID,
			TransactionDateTime: timeStamp.String(),
		},
		Body: dto.BodyReceipt{
			Items: itemList,
		},
		Summary: dto.SummaryReceipt{
			SubTotalAmount: subTotalSales,
			DiscountAmoutn: getTotalDiscount,
			TaxAmount:      0,
			TotalAmount:    totalSalesAfterDiscount,
			CashAmount:     totalSalesAfterDiscount,
			ChangeAmount:   0,
		},
	}

	err = utils.SendDigitalReceipt(receipt, ch, "email_queue")
	if err != nil {
		return nil, err
	}

	return &pb.CreatePosSalesResponse{
		PosSales: req.PosSales,
	}, nil
}

func (s *posSaleService) ReadAllPosSales(ctx context.Context, req *pb.ReadAllPosSalesRequest) (*pb.ReadAllPosSalesResponse, error) {
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
		return nil, errors.New("users cant read all sales transactions")
	}

	paginationResult, err := s.saleRepo.ReadAllPosSales(pagination, loginRole.Data.RoleName, req.JwtPayload)
	if err != nil {
		return nil, err
	}

	posSales := paginationResult.Records.([]entity.PosSale)
	pbPosSales := make([]*pb.PosSale, len(posSales))

	for i, posSale := range posSales {
		pbPosSales[i] = &pb.PosSale{
			SaleId:          posSale.SaleID.String(),
			ReceiptId:       posSale.ReceiptID,
			ProductId:       posSale.ProductID.String(),
			CustomerId:      posSale.CustomerID.String(),
			Quantity:        int32(posSale.Quantity),
			Price:           posSale.Price,
			SaleDate:        timestamppb.New(posSale.SaleDate),
			TotalPrice:      posSale.TotalPrice,
			StoreId:         posSale.StoreID.String(),
			CashierId:       posSale.CashierID.String(),
			PaymentMethodId: posSale.PaymentMethodID.String(),
			BranchId:        posSale.BranchID.String(),
			CompanyId:       posSale.CompanyID.String(),
			CreatedAt:       timestamppb.New(posSale.CreatedAt),
			CreatedBy:       posSale.CreatedBy.String(),
			UpdatedAt:       timestamppb.New(posSale.UpdatedAt),
			UpdatedBy:       posSale.UpdatedBy.String(),
		}
	}

	return &pb.ReadAllPosSalesResponse{
		PosSales: pbPosSales,
		Limit:    int32(pagination.Limit),
		Page:     int32(pagination.Page),
		MaxPage:  int32(paginationResult.TotalPages),
		Count:    paginationResult.TotalRecords,
	}, nil
}

func (s *posSaleService) ReadPosSale(ctx context.Context, req *pb.ReadPosSaleRequest) (*pb.ReadPosSaleResponse, error) {
	// Extract role ID from JWT payload
	jwtRoleID := req.JwtPayload.Role
	token := req.JwtToken

	// Get user login role name
	loginRole, err := utils.GetPosRole(jwtRoleID, token)
	if err != nil {
		return nil, err
	}

	if !utils.IsCompanyOrBranchOrStoreUser(loginRole.Data.RoleName) {
		return nil, errors.New("users cant read sales transactions")
	}

	posSale, err := s.saleRepo.ReadPosSale(req.SaleId)
	if err != nil {
		return nil, err
	}

	companyRole := os.Getenv("COMPANY_USER_ROLE")
	branchRole := os.Getenv("BRANCH_USER_ROLE")
	storeRole := os.Getenv("STORE_USER_ROLE")

	if loginRole.Data.RoleName == companyRole {
		if !utils.VerifyCompanyUserAccess(loginRole.Data.RoleName, posSale.CompanyId, req.JwtPayload.CompanyId) {
			return nil, errors.New("company users can only retrieve sales data within their company")
		}
	}

	if loginRole.Data.RoleName == branchRole {
		if !utils.VerifyBranchUserAccess(loginRole.Data.RoleName, posSale.BranchId, req.JwtPayload.BranchId) {
			return nil, errors.New("branch users can only retrieve sales data within their branch")
		}
	}

	if loginRole.Data.RoleName == storeRole {
		if !utils.VerifyStoreUserAccess(loginRole.Data.RoleName, posSale.StoreId, req.JwtPayload.StoreId) {
			return nil, errors.New("store users can only retrieve sales data within their branch")
		}
	}

	return &pb.ReadPosSaleResponse{
		PosSale: posSale,
	}, nil
}

func (s *posSaleService) UpdatePosSale(ctx context.Context, req *pb.UpdatePosSaleRequest) (*pb.UpdatePosSaleResponse, error) {
	// Extract role ID from JWT payload
	jwtRoleID := req.JwtPayload.Role
	token := req.JwtToken

	// Get user login role name
	loginRole, err := utils.GetPosRole(jwtRoleID, token)
	if err != nil {
		return nil, err
	}

	if !utils.IsBranchUser(loginRole.Data.RoleName) {
		return nil, errors.New("store users cant update sales data")
	}

	// Get the sale to be updated
	posSale, err := s.saleRepo.ReadPosSale(req.PosSale.SaleId)
	if err != nil {
		return nil, err
	}

	branchRole := os.Getenv("BRANCH_USER_ROLE")

	if loginRole.Data.RoleName == branchRole {
		if !utils.VerifyBranchUserAccess(loginRole.Data.RoleName, posSale.BranchId, req.JwtPayload.BranchId) {
			return nil, errors.New("branch users can only update sales data within their branch")
		}
	}

	now := timestamppb.New(time.Now())
	req.PosSale.UpdatedAt = now

	newSaleData := &entity.PosSale{
		SaleID:          uuid.MustParse(posSale.SaleId), // auto
		ReceiptID:       posSale.ReceiptId,              // auto
		ProductID:       uuid.MustParse(req.PosSale.ProductId),
		CustomerID:      uuid.MustParse(req.PosSale.CustomerId),
		Quantity:        int(req.PosSale.Quantity),
		Price:           posSale.Price, // auto
		SaleDate:        req.PosSale.SaleDate.AsTime(),
		TotalPrice:      posSale.Price * float64(req.PosSale.Quantity), // auto
		StoreID:         uuid.MustParse(req.JwtPayload.StoreId),        // auto
		CashierID:       uuid.MustParse(req.JwtPayload.UserId),         // auto
		PaymentMethodID: uuid.MustParse(posSale.PaymentMethodId),
		BranchID:        uuid.MustParse(req.JwtPayload.BranchId),
		CompanyID:       uuid.MustParse(req.JwtPayload.CompanyId),
		CreatedAt:       posSale.CreatedAt.AsTime(),
		CreatedBy:       uuid.MustParse(posSale.CreatedBy),
		UpdatedAt:       req.PosSale.UpdatedAt.AsTime(),
		UpdatedBy:       uuid.MustParse(req.JwtPayload.UserId),
	}

	// Update the sale
	posSale, err = s.saleRepo.UpdatePosSale(newSaleData)
	if err != nil {
		return nil, err
	}

	return &pb.UpdatePosSaleResponse{
		PosSale: posSale,
	}, nil
}

func (s *posSaleService) DeletePosSale(ctx context.Context, req *pb.DeletePosSaleRequest) (*pb.DeletePosSaleResponse, error) {
	// Extract role ID from JWT payload
	jwtRoleID := req.JwtPayload.Role
	token := req.JwtToken

	// Get user login role name
	loginRole, err := utils.GetPosRole(jwtRoleID, token)
	if err != nil {
		return nil, err
	}

	if !utils.IsBranchUser(loginRole.Data.RoleName) {
		return nil, errors.New("store users cant delete sales data")
	}

	// Get the sale to be updated
	posSale, err := s.saleRepo.ReadPosSale(req.SaleId)
	if err != nil {
		return nil, err
	}

	branchRole := os.Getenv("BRANCH_USER_ROLE")

	if loginRole.Data.RoleName == branchRole {
		if !utils.VerifyBranchUserAccess(loginRole.Data.RoleName, posSale.BranchId, req.JwtPayload.BranchId) {
			return nil, errors.New("branch users can only delete sales data within their branch")
		}
	}

	// Delete the sale
	err = s.saleRepo.DeletePosSale(req.SaleId)
	if err != nil {
		return nil, err
	}

	return &pb.DeletePosSaleResponse{
		Success: true,
	}, nil
}
