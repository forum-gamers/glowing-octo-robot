package transaction

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func CheckCurrency(curr string) bool {
	for _, val := range []TransactionCurrency{
		RUPIAH, US_DOLLAR,
	} {
		if val == curr {
			return true
		}
	}
	return false
}

func CheckTransactionType(tr string) bool {
	for _, val := range []TransactionType{
		PAYMENT, TOP_UP,
	} {
		if val == tr {
			return true
		}
	}
	return false
}

func CheckTransactionStatus(transactionStatus TransactionStatus) error {
	switch transactionStatus {
	case COMPLETED:
		return status.Error(codes.FailedPrecondition, "transaction is already completed")
	case FAILED:
		return status.Error(codes.FailedPrecondition, "transaction is failed")
	case CANCEL:
		return status.Error(codes.FailedPrecondition, "transaction is already canceled")
	default:
		return nil
	}
}

func IsValidStatus(stat string) bool {
	for _, val := range []TransactionStatus{
		PENDING, COMPLETED, FAILED, CANCEL, REFUND, SETTLEMENT, DENY, EXPIRE,
	} {
		if val == stat {
			return true
		}
	}
	return false
}

func CountFee(amount float64, trType TransactionType) float64 {
	if trType == TOP_UP {
		return 0
	}

	return amount * 0.05
}
