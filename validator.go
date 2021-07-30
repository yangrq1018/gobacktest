package gobacktest

type Validator interface {
	Validate(FillEvent, PortfolioHandler) (bool, error)
}
