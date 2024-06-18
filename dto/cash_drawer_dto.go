package dto

import "errors"

// CASH_DRAWER Failed Messages
const (
	MESSAGE_FAILED_CREATE_CASH_DRAWER = "failed to create cash drawer"
	MESSAGE_FAILED_UPDATE_CASH_DRAWER = "failed to update cash drawer"
	MESSAGE_FAILED_DELETE_CASH_DRAWER = "failed to delete cash drawer"
	MESSAGE_FAILED_GET_CASH_DRAWER    = "failed to get cash drawer"
)

// CASH_DRAWER Success Messages
const (
	MESSAGE_SUCCESS_CREATE_CASH_DRAWER = "success create cash drawer"
	MESSAGE_SUCCESS_UPDATE_CASH_DRAWER = "success update cash drawer"
	MESSAGE_SUCCESS_DELETE_CASH_DRAWER = "success delete cash drawer"
	MESSAGE_SUCCESS_GET_CASH_DRAWER    = "success get cash drawer"
)

// CASH_DRAWER Custom Errors
var (
	ErrCreateCashDrawer = errors.New(MESSAGE_FAILED_CREATE_CASH_DRAWER)
	ErrUpdateCashDrawer = errors.New(MESSAGE_FAILED_UPDATE_CASH_DRAWER)
	ErrDeleteCashDrawer = errors.New(MESSAGE_FAILED_DELETE_CASH_DRAWER)
	ErrGetCashDrawer    = errors.New(MESSAGE_FAILED_GET_CASH_DRAWER)
)
