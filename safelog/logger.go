package safelog

import (
	"bytes"
	"fmt"
	"io"
	"sync"

	"devicecert/domain"
)

type Logger struct {
	mu      sync.Mutex
	out     io.Writer
	entries []string
}

func New(out io.Writer) *Logger {
	return &Logger{out: out, entries: make([]string, 0)}
}

func (l *Logger) Write(event domain.AuditEvent) {
	if l == nil {
		return
	}
	line := FormatAudit(event)
	l.mu.Lock()
	l.entries = append(l.entries, line)
	if l.out != nil {
		_, _ = io.WriteString(l.out, line+"\n")
	}
	l.mu.Unlock()
}

func (l *Logger) Lines() []string {
	if l == nil {
		return nil
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	copyLines := make([]string, len(l.entries))
	copy(copyLines, l.entries)
	return copyLines
}

func (l *Logger) String() string {
	return bytes.NewBufferString(joinLines(l.Lines())).String()
}

func joinLines(lines []string) string {
	result := ""
	for index, line := range lines {
		if index > 0 {
			result += "\n"
		}
		result += line
	}
	return result
}

func (l *Logger) Count() int {
	if l == nil {
		return 0
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.entries)
}

func (l *Logger) ContainsPrivateKey(private string) bool {
	if private == "" {
		return false
	}
	for _, line := range l.Lines() {
		if line == private || contains(line, private) {
			return true
		}
	}
	return false
}

func contains(value, needle string) bool {
	for offset := 0; offset+len(needle) <= len(value); offset++ {
		if value[offset:offset+len(needle)] == needle {
			return true
		}
	}
	return false
}

func FormatCount(count int) string { return fmt.Sprintf("events=%d", count) }
