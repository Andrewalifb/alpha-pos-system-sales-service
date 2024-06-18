package dto

import "errors"

// CUSTOMER Failed Messages
const (
	MESSAGE_FAILED_CREATE_CUSTOMER = "failed to create customer"
	MESSAGE_FAILED_UPDATE_CUSTOMER = "failed to update customer"
	MESSAGE_FAILED_DELETE_CUSTOMER = "failed to delete customer"
	MESSAGE_FAILED_GET_CUSTOMER    = "failed to get customer"
)

// CUSTOMER Success Messages
const (
	MESSAGE_SUCCESS_CREATE_CUSTOMER = "success create customer"
	MESSAGE_SUCCESS_UPDATE_CUSTOMER = "success update customer"
	MESSAGE_SUCCESS_DELETE_CUSTOMER = "success delete customer"
	MESSAGE_SUCCESS_GET_CUSTOMER    = "success get customer"
)

// CUSTOMER Custom Errors
var (
	ErrCreateCustomer = errors.New(MESSAGE_FAILED_CREATE_CUSTOMER)
	ErrUpdateCustomer = errors.New(MESSAGE_FAILED_UPDATE_CUSTOMER)
	ErrDeleteCustomer = errors.New(MESSAGE_FAILED_DELETE_CUSTOMER)
	ErrGetCustomer    = errors.New(MESSAGE_FAILED_GET_CUSTOMER)
)
