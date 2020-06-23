package klog

import (
	"strconv"
	"testing"
)

type testHook struct {
	severityLevel string
	fired         map[string]int
}

func newTestHook(severityLevel string) *testHook {
	return &testHook{
		severityLevel: severityLevel,
		fired:         make(map[string]int),
	}
}

func (t *testHook) SeverityLevel() string {
	return t.severityLevel
}

func (t *testHook) Fire(s string, args ...interface{}) error {
	if _, ok := t.fired[s]; !ok {
		t.fired[s] = 1
		return nil
	}
	t.fired[s]++

	return nil
}

// validate takes the execpted value, and returns if the result
// was favourable and the got value as a string
func (t *testHook) validateFiring(
	expectedKey string, expectedValue int) (bool, string) {

	val, ok := t.fired[expectedKey]
	if !ok {
		return false, ""
	}

	if val != expectedValue {
		return false, strconv.Itoa(val)
	}

	return true, strconv.Itoa(val)
}

func TestCanFireHooks(t *testing.T) {
	hook := newTestHook(InfoSeverityLevel)

	AddHook(hook)

	Info("fire info 1")
	Info("fire info 2")
	expected, got := hook.validateFiring(InfoSeverityLevel, 2)
	if !expected {
		t.Errorf("unexpected firing result: got:\n\t%s \nwant:\t%d", got, 2)
	}

	Error("fire error 3")
	expected, got = hook.validateFiring(ErrorSeverityLevel, 1)
	if !expected {
		t.Errorf("unexpected firing result: got:\n\t%s \nwant:\t%d", got, 1)
	}

	_, ok := hook.fired[FatalSeverityLevel]
	if ok {
		t.Errorf("hook firing %s was not expected", FatalSeverityLevel)
	}
}

func TestAddingUnsupportedHook(t *testing.T) {
	severityLevel := "PANIC"

	hook := newTestHook(severityLevel)

	hooks := Hooks{}
	err := hooks.Add(hook)
	if err == nil {
		t.Errorf(
			"error was expected as severity: %s is not supported",
			severityLevel,
		)
	}
}

func TestCanFireMultipleHooks(t *testing.T) {
	hook1 := newTestHook(InfoSeverityLevel)
	hook2 := newTestHook(ErrorSeverityLevel)

	AddHook(hook1)
	AddHook(hook2)

	Info("fire info 1")
	Error("fire error 1")
	expected, got := hook1.validateFiring(InfoSeverityLevel, 1)
	if !expected {
		t.Errorf("unexpected firing result: got:\n\t%s \nwant:\t%d", got, 1)
	}
	expected, got = hook1.validateFiring(ErrorSeverityLevel, 1)
	if !expected {
		t.Errorf("unexpected firing result: got:\n\t%s \nwant:\t%d", got, 1)
	}
	expected, got = hook2.validateFiring(ErrorSeverityLevel, 1)
	if !expected {
		t.Errorf("unexpected firing result: got:\n\t%s \nwant:\t%d", got, 1)
	}
	_, ok := hook2.fired[InfoSeverityLevel]
	if ok {
		t.Errorf("hook2 firing %s was not expected", InfoSeverityLevel)
	}
}
