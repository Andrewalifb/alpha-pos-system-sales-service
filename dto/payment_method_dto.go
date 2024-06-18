package dto

import "errors"

// PAYMENT_METHOD Failed Messages
const (
	MESSAGE_FAILED_CREATE_PAYMENT_METHOD = "failed to create payment method"
	MESSAGE_FAILED_UPDATE_PAYMENT_METHOD = "failed to update payment method"
	MESSAGE_FAILED_DELETE_PAYMENT_METHOD = "failed to delete payment method"
	MESSAGE_FAILED_GET_PAYMENT_METHOD    = "failed to get payment method"
)

// PAYMENT_METHOD Success Messages
const (
	MESSAGE_SUCCESS_CREATE_PAYMENT_METHOD = "success create payment method"
	MESSAGE_SUCCESS_UPDATE_PAYMENT_METHOD = "success update payment method"
	MESSAGE_SUCCESS_DELETE_PAYMENT_METHOD = "success delete payment method"
	MESSAGE_SUCCESS_GET_PAYMENT_METHOD    = "success get payment method"
)

// PAYMENT_METHOD Custom Errors
var (
	ErrCreatePaymentMethod = errors.New(MESSAGE_FAILED_CREATE_PAYMENT_METHOD)
	ErrUpdatePaymentMethod = errors.New(MESSAGE_FAILED_UPDATE_PAYMENT_METHOD)
	ErrDeletePaymentMethod = errors.New(MESSAGE_FAILED_DELETE_PAYMENT_METHOD)
	ErrGetPaymentMethod    = errors.New(MESSAGE_FAILED_GET_PAYMENT_METHOD)
)
