package util

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/dustin/go-humanize"
)

const (
	LanguageVN               = "vi-VN"
	LanguageEN               = "en-US"
	DateLayoutVN             = "02/01/2006"
	DateLayoutEN             = "02-Jan-2006"
	TimeLayoutMinuteIn24Hour = "15:04"
	TimeLayoutSecondIn24Hour = "15:04:05"

	CurrencyVND = "VND"
	CurrencyUSD = "USD"
	CurrencyAUD = "AUD"
	CurrencyEUR = "EUR"
	CurrencyCAD = "CAD"
	CurrencyGBP = "GBP"
	CurrencyJPY = "JPY"
	CurrencySGD = "SGD"
	CurrencyCHF = "CHF"

	SymbolCurrencyVND = "đ"
	SymbolCurrencyUSD = "$"
	SymbolCurrencyAUD = "A$"
	SymbolCurrencyEUR = "€"
	SymbolCurrencyCAD = "C$"
	SymbolCurrencyGBP = "£"
	SymbolCurrencyJPY = "¥"
	SymbolCurrencySGD = "S$"
	SymbolCurrencyCHF = "₣"
)

var (
	currencies = map[string]string{
		CurrencyVND: SymbolCurrencyVND,
		CurrencyUSD: SymbolCurrencyUSD,
		CurrencyAUD: SymbolCurrencyAUD,
		CurrencyEUR: SymbolCurrencyEUR,
		CurrencyCAD: SymbolCurrencyCAD,
		CurrencyGBP: SymbolCurrencyGBP,
		CurrencyJPY: SymbolCurrencyJPY,
		CurrencySGD: SymbolCurrencySGD,
		CurrencyCHF: SymbolCurrencyCHF,
	}
)

func FormatDate(t time.Time, language string) string {
	switch language {
	case LanguageEN:
		return t.Format(DateLayoutEN)
	case LanguageVN:
		return t.Format(DateLayoutVN)
	}
	return ""
}

func FormatDateTimeLayoutMinuteIn24H(t time.Time, language string) string {
	switch language {
	case LanguageEN:
		return t.Format(fmt.Sprintf("%s %s", DateLayoutEN, TimeLayoutMinuteIn24Hour))
	case LanguageVN:
		return t.Format(fmt.Sprintf("%s %s", DateLayoutVN, TimeLayoutMinuteIn24Hour))
	}
	return ""
}

func FormatDateTimeLayoutSecondIn24H(t time.Time, language string) string {
	switch language {
	case LanguageEN:
		return t.Format(fmt.Sprintf("%s %s", DateLayoutEN, TimeLayoutSecondIn24Hour))
	case LanguageVN:
		return t.Format(fmt.Sprintf("%s %s", DateLayoutVN, TimeLayoutSecondIn24Hour))
	}
	return ""
}

func FormatTimeLayoutMinuteIn24H(t time.Time) string {
	return t.Format(TimeLayoutMinuteIn24Hour)
}

func FormatTimeLayoutSecondIn24H(t time.Time, language string) string {
	return t.Format(TimeLayoutSecondIn24Hour)
}

func ParseAndFormatMoney(value, currency string) string {
	v, _ := strconv.ParseFloat(value, 64)
	return FormatMoney(v, currency)
}

func FormatMoney(value float64, currency string) string {
	switch currency {
	case CurrencyVND:
		return humanize.CommafWithDigits(math.Round(value), 0)
	default:
		return humanize.CommafWithDigits(value, 2)
	}
}

func FormatCurrency(value, currency string) (formattedCurrency string) {
	formattedCurrency = ParseAndFormatMoney(value, currency)
	if v, ok := currencies[currency]; ok {
		formattedCurrency += " " + v
	}
	return formattedCurrency
}

func FormatCurrencyFromFloat(value float64, currency string) (formattedCurrency string) {
	formattedCurrency = FormatMoney(value, currency)
	if v, ok := currencies[currency]; ok {
		formattedCurrency += " " + v
	}
	return formattedCurrency
}
