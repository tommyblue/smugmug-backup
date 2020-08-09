package testutil

import "github.com/sirupsen/logrus"

// LessLogging is a test helper that removes logging (in fact it sets its
// level to Panic). It returns a function which when called, resets it to its
// previous level. Its useful to be called as follows in test/benchmarks:
//
//  func TestFoo(t *testing.T) {
//      defer DisableLogging()()
//
//      // logging is set to Panic for the whole test
//  }
func DisableLogging() (reset func()) {
	lvl := logrus.GetLevel()
	logrus.SetLevel(logrus.PanicLevel)
	return func() { logrus.SetLevel(lvl) }
}

// LessLogging is a test helper that decreases logging (in fact it sets its
// level to Error). It returns a function which when called, resets it to its
// previous level. Its useful to be called as follows in test/benchmarks:
//
//  func TestFoo(t *testing.T) {
//      defer LessLogging()()
//
//      // logging is set to Error for the whole test
//  }
func LessLogging() (reset func()) {
	lvl := logrus.GetLevel()
	logrus.SetLevel(logrus.ErrorLevel)
	return func() { logrus.SetLevel(lvl) }
}
