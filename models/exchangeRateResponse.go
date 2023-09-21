package models 

type ExchangeRateResponse struct {
	Date string `json:"date"`
	Base string `json:"base"`
	Rates map[string]string `json:"rates"`
}

type ExchangeRateMessage struct {
	Currency string `json:"currency"`
	Url string `json:"url"`
}

type ExchangeRateResult struct {
	Date string `json:"date"`
	Base string `json:"base"`
	Rate string `json:"rate"`
}