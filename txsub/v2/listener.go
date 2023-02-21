package txsub

type listener chan<- fullResult

func send(l chan fullResult, res fullResult) <-chan fullResult {
	l <- res
	close(l)
	return l
}
