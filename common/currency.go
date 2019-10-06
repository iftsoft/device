package common

type DevAmount   float32
type DevCounter  int32
type DevCurrency uint16

const (
	CurrencyNOT	DevCurrency = 0		// Currency not used - Валюта не используется
	CurrencyAMD	=	51		// Armenian Dram - Армянский драм
	CurrencyAUD	=	36		// Australian Dollar - Австралийский доллар
	CurrencyATS	=	40		// Austrian Schilling - Австрийский шиллинг
	CurrencyAZN	=	31		// Azerbaijan Manat - Азербайджанский манат
	CurrencyBYR	=	974		// Belarussian Ruble - Белорусский рубль
	CurrencyBEF	=	56		// Belgian Franc - Бельгийский франк
	CurrencyBRL	=	986		// Brazilian Real - Бразильский реал
	CurrencyGBP	=	826		// British Pound - Британский фунт
	CurrencyCAD	=	124		// Canadian Dollar - Канадский доллар
	CurrencyCNY	=	156		// China Yuan - Китайский юань
	CurrencyCZK	=	203		// Czech Koruna - Чешская крона
	CurrencyDKK	=	208		// Danish Krone - Датская крона
	CurrencyNLG	=	528		// Dutch Guilder - Нидерландский гульден
	CurrencyEEK	=	233		// Estonian Kroon - Эстонская крона
	CurrencyEUR	=	978		// European Euro - Единая европейская валюта
	CurrencyFIM	=	246		// Finnish Mark - Финская марка
	CurrencyFRF	=	250		// French Franc - Французский франк
	CurrencyGEL	=	981		// Georgia Lari - Грузинская лари
	CurrencyDEM	=	276		// German Mark - Немецкая марка
	CurrencyGRD	=	300		// Greek Drachma - Греческая драхма
	CurrencyHKD	=	344		// Hong Kong Dollar - Гонконгский доллар
	CurrencyHUF	=	348		// Hungarian Forint - Венгерский форинт
	CurrencyINR	=	356		// Indian Rupee - Индийская рупия
	CurrencyIEP	=	372		// Irish Punt - Ирландский фунт
	CurrencyILS	=	376		// Israeli Sheqel - Израильский шекель
	CurrencyITL	=	380		// Italian Lira - Итальянская лира
	CurrencyJPY	=	392		// Japanese Yen - Японская йена
	CurrencyKZT	=	398		// Kazakhstan Tenge - Казахский тенге
	CurrencyKGS	=	417		// Kyrgyzstan Som - Киргизский сом
	CurrencyLVL	=	428		// Latvian Lat - Латвийский лат
	CurrencyLTL	=	440		// Lithuanian Lita - Литовский лит
	CurrencyMDL	=	498		// Moldovan Leu - Молдавский лей
	CurrencyMXN	=	484		// Mexican Peso - Мексиканский песо
	CurrencyNZD	=	554		// New Zealand Dollar - Новозеландский доллар
	CurrencyNOK	=	578		// Norway Krone - Норвежская крона
	CurrencyPLN	=	985		// Polish Zloty - Польский злотый
	CurrencyPTE	=	620		// Portuguese Escudo - Португальское эскудо
	CurrencyROL	=	642		// Romania Leu - Румынский лей
	CurrencyRUB	=	643		// Russian Rouble - Российский рубль (новый)
	CurrencySGD	=	702		// Singapore Dollar - Сингапурский доллар
	CurrencySKK	=	703		// Slovak Koruna - Словацкая крона
	CurrencyZAR	=	710		// South African Rand - Южноафриканский ранд
	CurrencyESP	=	724		// Spanish Peseta - Испанская песета
	CurrencySEK	=	752		// Swedish Krona - Шведская крона
	CurrencyCHF	=	756		// Swiss Franc - Швейцарский франк
	CurrencyUAH	=	980		// Ukraine Hryvnia - Украинская гривна
	CurrencyUSD	=	840		// United States Dollar - Американский доллар
	CurrencyUZS	=	860		// Uzbekistan Sum - Узбекский сум
	CurrencyXAG	=	961		// Silver - Серебро
	CurrencyXAU	=	959		// Gold - Золото
	CurrencyXPD	=	964		// Palladium - Палладий
	CurrencyXPT	=	962		// Platinum - Платина
	CurrencyXXX	=	999		// Not currency (Token) - Нет валюты (жетон)
)

// String returns a string explaining of the currency code
func (e DevCurrency) String() string {
	switch e {
	case CurrencyNOT:		return "Currency not used"
	case CurrencyAMD:		return "Armenian Dram"			// Armenian Dram - Армянский драм
	case CurrencyAUD:		return "Australian Dollar"		// Australian Dollar - Австралийский доллар
	case CurrencyATS:		return "Austrian Schilling"		// Austrian Schilling - Австрийский шиллинг
	case CurrencyAZN:		return "Azerbaijan Manat"		// Azerbaijan Manat - Азербайджанский манат
	case CurrencyBYR:		return "Belarussian Ruble"		// Belarussian Ruble - Белорусский рубль
	case CurrencyBEF:		return "Belgian Franc"			// Belgian Franc - Бельгийский франк
	case CurrencyBRL:		return "Brazilian Real"			// Brazilian Real - Бразильский реал
	case CurrencyGBP:		return "British Pound"			// British Pound - Британский фунт
	case CurrencyCAD:		return "Canadian Dollar"		// Canadian Dollar - Канадский доллар
	case CurrencyCNY:		return "China Yuan"				// China Yuan - Китайский юань
	case CurrencyCZK:		return "Czech Koruna"			// Czech Koruna - Чешская крона
	case CurrencyDKK:		return "Danish Krone"			// Danish Krone - Датская крона
	case CurrencyNLG:		return "Dutch Guilder"			// Dutch Guilder - Нидерландский гульден
	case CurrencyEEK:		return "Estonian Kroon"			// Estonian Kroon - Эстонская крона
	case CurrencyEUR:		return "European Euro"			// European Euro - Единая европейская валюта
	case CurrencyFIM:		return "Finnish Mark"			// Finnish Mark - Финская марка
	case CurrencyFRF:		return "French Franc"			// French Franc - Французский франк
	case CurrencyGEL:		return "Georgia Lari"			// Georgia Lari - Грузинская лари
	case CurrencyDEM:		return "German Mark"			// German Mark - Немецкая марка
	case CurrencyGRD:		return "Greek Drachma"			// Greek Drachma - Греческая драхма
	case CurrencyHKD:		return "Hong Kong Dollar"		// Hong Kong Dollar - Гонконгский доллар
	case CurrencyHUF:		return "Hungarian Forint"		// Hungarian Forint - Венгерский форинт
	case CurrencyINR:		return "Indian Rupee"			// Indian Rupee - Индийская рупия
	case CurrencyIEP:		return "Irish Punt"				// Irish Punt - Ирландский фунт
	case CurrencyILS:		return "Israeli Sheqel"			// Israeli Sheqel - Израильский шекель
	case CurrencyITL:		return "Italian Lira"			// Italian Lira - Итальянская лира
	case CurrencyJPY:		return "Japanese Yen"			// Japanese Yen - Японская йена
	case CurrencyKZT:		return "Kazakhstan Tenge"		// Kazakhstan Tenge - Казахский тенге
	case CurrencyKGS:		return "Kyrgyzstan Som"			// Kyrgyzstan Som - Киргизский сом
	case CurrencyLVL:		return "Latvian Lat"			// Latvian Lat - Латвийский лат
	case CurrencyLTL:		return "Lithuanian Lita"		// Lithuanian Lita - Литовский лит
	case CurrencyMDL:		return "Moldovan Leu"			// Moldovan Leu - Молдавский лей
	case CurrencyMXN:		return "Mexican Peso"			// Mexican Peso - Мексиканский песо
	case CurrencyNZD:		return "New Zealand Dollar"		// New Zealand Dollar - Новозеландский доллар
	case CurrencyNOK:		return "Norway Krone"			// Norway Krone - Норвежская крона
	case CurrencyPLN:		return "Polish Zloty"			// Polish Zloty - Польский злотый
	case CurrencyPTE:		return "Portuguese Escudo"		// Portuguese Escudo - Португальское эскудо
	case CurrencyROL:		return "Romania Leu"			// Romania Leu - Румынский лей
	case CurrencyRUB:		return "Russian Rouble"			// Russian Rouble - Российский рубль (новый)
	case CurrencySGD:		return "Singapore Dollar"		// Singapore Dollar - Сингапурский доллар
	case CurrencySKK:		return "Slovak Koruna"			// Slovak Koruna - Словацкая крона
	case CurrencyZAR:		return "South African Rand"		// South African Rand - Южноафриканский ранд
	case CurrencyESP:		return "Spanish Peseta"			// Spanish Peseta - Испанская песета
	case CurrencySEK:		return "Swedish Krona"			// Swedish Krona - Шведская крона
	case CurrencyCHF:		return "Swiss Franc"			// Swiss Franc - Швейцарский франк
	case CurrencyUAH:		return "Ukraine Hryvnia"		// Ukraine Hryvnia - Украинская гривна
	case CurrencyUSD:		return "United States Dollar"	// United States Dollar - Американский доллар
	case CurrencyUZS:		return "Uzbekistan Sum"			// Uzbekistan Sum - Узбекский сум
	case CurrencyXAG:		return "Silver"					// Silver - Серебро
	case CurrencyXAU:		return "Gold"					// Gold - Золото
	case CurrencyXPD:		return "Palladium"				// Palladium - Палладий
	case CurrencyXPT:		return "Platinum"				// Platinum - Платина
	case CurrencyXXX:		return "Not currency (Token)"	// Not currency (Token) - Нет валюты (жетон)
	default:				return "Unknown Currency"
	}
}

// IsoCode returns the text ISO code of the currency code
func (e DevCurrency) IsoCode() string {
	switch e {
	case CurrencyNOT:		return "NOT"
	case CurrencyAMD:		return "AMD"		// Armenian Dram - Армянский драм
	case CurrencyAUD:		return "AUD"		// Australian Dollar - Австралийский доллар
	case CurrencyATS:		return "ATS"		// Austrian Schilling - Австрийский шиллинг
	case CurrencyAZN:		return "AZN"		// Azerbaijan Manat - Азербайджанский манат
	case CurrencyBYR:		return "BYR"		// Belarussian Ruble - Белорусский рубль
	case CurrencyBEF:		return "BEF"		// Belgian Franc - Бельгийский франк
	case CurrencyBRL:		return "BRL"		// Brazilian Real - Бразильский реал
	case CurrencyGBP:		return "GBP"		// British Pound - Британский фунт
	case CurrencyCAD:		return "CAD"		// Canadian Dollar - Канадский доллар
	case CurrencyCNY:		return "CNY"		// China Yuan - Китайский юань
	case CurrencyCZK:		return "CZK"		// Czech Koruna - Чешская крона
	case CurrencyDKK:		return "DKK"		// Danish Krone - Датская крона
	case CurrencyNLG:		return "NLG"		// Dutch Guilder - Нидерландский гульден
	case CurrencyEEK:		return "EEK"		// Estonian Kroon - Эстонская крона
	case CurrencyEUR:		return "EUR"		// European Euro - Единая европейская валюта
	case CurrencyFIM:		return "FIM"		// Finnish Mark - Финская марка
	case CurrencyFRF:		return "FRF"		// French Franc - Французский франк
	case CurrencyGEL:		return "GEL"		// Georgia Lari - Грузинская лари
	case CurrencyDEM:		return "DEM"		// German Mark - Немецкая марка
	case CurrencyGRD:		return "GRD"		// Greek Drachma - Греческая драхма
	case CurrencyHKD:		return "HKD"		// Hong Kong Dollar - Гонконгский доллар
	case CurrencyHUF:		return "HUF"		// Hungarian Forint - Венгерский форинт
	case CurrencyINR:		return "INR"		// Indian Rupee - Индийская рупия
	case CurrencyIEP:		return "IEP"		// Irish Punt - Ирландский фунт
	case CurrencyILS:		return "ILS"		// Israeli Sheqel - Израильский шекель
	case CurrencyITL:		return "ITL"		// Italian Lira - Итальянская лира
	case CurrencyJPY:		return "JPY"		// Japanese Yen - Японская йена
	case CurrencyKZT:		return "KZT"		// Kazakhstan Tenge - Казахский тенге
	case CurrencyKGS:		return "KGS"		// Kyrgyzstan Som - Киргизский сом
	case CurrencyLVL:		return "LVL"		// Latvian Lat - Латвийский лат
	case CurrencyLTL:		return "LTL"		// Lithuanian Lita - Литовский лит
	case CurrencyMDL:		return "MDL"		// Moldovan Leu - Молдавский лей
	case CurrencyMXN:		return "MXN"		// Mexican Peso - Мексиканский песо
	case CurrencyNZD:		return "NZD"		// New Zealand Dollar - Новозеландский доллар
	case CurrencyNOK:		return "NOK"		// Norway Krone - Норвежская крона
	case CurrencyPLN:		return "PLN"		// Polish Zloty - Польский злотый
	case CurrencyPTE:		return "PTE"		// Portuguese Escudo - Португальское эскудо
	case CurrencyROL:		return "ROL"		// Romania Leu - Румынский лей
	case CurrencyRUB:		return "RUB"		// Russian Rouble - Российский рубль (новый)
	case CurrencySGD:		return "SGD"		// Singapore Dollar - Сингапурский доллар
	case CurrencySKK:		return "SKK"		// Slovak Koruna - Словацкая крона
	case CurrencyZAR:		return "ZAR"		// South African Rand - Южноафриканский ранд
	case CurrencyESP:		return "ESP"		// Spanish Peseta - Испанская песета
	case CurrencySEK:		return "SEK"		// Swedish Krona - Шведская крона
	case CurrencyCHF:		return "CHF"		// Swiss Franc - Швейцарский франк
	case CurrencyUAH:		return "UAH"		// Ukraine Hryvnia - Украинская гривна
	case CurrencyUSD:		return "USD"		// United States Dollar - Американский доллар
	case CurrencyUZS:		return "UZS"		// Uzbekistan Sum - Узбекский сум
	case CurrencyXAG:		return "XAG"		// Silver - Серебро
	case CurrencyXAU:		return "XAU"		// Gold - Золото
	case CurrencyXPD:		return "XPD"		// Palladium - Палладий
	case CurrencyXPT:		return "XPT"		// Platinum - Платина
	case CurrencyXXX:		return "XXX"		// Not currency (Token) - Нет валюты (жетон)
	default:				return "???"
	}
}

