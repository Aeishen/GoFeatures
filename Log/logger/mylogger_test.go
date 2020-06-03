package logger

import "testing"

func TestConstLevel(t *testing.T){
	t.Logf("%v, %T", DebugLevel, DebugLevel)
	t.Logf("%v, %T", InfoLevel, InfoLevel)
	t.Logf("%v, %T", WarnLevel, WarnLevel)
	t.Logf("%v, %T", ErrorLevel, ErrorLevel)
	t.Logf("%v, %T", FatalLevel, FatalLevel)

}
//go test -v