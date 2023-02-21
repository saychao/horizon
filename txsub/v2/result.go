package txsub

// Result - struct providing with all required info to consider result of submission
type Result struct {
	Hash           string
	LedgerSequence int32
	TransactionID  int64
	EnvelopeXDR    string
	ResultXDR      string
	ResultMetaXDR  string
}

type fullResult struct {
	Result
	Err error
}

func (r fullResult) unwrap() (*Result, error) {
	return &r.Result, r.Err
}
