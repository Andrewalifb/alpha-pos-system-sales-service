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

type PosOnlinePaymentController interface {
	HandleCreatePosOnlinePaymentRequest(c *gin.Context)
	HandleReadPosOnlinePaymentRequest(c *gin.Context)
	HandleUpdatePosOnlinePaymentRequest(c *gin.Context)
	HandleDeletePosOnlinePaymentRequest(c *gin.Context)
	HandleReadAllPosOnlinePaymentsRequest(c *gin.Context)
}

type posOnlinePaymentController struct {
	service pb.PosOnlinePaymentServiceClient
}

func NewPosOnlinePaymentController(service pb.PosOnlinePaymentServiceClient) PosOnlinePaymentController {
	return &posOnlinePaymentController{
		service: service,
	}
}

func (p *posOnlinePaymentController) HandleCreatePosOnlinePaymentRequest(ctx *gin.Context) {
	var req pb.CreatePosOnlinePaymentRequest

	if err := ctx.ShouldBindJSON(&req.PosOnlinePayment); err != nil {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_ONLINE_PAYMENT, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	getJwtPayload, exist := ctx.Get("user")
	if !exist {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_ONLINE_PAYMENT, "Jwt Payload is Empty", nil)
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

	res, err := p.service.CreatePosOnlinePayment(ctx, &req)
	if err != nil {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_ONLINE_PAYMENT, err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	successResponse := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_CREATE_ONLINE_PAYMENT, res)
	ctx.JSON(http.StatusOK, successResponse)
}

func (p *posOnlinePaymentController) HandleReadPosOnlinePaymentRequest(ctx *gin.Context) {
	paymentID := ctx.Param("id")

	var req pb.ReadPosOnlinePaymentRequest
	req.PaymentId = paymentID

	getJwtPayload, exist := ctx.Get("user")
	if !exist {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_ONLINE_PAYMENT, "Jwt Payload is Empty", nil)
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

	res, err := p.service.ReadPosOnlinePayment(ctx, &req)
	if err != nil {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_ONLINE_PAYMENT, err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	successResponse := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_ONLINE_PAYMENT, res)
	ctx.JSON(http.StatusOK, successResponse)
}

func (p *posOnlinePaymentController) HandleUpdatePosOnlinePaymentRequest(ctx *gin.Context) {
	var req pb.UpdatePosOnlinePaymentRequest
	paymentID := ctx.Param("id")
	if err := ctx.ShouldBindJSON(&req.PosOnlinePayment); err != nil {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_ONLINE_PAYMENT, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	getJwtPayload, exist := ctx.Get("user")
	if !exist {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_ONLINE_PAYMENT, "Jwt Payload is Empty", nil)
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
	req.PosOnlinePayment.PaymentId = paymentID
	res, err := p.service.UpdatePosOnlinePayment(ctx, &req)
	if err != nil {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_ONLINE_PAYMENT, err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	successResponse := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_UPDATE_ONLINE_PAYMENT, res)
	ctx.JSON(http.StatusOK, successResponse)
}

func (p *posOnlinePaymentController) HandleDeletePosOnlinePaymentRequest(ctx *gin.Context) {
	paymentID := ctx.Param("id")

	var req pb.DeletePosOnlinePaymentRequest
	req.PaymentId = paymentID

	getJwtPayload, exist := ctx.Get("user")
	if !exist {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_ONLINE_PAYMENT, "Jwt Payload is Empty", nil)
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

	_, err := p.service.DeletePosOnlinePayment(ctx, &req)
	if err != nil {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_ONLINE_PAYMENT, err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	successResponse := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_DELETE_ONLINE_PAYMENT, nil)
	ctx.JSON(http.StatusOK, successResponse)
}

func (p *posOnlinePaymentController) HandleReadAllPosOnlinePaymentsRequest(ctx *gin.Context) {
	limitQuery := ctx.Query("limit")
	pageQuery := ctx.Query("page")

	if (limitQuery == "" && pageQuery != "") || (limitQuery != "" && pageQuery == "") {
		errorResponse := utils.BuildResponseFailed("Both limit and page must be provided", "Value Is Empty", pb.CreatePosOnlinePaymentResponse{})
		ctx.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	var req pb.ReadAllPosOnlinePaymentsRequest

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
			errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_ONLINE_PAYMENT, "Jwt Payload is Empty", nil)
			ctx.JSON(http.StatusBadRequest, errorResponse)
			return
		}

		req = pb.ReadAllPosOnlinePaymentsRequest{
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
	res, err := p.service.ReadAllPosOnlinePayments(ctx, &req)
	if err != nil {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_ONLINE_PAYMENT, err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	successResponse := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_ONLINE_PAYMENT, res)
	ctx.JSON(http.StatusOK, successResponse)
}
