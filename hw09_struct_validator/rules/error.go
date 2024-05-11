package rules

type ErrCode int

const (
	ErrMin ErrCode = iota
	ErrMax
	ErrLength
	ErrRegexp
	ErrIn
)

func (e ErrCode) Error() string {
	switch e {
	case ErrMin:
		return "less than min"
	case ErrMax:
		return "greater than max"
	case ErrLength:
		return "wrong length"
	case ErrRegexp:
		return "not match regexp"
	case ErrIn:
		return "not in enum"
	default:
		panic("unsupported errCode")
	}
}
