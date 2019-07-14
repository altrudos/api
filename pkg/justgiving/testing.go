package justgiving

type fixtures struct {
	CharityId             int
	DonationReferenceCode string
}

var Fixtures = fixtures{
	CharityId:             2050,
	DonationReferenceCode: "ch-1559330147572782000",
}

//TODO: Use env to override maybe?
func GetTestAppId() string {
	return "a7f36da5"
}

func GetTestJG() *JustGiving {
	var JG = &JustGiving{
		Mode:  ModeStaging,
		AppId: GetTestAppId(),
	}
	return JG
}
