package controller

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	pb "github.com/Andrewalifb/alpha-pos-system-sales-service/api/proto"
	"github.com/Andrewalifb/alpha-pos-system-sales-service/dto"
	"github.com/Andrewalifb/alpha-pos-system-sales-service/utils"

	"github.com/gin-gonic/gin"
)

type PosCustomerController interface {
	HandleCreatePosCustomerRequest(c *gin.Context)
	HandleReadPosCustomerRequest(c *gin.Context)
	HandleUpdatePosCustomerRequest(c *gin.Context)
	HandleDeletePosCustomerRequest(c *gin.Context)
	HandleReadAllPosCustomersRequest(c *gin.Context)
}

type posCustomerController struct {
	service pb.PosCustomerServiceClient
}

func NewPosCustomerController(service pb.PosCustomerServiceClient) PosCustomerController {
	return &posCustomerController{
		service: service,
	}
}

func (p *posCustomerController) HandleCreatePosCustomerRequest(ctx *gin.Context) {
	var req pb.CreatePosCustomerRequest

	if err := ctx.ShouldBindJSON(&req.PosCustomer); err != nil {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_CUSTOMER, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	birthOfDate, err := time.Parse("2006-01-02", req.PosCustomer.DateOfBirth)
	if err != nil {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_CUSTOMER, "Invalid birth of date format", nil)
		ctx.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	req.PosCustomer.DateOfBirth = birthOfDate.Format("2006-01-02")

	getJwtPayload, exist := ctx.Get("user")
	if !exist {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_CUSTOMER, "Jwt Payload is Empty", nil)
		ctx.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	req.JwtPayload = getJwtPayload.(*pb.JWTPayload)

	authHeader := ctx.GetHeader("Authorization")
	bearerToken := strings.Split(authHeader, " ")
	if len(bearerToken) != 2 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization token format"})
		ctx.Abort()
		return
	}

	token := bearerToken[1]
	req.JwtToken = token

	res, err := p.service.CreatePosCustomer(ctx, &req)
	if err != nil {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_CUSTOMER, err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	successResponse := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_CREATE_CUSTOMER, res)
	ctx.JSON(http.StatusOK, successResponse)
}

func (p *posCustomerController) HandleReadPosCustomerRequest(ctx *gin.Context) {
	customerID := ctx.Param("id")

	var req pb.ReadPosCustomerRequest
	req.CustomerId = customerID

	getJwtPayload, exist := ctx.Get("user")
	if !exist {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_CUSTOMER, "Jwt Payload is Empty", nil)
		ctx.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	req.JwtPayload = getJwtPayload.(*pb.JWTPayload)

	authHeader := ctx.GetHeader("Authorization")
	bearerToken := strings.Split(authHeader, " ")
	if len(bearerToken) != 2 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization token format"})
		ctx.Abort()
		return
	}

	token := bearerToken[1]
	req.JwtToken = token

	res, err := p.service.ReadPosCustomer(ctx, &req)
	if err != nil {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_CUSTOMER, err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	successResponse := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_CUSTOMER, res)
	ctx.JSON(http.StatusOK, successResponse)
}

func (p *posCustomerController) HandleUpdatePosCustomerRequest(ctx *gin.Context) {
	var req pb.UpdatePosCustomerRequest
	customerID := ctx.Param("id")
	if err := ctx.ShouldBindJSON(&req.PosCustomer); err != nil {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_CUSTOMER, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	getJwtPayload, exist := ctx.Get("user")
	if !exist {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_CUSTOMER, "Jwt Payload is Empty", nil)
		ctx.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	req.JwtPayload = getJwtPayload.(*pb.JWTPayload)

	authHeader := ctx.GetHeader("Authorization")
	bearerToken := strings.Split(authHeader, " ")
	if len(bearerToken) != 2 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization token format"})
		ctx.Abort()
		return
	}

	token := bearerToken[1]
	req.JwtToken = token
	req.PosCustomer.CustomerId = customerID
	res, err := p.service.UpdatePosCustomer(ctx, &req)
	if err != nil {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_CUSTOMER, err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	successResponse := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_UPDATE_CUSTOMER, res)
	ctx.JSON(http.StatusOK, successResponse)
}

func (p *posCustomerController) HandleDeletePosCustomerRequest(ctx *gin.Context) {
	customerID := ctx.Param("id")

	var req pb.DeletePosCustomerRequest
	req.CustomerId = customerID

	getJwtPayload, exist := ctx.Get("user")
	if !exist {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_CUSTOMER, "Jwt Payload is Empty", nil)
		ctx.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	req.JwtPayload = getJwtPayload.(*pb.JWTPayload)

	authHeader := ctx.GetHeader("Authorization")
	bearerToken := strings.Split(authHeader, " ")
	if len(bearerToken) != 2 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization token format"})
		ctx.Abort()
		return
	}

	token := bearerToken[1]
	req.JwtToken = token

	_, err := p.service.DeletePosCustomer(ctx, &req)
	if err != nil {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_CUSTOMER, err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	successResponse := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_DELETE_CUSTOMER, nil)
	ctx.JSON(http.StatusOK, successResponse)
}

func (p *posCustomerController) HandleReadAllPosCustomersRequest(ctx *gin.Context) {
	limitQuery := ctx.Query("limit")
	pageQuery := ctx.Query("page")

	if (limitQuery == "" && pageQuery != "") || (limitQuery != "" && pageQuery == "") {
		errorResponse := utils.BuildResponseFailed("Both limit and page must be provided", "Value Is Empty", pb.CreatePosCustomerResponse{})
		ctx.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	var req pb.ReadAllPosCustomersRequest

	if limitQuery != "" && pageQuery != "" {
		limit, err := strconv.Atoi(limitQuery)
		if err != nil {
			errorResponse := utils.BuildResponseFailed("Invalid limit value", err.Error(), nil)
			ctx.JSON(http.StatusBadRequest, errorResponse)
			return
		}

		page, err := strconv.Atoi(pageQuery)
		if err != nil {
			errorResponse := utils.BuildResponseFailed("Invalid page value", err.Error(), nil)
			ctx.JSON(http.StatusBadRequest, errorResponse)
			return
		}

		getJwtPayload, exist := ctx.Get("user")
		if !exist {
			errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_CUSTOMER, "Jwt Payload is Empty", nil)
			ctx.JSON(http.StatusBadRequest, errorResponse)
			return
		}

		req = pb.ReadAllPosCustomersRequest{
			Limit:      int32(limit),
			Page:       int32(page),
			JwtPayload: getJwtPayload.(*pb.JWTPayload),
		}
	}
	authHeader := ctx.GetHeader("Authorization")
	bearerToken := strings.Split(authHeader, " ")
	if len(bearerToken) != 2 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization token format"})
		ctx.Abort()
		return
	}

	token := bearerToken[1]
	req.JwtToken = token
	res, err := p.service.ReadAllPosCustomers(ctx, &req)
	if err != nil {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_CUSTOMER, err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	successResponse := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_CUSTOMER, res)
	ctx.JSON(http.StatusOK, successResponse)
}
