package batchy

// ErrBatcherStopped indicates that the batcher is not accepting new items
//   because it has been stopped
const ErrBatcherStopped = err("Batcher Stopped")

type err string

func (e err) Error() string {
	return string(e)
}
