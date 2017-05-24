package vestigo

import "net/http"

type MockInterceptor struct {
	before          bool
	intercept       bool
	after           bool
	CalledIntercept int
}

func (m *MockInterceptor) Before() bool {
	return m.before
}

func (m *MockInterceptor) After() bool {
	return m.after
}

func (m *MockInterceptor) Intercept(w http.ResponseWriter, r *http.Request) bool {
	m.CalledIntercept += 1
	return m.intercept
}

func NewMockInterceptor(before, intercept, after bool) *MockInterceptor {
	return &MockInterceptor{before: before, intercept: intercept, after: after}
}
