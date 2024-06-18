package dto

import "errors"

// ONLINE_PAYMENT Failed Messages
const (
	MESSAGE_FAILED_CREATE_ONLINE_PAYMENT = "failed to create online payment"
	MESSAGE_FAILED_UPDATE_ONLINE_PAYMENT = "failed to update online payment"
	MESSAGE_FAILED_DELETE_ONLINE_PAYMENT = "failed to delete online payment"
	MESSAGE_FAILED_GET_ONLINE_PAYMENT    = "failed to get online payment"
)

// ONLINE_PAYMENT Success Messages
const (
	MESSAGE_SUCCESS_CREATE_ONLINE_PAYMENT = "success create online payment"
	MESSAGE_SUCCESS_UPDATE_ONLINE_PAYMENT = "success update online payment"
	MESSAGE_SUCCESS_DELETE_ONLINE_PAYMENT = "success delete online payment"
	MESSAGE_SUCCESS_GET_ONLINE_PAYMENT    = "success get online payment"
)

// ONLINE_PAYMENT Custom Errors
var (
	ErrCreateOnlinePayment = errors.New(MESSAGE_FAILED_CREATE_ONLINE_PAYMENT)
	ErrUpdateOnlinePayment = errors.New(MESSAGE_FAILED_UPDATE_ONLINE_PAYMENT)
	ErrDeleteOnlinePayment = errors.New(MESSAGE_FAILED_DELETE_ONLINE_PAYMENT)
	ErrGetOnlinePayment    = errors.New(MESSAGE_FAILED_GET_ONLINE_PAYMENT)
)
