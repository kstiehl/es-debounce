package helper

type TestStream struct {
	StreamName string
}

func (t TestStream) Name() string {
	return t.StreamName
}
