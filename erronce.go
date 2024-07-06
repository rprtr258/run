package run

type Once struct {
	err error
}

func NewErrOnce() *Once {
	return &Once{}
}

func (o *Once) Run(f func() error) *Once {
	if o.err == nil {
		o.err = f()
	}
	return o
}

func Init[T any](o *Once, f func() (T, error)) T {
	if o.err != nil {
		return *new(T)
	}
	var t T
	t, o.err = f()
	return t
}

func (o *Once) Error() error {
	return o.err
}
