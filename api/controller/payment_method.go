package controller

import (
	"net/http"
	"strconv"
	"strings"

	pb "github.com/Andrewalifb/alpha-pos-system-sales-service/api/proto"
	"github.com/Andrewalifb/alpha-pos-system-sales-service/dto"
	"github.com/Andrewalifb/alpha-pos-system-sales-service/utils"

	"github.com/gin-gonic/gin"
)

type PosPaymentMethodController interface {
	HandleCreatePosPaymentMethodRequest(c *gin.Context)
	HandleReadPosPaymentMethodRequest(c *gin.Context)
	HandleUpdatePosPaymentMethodRequest(c *gin.Context)
	HandleDeletePosPaymentMethodRequest(c *gin.Context)
	HandleReadAllPosPaymentMethodsRequest(c *gin.Context)
}

type posPaymentMethodController struct {
	service pb.PosPaymentMethodServiceClient
}

func NewPosPaymentMethodController(service pb.PosPaymentMethodServiceClient) PosPaymentMethodController {
	return &posPaymentMethodController{
		service: service,
	}
}

func (c *posPaymentMethodController) HandleCreatePosPaymentMethodRequest(ctx *gin.Context) {
	// Declare req body Pos Payment Method
	var req pb.CreatePosPaymentMethodRequest

	// First, binding payment method data
	if err := ctx.ShouldBindJSON(&req.PosPaymentMethod); err != nil {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_PAYMENT_METHOD, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	// Get JWT Payload data from middleware
	getJwtPayload, exist := ctx.Get("user")
	if !exist {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_PAYMENT_METHOD, "Jwt Payload is Empty", nil)
		ctx.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	// Add JWT payload from middleware into req body
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

	// Service call
	res, err := c.service.CreatePosPaymentMethod(ctx, &req)
	if err != nil {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_PAYMENT_METHOD, err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	// Success response
	successResponse := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_CREATE_PAYMENT_METHOD, res)
	ctx.JSON(http.StatusOK, successResponse)
}

func (c *posPaymentMethodController) HandleReadPosPaymentMethodRequest(ctx *gin.Context) {
	var req pb.ReadPosPaymentMethodRequest

	// Get payment method ID from URL
	paymentMethodID := ctx.Param("id")
	req.PaymentMethodId = paymentMethodID

	// Get JWT Payload data from middleware
	getJwtPayload, exist := ctx.Get("user")
	if !exist {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_PAYMENT_METHOD, "Jwt Payload is Empty", nil)
		ctx.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	// Add JWT payload from middleware into req body
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

	res, err := c.service.ReadPosPaymentMethod(ctx, &req)
	if err != nil {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_PAYMENT_METHOD, err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	successResponse := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_PAYMENT_METHOD, res)
	ctx.JSON(http.StatusOK, successResponse)
}

func (c *posPaymentMethodController) HandleUpdatePosPaymentMethodRequest(ctx *gin.Context) {
	var req pb.UpdatePosPaymentMethodRequest
	paymentMethodID := ctx.Param("id")
	if err := ctx.ShouldBindJSON(&req.PosPaymentMethod); err != nil {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_PAYMENT_METHOD, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	// Get JWT Payload data from middleware
	getJwtPayload, exist := ctx.Get("user")
	if !exist {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_PAYMENT_METHOD, "Jwt Payload is Empty", nil)
		ctx.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	// Add JWT payload from middleware into req body
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
	req.PosPaymentMethod.PaymentMethodId = paymentMethodID
	res, err := c.service.UpdatePosPaymentMethod(ctx, &req)
	if err != nil {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_PAYMENT_METHOD, err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	successResponse := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_UPDATE_PAYMENT_METHOD, res)
	ctx.JSON(http.StatusOK, successResponse)
}

func (c *posPaymentMethodController) HandleDeletePosPaymentMethodRequest(ctx *gin.Context) {
	var req pb.DeletePosPaymentMethodRequest

	// Get payment method ID from URL
	paymentMethodID := ctx.Param("id")
	req.PaymentMethodId = paymentMethodID

	// Get JWT Payload data from middleware
	getJwtPayload, exist := ctx.Get("user")
	if !exist {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_PAYMENT_METHOD, "Jwt Payload is Empty", nil)
		ctx.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	// Add JWT payload from middleware into req body
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

	res, err := c.service.DeletePosPaymentMethod(ctx, &req)
	if err != nil {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_PAYMENT_METHOD, err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	successResponse := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_DELETE_PAYMENT_METHOD, res)
	ctx.JSON(http.StatusOK, successResponse)
}

func (c *posPaymentMethodController) HandleReadAllPosPaymentMethodsRequest(ctx *gin.Context) {
	limitQuery := ctx.Query("limit")
	pageQuery := ctx.Query("page")

	if (limitQuery == "" && pageQuery != "") || (limitQuery != "" && pageQuery == "") {
		errorResponse := utils.BuildResponseFailed("Both limit and page must be provided", "Value Is Empty", pb.CreatePosPaymentMethodResponse{})
		ctx.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	var req pb.ReadAllPosPaymentMethodsRequest

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

		// Get JWT Payload data from middleware
		getJwtPayload, exist := ctx.Get("user")
		if !exist {
			errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_PAYMENT_METHOD, "Jwt Payload is Empty", nil)
			ctx.JSON(http.StatusBadRequest, errorResponse)
			return
		}

		req = pb.ReadAllPosPaymentMethodsRequest{
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

	res, err := c.service.ReadAllPosPaymentMethods(ctx, &req)
	if err != nil {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_PAYMENT_METHOD, err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	successResponse := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_PAYMENT_METHOD, res)
	ctx.JSON(http.StatusOK, successResponse)
}
