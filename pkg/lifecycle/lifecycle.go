package lifecycle

// LifeCycle service's lifecycle
type LifeCycle interface {
	Start() error
	Stop() error
}
