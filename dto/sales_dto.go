package dto

import "errors"

// SALES Failed Messages
const (
	MESSAGE_FAILED_CREATE_SALES = "failed to create sales"
	MESSAGE_FAILED_UPDATE_SALES = "failed to update sales"
	MESSAGE_FAILED_DELETE_SALES = "failed to delete sales"
	MESSAGE_FAILED_GET_SALES    = "failed to get sales"
)

// SALES Success Messages
const (
	MESSAGE_SUCCESS_CREATE_SALES = "success create sales"
	MESSAGE_SUCCESS_UPDATE_SALES = "success update sales"
	MESSAGE_SUCCESS_DELETE_SALES = "success delete sales"
	MESSAGE_SUCCESS_GET_SALES    = "success get sales"
)

// SALES Custom Errors
var (
	ErrCreateSales = errors.New(MESSAGE_FAILED_CREATE_SALES)
	ErrUpdateSales = errors.New(MESSAGE_FAILED_UPDATE_SALES)
	ErrDeleteSales = errors.New(MESSAGE_FAILED_DELETE_SALES)
	ErrGetSales    = errors.New(MESSAGE_FAILED_GET_SALES)
)
