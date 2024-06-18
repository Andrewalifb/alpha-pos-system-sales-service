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

type PosSaleController interface {
	HandleCreatePosSaleRequest(c *gin.Context)
	HandleReadPosSaleRequest(c *gin.Context)
	HandleUpdatePosSaleRequest(c *gin.Context)
	HandleDeletePosSaleRequest(c *gin.Context)
	HandleReadAllPosSalesRequest(c *gin.Context)
}

type posSaleController struct {
	service pb.PosSaleServiceClient
}

func NewPosSaleController(service pb.PosSaleServiceClient) PosSaleController {
	return &posSaleController{
		service: service,
	}
}

func (p *posSaleController) HandleCreatePosSaleRequest(ctx *gin.Context) {
	var req pb.CreatePosSalesRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_SALES, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	getJwtPayload, exist := ctx.Get("user")
	if !exist {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_SALES, "Jwt Payload is Empty", nil)
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

	res, err := p.service.CreatePosSales(ctx, &req)
	if err != nil {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_SALES, err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	successResponse := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_CREATE_SALES, res)
	ctx.JSON(http.StatusOK, successResponse)
}

func (p *posSaleController) HandleReadPosSaleRequest(ctx *gin.Context) {
	saleID := ctx.Param("id")

	var req pb.ReadPosSaleRequest
	req.SaleId = saleID

	getJwtPayload, exist := ctx.Get("user")
	if !exist {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_SALES, "Jwt Payload is Empty", nil)
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

	res, err := p.service.ReadPosSale(ctx, &req)
	if err != nil {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_SALES, err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	successResponse := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_SALES, res)
	ctx.JSON(http.StatusOK, successResponse)
}

func (p *posSaleController) HandleUpdatePosSaleRequest(ctx *gin.Context) {
	var req pb.UpdatePosSaleRequest
	saleID := ctx.Param("id")
	if err := ctx.ShouldBindJSON(&req.PosSale); err != nil {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_SALES, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	getJwtPayload, exist := ctx.Get("user")
	if !exist {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_SALES, "Jwt Payload is Empty", nil)
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
	req.PosSale.SaleId = saleID
	res, err := p.service.UpdatePosSale(ctx, &req)
	if err != nil {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_SALES, err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	successResponse := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_UPDATE_SALES, res)
	ctx.JSON(http.StatusOK, successResponse)
}

func (p *posSaleController) HandleDeletePosSaleRequest(ctx *gin.Context) {
	saleID := ctx.Param("id")

	var req pb.DeletePosSaleRequest
	req.SaleId = saleID

	getJwtPayload, exist := ctx.Get("user")
	if !exist {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_SALES, "Jwt Payload is Empty", nil)
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

	_, err := p.service.DeletePosSale(ctx, &req)
	if err != nil {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_SALES, err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	successResponse := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_DELETE_SALES, nil)
	ctx.JSON(http.StatusOK, successResponse)
}

func (p *posSaleController) HandleReadAllPosSalesRequest(ctx *gin.Context) {
	limitQuery := ctx.Query("limit")
	pageQuery := ctx.Query("page")

	if (limitQuery == "" && pageQuery != "") || (limitQuery != "" && pageQuery == "") {
		errorResponse := utils.BuildResponseFailed("Both limit and page must be provided", "Value Is Empty", pb.CreatePosSalesResponse{})
		ctx.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	var req pb.ReadAllPosSalesRequest

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
			errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_SALES, "Jwt Payload is Empty", nil)
			ctx.JSON(http.StatusBadRequest, errorResponse)
			return
		}

		req = pb.ReadAllPosSalesRequest{
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
	res, err := p.service.ReadAllPosSales(ctx, &req)
	if err != nil {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_SALES, err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	successResponse := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_SALES, res)
	ctx.JSON(http.StatusOK, successResponse)
}
