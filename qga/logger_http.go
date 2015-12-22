package qga

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"
)

const (
	agentLogger = "http://169.254.169.254/agent/log"

//	agentLogger = "http://api.ix.clodo.ru/servers/agent_log/?vps=5591-444&token=a708b02ff3df5eef61d70254b7ee3354"
)

// Logger struct
type Logger struct {
	w      *http.Client
	buffer *bytes.Buffer
	//	cr     *gzip.Reader
	//	cw     *gzip.Writer
	m sync.Mutex
}

// NewLogger create new http logger
func NewLogger(c *http.Client) *Logger {
	l := &Logger{}
	if c == nil {
		httpTransport := &http.Transport{
			Dial:            (&net.Dialer{DualStack: true}).Dial,
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		dt, _ := time.ParseDuration("20s")
		l.w = &http.Client{Transport: httpTransport, Timeout: dt}
	} else {
		l.w = c
	}

	l.buffer = bytes.NewBuffer(nil)
	//	l.cr, _ = gzip.NewReader(l.buffer)
	//	l.cw = gzip.NewWriter(l.buffer)
	return l
}

// Close closes logger
func (l *Logger) Close() error {
	l.buffer.Reset()
	//	l.cr.Close()
	//	l.cw.Close()
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
	l.m.Lock()
	defer l.m.Unlock()
	l.reset()
	l.buffer.Write([]byte("debug: " + msg))
	_, err := l.w.Post(agentLogger, "text/plain", l.buffer)
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
	l.m.Lock()
	defer l.m.Unlock()
	l.reset()
	l.buffer.Write([]byte("error: " + msg))
	_, err := l.w.Post(agentLogger, "text/plain", l.buffer)
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
	l.m.Lock()
	defer l.m.Unlock()
	l.reset()
	l.buffer.Write([]byte("info: " + msg))
	_, err := l.w.Post(agentLogger, "text/plain", l.buffer)
	return err
}

// Infof send info formatted message
func (l *Logger) Infof(f string, msg string) error {
	return l.Info(fmt.Sprintf(f, msg))
}

func (l *Logger) reset() {
	if l.buffer.Len() > 0 {
		l.buffer.Reset()
		//l.cw.Reset(l.buffer)
		//l.cr.Reset(l.buffer)
	}
}
