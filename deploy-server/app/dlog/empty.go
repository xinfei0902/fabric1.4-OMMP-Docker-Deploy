package dlog

type NullOutput struct{}

func (*NullOutput) Write(p []byte) (n int, err error) {
	return len(p), nil
}
