package dto

import "errors"

// INVOICES Failed Messages
const (
	MESSAGE_FAILED_CREATE_INVOICES = "failed to create invoices"
	MESSAGE_FAILED_UPDATE_INVOICES = "failed to update invoices"
	MESSAGE_FAILED_DELETE_INVOICES = "failed to delete invoices"
	MESSAGE_FAILED_GET_INVOICES    = "failed to get invoices"
)

// INVOICES Success Messages
const (
	MESSAGE_SUCCESS_CREATE_INVOICES = "success create invoices"
	MESSAGE_SUCCESS_UPDATE_INVOICES = "success update invoices"
	MESSAGE_SUCCESS_DELETE_INVOICES = "success delete invoices"
	MESSAGE_SUCCESS_GET_INVOICES    = "success get invoices"
)

// INVOICES Custom Errors
var (
	ErrCreateInvoices = errors.New(MESSAGE_FAILED_CREATE_INVOICES)
	ErrUpdateInvoices = errors.New(MESSAGE_FAILED_UPDATE_INVOICES)
	ErrDeleteInvoices = errors.New(MESSAGE_FAILED_DELETE_INVOICES)
	ErrGetInvoices    = errors.New(MESSAGE_FAILED_GET_INVOICES)
)
