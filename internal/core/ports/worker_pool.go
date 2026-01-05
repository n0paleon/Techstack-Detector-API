package ports

type WorkerPool interface {
	Submit(fn func()) error
	Release()
}
