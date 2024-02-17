package types

const (
	FailedVerify = iota
	Expired
)

var VerifyTokenErrMap = map[int]string{
	FailedVerify: "Not Existed token at server",
	Expired:      "Expired Token Login Again",
}
