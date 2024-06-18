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

type PosReturnController interface {
	HandleCreatePosReturnRequest(c *gin.Context)
	HandleReadPosReturnRequest(c *gin.Context)
	HandleUpdatePosReturnRequest(c *gin.Context)
	HandleDeletePosReturnRequest(c *gin.Context)
	HandleReadAllPosReturnsRequest(c *gin.Context)
}

type posReturnController struct {
	service pb.PosReturnServiceClient
}

func NewPosReturnController(service pb.PosReturnServiceClient) PosReturnController {
	return &posReturnController{
		service: service,
	}
}

func (p *posReturnController) HandleCreatePosReturnRequest(ctx *gin.Context) {
	var req pb.CreatePosReturnRequest

	if err := ctx.ShouldBindJSON(&req.PosReturn); err != nil {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_RETURN, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	getJwtPayload, exist := ctx.Get("user")
	if !exist {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_RETURN, "Jwt Payload is Empty", nil)
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

	res, err := p.service.CreatePosReturn(ctx, &req)
	if err != nil {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_RETURN, err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	successResponse := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_CREATE_RETURN, res)
	ctx.JSON(http.StatusOK, successResponse)
}

func (p *posReturnController) HandleReadPosReturnRequest(ctx *gin.Context) {
	returnID := ctx.Param("id")

	var req pb.ReadPosReturnRequest
	req.ReturnId = returnID

	getJwtPayload, exist := ctx.Get("user")
	if !exist {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_RETURN, "Jwt Payload is Empty", nil)
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

	res, err := p.service.ReadPosReturn(ctx, &req)
	if err != nil {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_RETURN, err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	successResponse := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_RETURN, res)
	ctx.JSON(http.StatusOK, successResponse)
}

func (p *posReturnController) HandleUpdatePosReturnRequest(ctx *gin.Context) {
	var req pb.UpdatePosReturnRequest
	returnID := ctx.Param("id")

	if err := ctx.ShouldBindJSON(&req.PosReturn); err != nil {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_RETURN, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	getJwtPayload, exist := ctx.Get("user")
	if !exist {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_RETURN, "Jwt Payload is Empty", nil)
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
	req.PosReturn.ReturnId = returnID
	res, err := p.service.UpdatePosReturn(ctx, &req)
	if err != nil {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_RETURN, err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	successResponse := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_UPDATE_RETURN, res)
	ctx.JSON(http.StatusOK, successResponse)
}

func (p *posReturnController) HandleDeletePosReturnRequest(ctx *gin.Context) {
	returnID := ctx.Param("id")

	var req pb.DeletePosReturnRequest
	req.ReturnId = returnID

	getJwtPayload, exist := ctx.Get("user")
	if !exist {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_RETURN, "Jwt Payload is Empty", nil)
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

	_, err := p.service.DeletePosReturn(ctx, &req)
	if err != nil {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_RETURN, err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	successResponse := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_DELETE_RETURN, nil)
	ctx.JSON(http.StatusOK, successResponse)
}

func (p *posReturnController) HandleReadAllPosReturnsRequest(ctx *gin.Context) {
	limitQuery := ctx.Query("limit")
	pageQuery := ctx.Query("page")

	if (limitQuery == "" && pageQuery != "") || (limitQuery != "" && pageQuery == "") {
		errorResponse := utils.BuildResponseFailed("Both limit and page must be provided", "Value Is Empty", pb.CreatePosReturnResponse{})
		ctx.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	var req pb.ReadAllPosReturnsRequest

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
			errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_RETURN, "Jwt Payload is Empty", nil)
			ctx.JSON(http.StatusBadRequest, errorResponse)
			return
		}

		req = pb.ReadAllPosReturnsRequest{
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

	res, err := p.service.ReadAllPosReturns(ctx, &req)
	if err != nil {
		errorResponse := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_RETURN, err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	successResponse := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_RETURN, res)
	ctx.JSON(http.StatusOK, successResponse)
}
