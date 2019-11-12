package batchy

const ErrBatcherStopped = err("Batcher Stopped")

type err string

func (e err) Error() string {
	return string(e)
}
