package mocks

import mock "github.com/stretchr/testify/mock"

type Logger struct {
	mock.Mock
}

func (l *Logger) Printf(format string, v ...any) {
	l.Called(format, v)
}
