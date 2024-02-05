package utils

const (
	USD = "USD"
	EUR = "EUR"
	UAH = "UAH"
)

func GetSupportedCurrencies() []string {
	return []string{USD, EUR, UAH}
}

func IsCurrencySupported(currency string) bool {
	for _, cur := range GetSupportedCurrencies() {
		if cur == currency {
			return true
		}
	}
	return false
}
