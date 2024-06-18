package dto

import "errors"

// RETURN Failed Messages
const (
	MESSAGE_FAILED_CREATE_RETURN = "failed to create return"
	MESSAGE_FAILED_UPDATE_RETURN = "failed to update return"
	MESSAGE_FAILED_DELETE_RETURN = "failed to delete return"
	MESSAGE_FAILED_GET_RETURN    = "failed to get return"
)

// RETURN Success Messages
const (
	MESSAGE_SUCCESS_CREATE_RETURN = "success create return"
	MESSAGE_SUCCESS_UPDATE_RETURN = "success update return"
	MESSAGE_SUCCESS_DELETE_RETURN = "success delete return"
	MESSAGE_SUCCESS_GET_RETURN    = "success get return"
)

// RETURN Custom Errors
var (
	ErrCreateReturn = errors.New(MESSAGE_FAILED_CREATE_RETURN)
	ErrUpdateReturn = errors.New(MESSAGE_FAILED_UPDATE_RETURN)
	ErrDeleteReturn = errors.New(MESSAGE_FAILED_DELETE_RETURN)
	ErrGetReturn    = errors.New(MESSAGE_FAILED_GET_RETURN)
)
