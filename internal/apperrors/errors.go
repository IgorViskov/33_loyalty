package apperrors

import "errors"

var (
	ErrNotValidJSON                = errors.New("not valid json")
	ErrProcessStatusShouldBeString = errors.New("process status should be string")
	ErrInsertConflict              = errors.New("insert conflict")
	ErrNonImplemented              = errors.New("not implemented")
	InfoEmptyRunHost               = "run host is empty"
	ErrNotValidAccrualHost         = errors.New("accrual host not valid")
	ErrDBURIIsEmpty                = errors.New("db uri is empty")
	ErrPairLoginPasswordNotValid   = errors.New("pair login password not valid")
	ErrInvalidFormatRequest        = errors.New("invalid format request")
	ErrRecordNotFound              = errors.New("record not found")
	MsgOrderEasUploadedAnotherUser = "this number was uploaded by another user"
	MsgIncorrectOrderNumber        = "incorrect order number"
	ErrWithdrawalsNegative         = errors.New("accrual negative")
	ErrAmountNotPositive           = errors.New("amount not positive")
)
