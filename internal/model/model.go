package model

//Data ...
type Data struct {
	Cpf           string
	Private       string
	Incomplete    string
	LastPurchase  string
	AvgTicket     string
	LastTicket    string
	FrequentStore string
	LastStore     string
}

//MapData ...
type MapData struct {
	Key   int
	Value []string
}
