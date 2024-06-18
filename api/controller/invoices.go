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

type PosInvoiceController interface {
	HandleCreatePosInvoiceRequest(c *gin.Context)
	HandleReadPosInvoiceRequest(c *gin.Context)
	HandleUpdatePosInvoiceRequest(c *gin.Context)
	HandleDeletePosInvoiceRequest(c *gin.Context)
	HandleReadAllPosInvoicesRequest(c *gin.Context)
}

type posInvoiceController struct {
	service pb.PosInvoiceServiceClient
}

func NewPosInvoiceController(service pb.PosInvoiceServiceClient) PosInvoiceController {
	return &posInvoiceController{
		service: service,
	}
}

func (p *posInvoiceController) HandleCreatePosInvoiceRequest(ctx *gin.Context) {
	var req pb.CreatePosInvoiceRequest

	if err := ctx.ShouldBindJSON(&req.PosInvoice); err != nil {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_INVOICES, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	getJwtPayload, exist := ctx.Get("user")
	if !exist {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_INVOICES, "Jwt Payload is Empty", nil)
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

	res, err := p.service.CreatePosInvoice(ctx, &req)
	if err != nil {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_INVOICES, err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	successResponse := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_CREATE_INVOICES, res)
	ctx.JSON(http.StatusOK, successResponse)
}

func (p *posInvoiceController) HandleReadPosInvoiceRequest(ctx *gin.Context) {
	invoiceID := ctx.Param("id")

	var req pb.ReadPosInvoiceRequest
	req.InvoiceId = invoiceID

	getJwtPayload, exist := ctx.Get("user")
	if !exist {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_INVOICES, "Jwt Payload is Empty", nil)
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

	res, err := p.service.ReadPosInvoice(ctx, &req)
	if err != nil {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_INVOICES, err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	successResponse := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_INVOICES, res)
	ctx.JSON(http.StatusOK, successResponse)
}

func (p *posInvoiceController) HandleUpdatePosInvoiceRequest(ctx *gin.Context) {
	var req pb.UpdatePosInvoiceRequest
	invoiceID := ctx.Param("id")
	if err := ctx.ShouldBindJSON(&req.PosInvoice); err != nil {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_INVOICES, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	getJwtPayload, exist := ctx.Get("user")
	if !exist {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_INVOICES, "Jwt Payload is Empty", nil)
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
	req.PosInvoice.InvoiceId = invoiceID
	res, err := p.service.UpdatePosInvoice(ctx, &req)
	if err != nil {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_INVOICES, err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	successResponse := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_UPDATE_INVOICES, res)
	ctx.JSON(http.StatusOK, successResponse)
}

func (p *posInvoiceController) HandleDeletePosInvoiceRequest(ctx *gin.Context) {
	invoiceID := ctx.Param("id")

	var req pb.DeletePosInvoiceRequest
	req.InvoiceId = invoiceID

	getJwtPayload, exist := ctx.Get("user")
	if !exist {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_INVOICES, "Jwt Payload is Empty", nil)
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

	_, err := p.service.DeletePosInvoice(ctx, &req)
	if err != nil {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_INVOICES, err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	successResponse := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_DELETE_INVOICES, nil)
	ctx.JSON(http.StatusOK, successResponse)
}

func (p *posInvoiceController) HandleReadAllPosInvoicesRequest(ctx *gin.Context) {
	limitQuery := ctx.Query("limit")
	pageQuery := ctx.Query("page")

	if (limitQuery == "" && pageQuery != "") || (limitQuery != "" && pageQuery == "") {
		errorResponse := utils.BuildResponseFailed("Both limit and page must be provided", "Value Is Empty", pb.CreatePosInvoiceResponse{})
		ctx.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	var req pb.ReadAllPosInvoicesRequest

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
			errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_INVOICES, "Jwt Payload is Empty", nil)
			ctx.JSON(http.StatusBadRequest, errorResponse)
			return
		}

		req = pb.ReadAllPosInvoicesRequest{
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
	res, err := p.service.ReadAllPosInvoices(ctx, &req)
	if err != nil {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_INVOICES, err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	successResponse := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_INVOICES, res)
	ctx.JSON(http.StatusOK, successResponse)
}
