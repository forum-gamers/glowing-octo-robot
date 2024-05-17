package transaction

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
