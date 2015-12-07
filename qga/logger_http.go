package qga

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"
)

const (
	agentLogger = "http://169.254.169.254/agent/log"
)

// Logger struct
type Logger struct {
	w *http.Client
}

// NewLogger create new http logger
func NewLogger(c *http.Client) (*Logger, error) {
	l := &Logger{}
	if c == nil {
		httpTransport := &http.Transport{
			Dial:            (&net.Dialer{DualStack: true}).Dial,
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		dt, err := time.ParseDuration("10s")
		if err != nil {
			return nil, err
		}
		l.w = &http.Client{Transport: httpTransport, Timeout: dt}
	} else {
		l.w = c
	}

	return l, nil
}

// Close closes logger
func (l *Logger) Close() error {
	if l.w == nil {
		return nil
	}
	return nil
}

// Debug send debug message
func (l *Logger) Debug(msg string) error {
	if l.w == nil {
		return nil
	}
	_, err := l.w.Post(agentLogger, "text/plain", bytes.NewBufferString("debug: "+msg))
	return err
}

// Debugf send debug formatted message
func (l *Logger) Debugf(f string, msg string) error {
	return l.Debug(fmt.Sprintf(f, msg))
}

// Error send error message
func (l *Logger) Error(msg string) error {
	if l.w == nil {
		return nil
	}
	_, err := l.w.Post(agentLogger, "text/plain", bytes.NewBufferString("error: "+msg))
	return err
}

// Errorf send error formatted message
func (l *Logger) Errorf(f string, msg string) error {
	return l.Error(fmt.Sprintf(f, msg))
}

// Info send info message
func (l *Logger) Info(msg string) error {
	if l.w == nil {
		return nil
	}
	_, err := l.w.Post(agentLogger, "text/plain", bytes.NewBufferString("info: "+msg))
	return err
}

// Infof send info formatted message
func (l *Logger) Infof(f string, msg string) error {
	return l.Info(fmt.Sprintf(f, msg))
}
