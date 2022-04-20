package main

type Org struct {
	FullName                   string
	ShortName                  string
	INN                        int
	OGRN                       int
	KPP                        int
	Rating                     int
	INDUSTRY                   string
	Country                    string
	Region                     string
	City                       string
	Adress                     string
	Index                      int
	GenDirector                string
	ElectronicTradingPlatforms string
	www                        string
}

type Prod struct {
	RegistyNumber   string
	Name            string
	OKPD2           int
	TNVED           int
	madeAccordingTo string
	Points          int
}

type URL struct {
	Org   string
	Prods string
}
