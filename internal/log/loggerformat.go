package log

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// ANSI color codes for Spring Boot-style logging
const (
	ColorReset   = "\033[0m"
	ColorRed     = "\033[31m"
	ColorGreen   = "\033[32m"
	ColorYellow  = "\033[33m"
	ColorBlue    = "\033[34m"
	ColorPurple  = "\033[35m"
	ColorCyan    = "\033[36m"
	ColorWhite   = "\033[37m"
	ColorBold    = "\033[1m"
	ColorBgRed   = "\033[41m"
	ColorBgGreen = "\033[42m"
)

// SpringBootHandler implements a custom slog handler that formats logs
// similar to Spring Boot's logback pattern:
// [%boldGreen(%d{yyyy-MM-dd HH:mm:ss.SSS})] %highlight(%-5level) %magenta([%thread]) %cyan(%logger) [%boldYellow(%M:%L)] : %msg%n%throwable
type SpringBootHandler struct {
	writer      io.Writer
	level       slog.Level
	attrs       []slog.Attr
	group       string
	projectRoot string
}

func NewFormattedLogHandler(w io.Writer, level slog.Level) *SpringBootHandler {
	return &SpringBootHandler{
		writer:      w,
		level:       level,
		attrs:       make([]slog.Attr, 0),
		projectRoot: getProjectRoot(),
	}
}

// getProjectRoot finds the project root directory by looking for go.mod file
func getProjectRoot() string {
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}

	dir := wd
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return wd
}

// Enabled returns true if the handler should log at the given level
func (h *SpringBootHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level
}

// Handle formats and writes the log record in Spring Boot style
// Format: [timestamp] LEVEL [thread] logger [filename:method:line] : message
func (h *SpringBootHandler) Handle(_ context.Context, r slog.Record) error {
	// Get caller information
	var funcName, fileName, packageName, relativeFileName string
	var lineNum int

	if r.PC != 0 {
		frame, _ := runtime.CallersFrames([]uintptr{r.PC}).Next()
		funcName = frame.Function
		fileName = frame.File
		lineNum = frame.Line

		// Extract package and function name
		if idx := strings.LastIndex(funcName, "/"); idx != -1 {
			funcName = funcName[idx+1:]
		}
		if idx := strings.Index(funcName, "."); idx != -1 {
			packageName = funcName[:idx]
			funcName = funcName[idx+1:]
		}

		// Get relative path from project root
		if h.projectRoot != "" {
			if rel, err := filepath.Rel(h.projectRoot, fileName); err == nil {
				relativeFileName = rel
			} else {
				relativeFileName = filepath.Base(fileName)
			}
		} else {
			relativeFileName = filepath.Base(fileName)
		}

		// Remove .go extension for cleaner output
		if strings.HasSuffix(relativeFileName, ".go") {
			relativeFileName = relativeFileName[:len(relativeFileName)-3]
		}
	}

	var logLine strings.Builder

	// 1. Format timestamp with bold green color - [%boldGreen(%d{yyyy-MM-dd HH:mm:ss.SSS})]
	timestamp := r.Time.Format("2006-01-02 15:04:05.000")
	logLine.WriteString(fmt.Sprintf("[%s%s%s%s]",
		ColorBold, ColorGreen, timestamp, ColorReset))

	// 2. Format log level with highlight colors - %highlight(%-5level)
	levelStr := h.formatLevel(r.Level)
	logLine.WriteString(fmt.Sprintf(" %s", levelStr))

	// 3. Format thread (goroutine) info - %magenta([%thread])
	goroutineID := h.getGoroutineID()
	threadInfo := fmt.Sprintf(" %s[%s]%s", ColorPurple, goroutineID, ColorReset)
	logLine.WriteString(threadInfo)

	// 4. Format logger name (package name) - %cyan(%logger)
	logger := packageName
	if logger == "" {
		logger = "main"
	}
	loggerName := fmt.Sprintf(" %s%s%s", ColorCyan, logger, ColorReset)
	logLine.WriteString(loggerName)

	// 5. Format full-filename:method:line - [%boldYellow(full-filename:method:line)]
	if relativeFileName != "" && funcName != "" && lineNum > 0 {
		// Clean up function name to remove receiver type
		cleanFuncName := h.cleanFunctionName(funcName)
		methodInfo := fmt.Sprintf(" [%s%s:%s:%d%s]",
			ColorBold+ColorYellow, relativeFileName, cleanFuncName, lineNum, ColorReset)
		logLine.WriteString(methodInfo)
	}

	// 6. Message separator and content - : %msg
	logLine.WriteString(" : ")

	// Build the log message with attributes
	msg := r.Message

	// Add handler-level attributes
	var allAttrs []slog.Attr
	allAttrs = append(allAttrs, h.attrs...)

	// Add record-level attributes
	r.Attrs(func(a slog.Attr) bool {
		allAttrs = append(allAttrs, a)
		return true
	})

	// Format attributes
	if len(allAttrs) > 0 {
		var attrStrs []string
		for _, attr := range allAttrs {
			attrStrs = append(attrStrs, h.formatAttribute(attr))
		}
		if len(attrStrs) > 0 {
			msg = fmt.Sprintf("%s {%s}", msg, strings.Join(attrStrs, ", "))
		}
	}

	logLine.WriteString(msg)
	logLine.WriteString("\n")

	_, err := h.writer.Write([]byte(logLine.String()))
	return err
}

// formatLevel formats the log level with appropriate colors and highlighting
func (h *SpringBootHandler) formatLevel(level slog.Level) string {
	switch level {
	case slog.LevelError:
		return fmt.Sprintf("%s%s%-5s%s%s", ColorBold, ColorRed, "ERROR", ColorReset, "")
	case slog.LevelWarn:
		return fmt.Sprintf("%s%-5s%s", ColorYellow, "WARN", ColorReset)
	case slog.LevelInfo:
		return fmt.Sprintf("%s%-5s%s", ColorBlue, "INFO", ColorReset)
	case slog.LevelDebug:
		return fmt.Sprintf("%s%-5s%s", ColorCyan, "DEBUG", ColorReset)
	default:
		// Handle custom levels
		levelName := level.String()
		if len(levelName) > 5 {
			levelName = levelName[:5]
		}
		return fmt.Sprintf("%-5s", strings.ToUpper(levelName))
	}
}

// formatAttribute formats a single attribute for display
func (h *SpringBootHandler) formatAttribute(attr slog.Attr) string {
	key := attr.Key
	if h.group != "" {
		key = h.group + "." + key
	}

	return fmt.Sprintf("%s=%v", key, attr.Value)
}

// cleanFunctionName removes receiver types and package prefixes from function names
// Example: "(*application).healthcheckHandler" -> "healthcheckHandler"
func (h *SpringBootHandler) cleanFunctionName(funcName string) string {
	// Remove receiver type like "(*application)."
	if idx := strings.Index(funcName, ")."); idx != -1 {
		return funcName[idx+2:]
	}

	// Remove simple receiver type like "application."
	if idx := strings.LastIndex(funcName, "."); idx != -1 {
		return funcName[idx+1:]
	}

	return funcName
}

// getGoroutineID extracts the current goroutine ID for thread information
func (h *SpringBootHandler) getGoroutineID() string {
	buf := make([]byte, 64)
	buf = buf[:runtime.Stack(buf, false)]

	// Parse "goroutine 123 [running]:" to extract ID
	stack := string(buf)
	if strings.HasPrefix(stack, "goroutine ") {
		start := len("goroutine ")
		end := strings.Index(stack[start:], " ")
		if end > 0 {
			return stack[start : start+end]
		}
	}

	return "1"
}

// WithAttrs returns a new handler with the given attributes added
func (h *SpringBootHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newAttrs := make([]slog.Attr, len(h.attrs)+len(attrs))
	copy(newAttrs, h.attrs)
	copy(newAttrs[len(h.attrs):], attrs)

	return &SpringBootHandler{
		writer:      h.writer,
		level:       h.level,
		attrs:       newAttrs,
		group:       h.group,
		projectRoot: h.projectRoot,
	}
}

// WithGroup returns a new handler with the given group name
func (h *SpringBootHandler) WithGroup(name string) slog.Handler {
	group := name
	if h.group != "" {
		group = h.group + "." + name
	}

	return &SpringBootHandler{
		writer:      h.writer,
		level:       h.level,
		attrs:       h.attrs,
		group:       group,
		projectRoot: h.projectRoot,
	}
}
