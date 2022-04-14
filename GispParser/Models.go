package main

type Org struct {
	Name   string
	INN    int
	OGRN   int
	Adress string
}

type Prod struct {
	registryNumber  string
	Name            string
	OKPD2           int
	TNVED           int
	madeAccordingTo string
	Points          int
}
