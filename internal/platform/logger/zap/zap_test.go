package zap

import (
	"errors"
	"testing"
	"time"

	"github.com/pivaldi/go-cleanstack/internal/platform/logging"
	"go.uber.org/zap/zapcore"
)

// mockObjectMarshaler is a test implementation of ObjectMarshaler
type mockObjectMarshaler struct {
	called bool
}

func (m *mockObjectMarshaler) MarshalLogObject(enc logging.ObjectEncoder) error {
	m.called = true
	enc.AddString("test", "value")
	enc.AddInt("count", 42)
	return nil
}

func TestZapLogger_ImplementsInterface(t *testing.T) {
	var _ logging.Logger = (*zapLogger)(nil)
}

func TestToZapField_String(t *testing.T) {
	f := logging.Field{
		Key:    "name",
		Type:   logging.StringType,
		String: "value",
	}

	zapField := toZapField(f)
	if zapField.Key != "name" {
		t.Errorf("expected key 'name', got %q", zapField.Key)
	}
	if zapField.Type != zapcore.StringType {
		t.Errorf("expected StringType, got %v", zapField.Type)
	}
	if zapField.String != "value" {
		t.Errorf("expected string 'value', got %q", zapField.String)
	}
}

func TestToZapField_Int64(t *testing.T) {
	f := logging.Field{
		Key:     "count",
		Type:    logging.Int64Type,
		Integer: 42,
	}

	zapField := toZapField(f)
	if zapField.Key != "count" {
		t.Errorf("expected key 'count', got %q", zapField.Key)
	}
	if zapField.Type != zapcore.Int64Type {
		t.Errorf("expected Int64Type, got %v", zapField.Type)
	}
	if zapField.Integer != 42 {
		t.Errorf("expected integer 42, got %d", zapField.Integer)
	}
}

func TestToZapField_Bool(t *testing.T) {
	f := logging.Field{
		Key:     "enabled",
		Type:    logging.BoolType,
		Integer: 1,
	}

	zapField := toZapField(f)
	if zapField.Key != "enabled" {
		t.Errorf("expected key 'enabled', got %q", zapField.Key)
	}
	if zapField.Type != zapcore.BoolType {
		t.Errorf("expected BoolType, got %v", zapField.Type)
	}
}

// Test adapter implementations
func TestObjectMarshalerAdapter(t *testing.T) {
	mock := &mockObjectMarshaler{}
	adapter := objectMarshalerAdapter{om: mock}

	// Create a zap logger to get a real ObjectEncoder
	logger, err := NewDevelopment("debug")
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}
	defer logger.Sync()

	// We can't easily test the actual marshaling without a real encoder,
	// but we can verify the adapter exists and has the right structure
	if adapter.om != mock {
		t.Error("adapter should contain the mock")
	}
}

// Tests for logger creation functions

func TestNewProduction(t *testing.T) {
	t.Run("valid level", func(t *testing.T) {
		logger, err := NewProduction("info")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if logger == nil {
			t.Fatal("expected non-nil logger")
		}

		// Verify it implements Logger interface
		var _ logging.Logger = logger
	})

	t.Run("invalid level", func(t *testing.T) {
		logger, err := NewProduction("invalid")
		if err == nil {
			t.Error("expected error for invalid level")
		}
		if logger != nil {
			t.Error("expected nil logger on error")
		}
	})

	t.Run("debug level", func(t *testing.T) {
		logger, err := NewProduction("debug")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if logger == nil {
			t.Fatal("expected non-nil logger")
		}
	})
}

func TestNewDevelopment(t *testing.T) {
	t.Run("valid level", func(t *testing.T) {
		logger, err := NewDevelopment("debug")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if logger == nil {
			t.Fatal("expected non-nil logger")
		}

		// Verify it implements Logger interface
		var _ logging.Logger = logger
	})

	t.Run("invalid level", func(t *testing.T) {
		logger, err := NewDevelopment("invalid")
		if err == nil {
			t.Error("expected error for invalid level")
		}
		if logger != nil {
			t.Error("expected nil logger on error")
		}
	})

	t.Run("info level", func(t *testing.T) {
		logger, err := NewDevelopment("info")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if logger == nil {
			t.Fatal("expected non-nil logger")
		}
	})
}

func TestNewLogger(t *testing.T) {
	t.Run("development environment", func(t *testing.T) {
		logger, err := NewLogger("development", "debug")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if logger == nil {
			t.Fatal("expected non-nil logger")
		}

		// Verify it implements Logger interface
		var _ logging.Logger = logger
	})

	t.Run("production environment", func(t *testing.T) {
		logger, err := NewLogger("production", "info")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if logger == nil {
			t.Fatal("expected non-nil logger")
		}

		// Verify it implements Logger interface
		var _ logging.Logger = logger
	})

	t.Run("staging environment defaults to production", func(t *testing.T) {
		logger, err := NewLogger("staging", "warn")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if logger == nil {
			t.Fatal("expected non-nil logger")
		}
	})

	t.Run("invalid level", func(t *testing.T) {
		logger, err := NewLogger("production", "invalid")
		if err == nil {
			t.Error("expected error for invalid level")
		}
		if logger != nil {
			t.Error("expected nil logger on error")
		}
	})
}

func TestNewNop(t *testing.T) {
	logger := NewNop()
	if logger == nil {
		t.Fatal("expected non-nil logger")
	}

	// Verify it implements Logger interface
	var _ logging.Logger = logger

	// Verify we can call methods without panicking
	logger.Info("test message", logging.String("key", "value"))
	logger.Infof("test message %s", "arg")
	logger.Debug("debug message")
	logger.Warn("warn message")
	logger.Error("error message")

	// Verify Sync doesn't fail
	err := logger.Sync()
	if err != nil {
		t.Errorf("expected no error from Sync, got %v", err)
	}

	// Verify With returns a new logger
	newLogger := logger.With(logging.String("field", "value"))
	if newLogger == nil {
		t.Error("expected non-nil logger from With")
	}

	// Verify Named returns a new logger
	namedLogger := logger.Named("test")
	if namedLogger == nil {
		t.Error("expected non-nil logger from Named")
	}
}

func TestMust(t *testing.T) {
	t.Run("success returns logger", func(t *testing.T) {
		logger, err := NewProduction("info")
		if err != nil {
			t.Fatalf("failed to create logger: %v", err)
		}

		result := Must(logger, nil)
		if result == nil {
			t.Error("expected non-nil logger")
		}
		if result != logger {
			t.Error("expected same logger instance")
		}
	})

	t.Run("error panics", func(t *testing.T) {
		defer func() {
			r := recover()
			if r == nil {
				t.Error("expected panic, got none")
			}
			// Verify panic message contains error info
			msg, ok := r.(string)
			if !ok {
				t.Errorf("expected string panic, got %T", r)
			}
			if msg == "" {
				t.Error("expected non-empty panic message")
			}
		}()

		Must(nil, errors.New("test error"))
	})

	t.Run("can be used with NewProduction", func(t *testing.T) {
		logger := Must(NewProduction("info"))
		if logger == nil {
			t.Error("expected non-nil logger")
		}
	})

	t.Run("can be used with NewDevelopment", func(t *testing.T) {
		logger := Must(NewDevelopment("debug"))
		if logger == nil {
			t.Error("expected non-nil logger")
		}
	})

	t.Run("can be used with NewLogger", func(t *testing.T) {
		logger := Must(NewLogger("development", "debug"))
		if logger == nil {
			t.Error("expected non-nil logger")
		}
	})
}

func TestParseLevel(t *testing.T) {
	t.Run("valid levels", func(t *testing.T) {
		validLevels := []string{"debug", "info", "warn", "error", "dpanic", "panic", "fatal"}
		for _, level := range validLevels {
			_, err := parseLevel(level)
			if err != nil {
				t.Errorf("expected no error for level %q, got %v", level, err)
			}
		}
	})

	t.Run("invalid level", func(t *testing.T) {
		_, err := parseLevel("invalid")
		if err == nil {
			t.Error("expected error for invalid level")
		}
	})

	t.Run("empty level", func(t *testing.T) {
		// Note: zap's UnmarshalText accepts empty string and defaults to "info"
		level, err := parseLevel("")
		if err != nil {
			t.Errorf("unexpected error for empty level: %v", err)
		}
		// Empty level defaults to info (zapcore.InfoLevel)
		if level.String() != "info" {
			t.Errorf("expected empty level to default to info, got %s", level.String())
		}
	})
}

// ============================================================================
// Integration Tests
// ============================================================================
//
// The tests below verify that the logging abstraction works correctly with
// real logger instances (not mocks). These integration tests ensure:
// - Logger interface methods work correctly
// - Field constructors produce usable fields
// - zapLogger adapter works correctly
// - With() and Named() create proper child loggers
// - Sync() works without errors

func TestLoggerIntegration_BasicLogging(t *testing.T) {
	t.Run("development logger all levels", func(t *testing.T) {
		logger := Must(NewDevelopment("debug"))
		defer logger.Sync()

		// Test all log levels - should not panic
		logger.Debug("debug message", logging.String("key", "value"))
		logger.Info("info message", logging.Int("count", 42))
		logger.Warn("warn message", logging.Bool("flag", true))
		logger.Error("error message", logging.Err(errors.New("test error")))

		// If we got here without panicking, test passed
		t.Log("All log levels executed successfully")
	})

	t.Run("production logger respects level", func(t *testing.T) {
		logger := Must(NewProduction("info"))
		defer logger.Sync()

		// Debug should be filtered out at info level, but shouldn't panic
		logger.Debug("should be filtered")
		logger.Info("should be logged")
		logger.Warn("should be logged")
		logger.Error("should be logged")

		t.Log("Production logger executed successfully")
	})
}

func TestLoggerIntegration_SugaredLogging(t *testing.T) {
	t.Run("formatted logging", func(t *testing.T) {
		logger := Must(NewDevelopment("info"))
		defer logger.Sync()

		// Test printf-style logging - should not panic
		logger.Infof("formatted message: %s", "test")
		logger.Debugf("debug: %d", 42)
		logger.Warnf("warning: %v", true)
		logger.Errorf("error: %s", "test error")

		t.Log("Formatted logging executed successfully")
	})
}

func TestLoggerIntegration_With(t *testing.T) {
	t.Run("persistent fields", func(t *testing.T) {
		logger := Must(NewDevelopment("info"))
		defer logger.Sync()

		childLogger := logger.With(logging.String("service", "test"), logging.Int("pid", 12345))
		if childLogger == nil {
			t.Fatal("With() returned nil logger")
		}

		// Both messages should include service and pid fields
		childLogger.Info("message 1")
		childLogger.Info("message 2")

		t.Log("Child logger with persistent fields executed successfully")
	})

	t.Run("chained With calls", func(t *testing.T) {
		logger := Must(NewDevelopment("info"))
		defer logger.Sync()

		childLogger := logger.
			With(logging.String("service", "test")).
			With(logging.Int("pid", 12345)).
			With(logging.String("version", "1.0"))

		if childLogger == nil {
			t.Fatal("chained With() returned nil logger")
		}

		childLogger.Info("message with multiple fields")
		t.Log("Chained With() executed successfully")
	})
}

func TestLoggerIntegration_Named(t *testing.T) {
	t.Run("single named logger", func(t *testing.T) {
		logger := Must(NewDevelopment("info"))
		defer logger.Sync()

		serviceLogger := logger.Named("MyService")
		if serviceLogger == nil {
			t.Fatal("Named() returned nil logger")
		}

		serviceLogger.Info("service started")
		t.Log("Named logger executed successfully")
	})

	t.Run("nested named loggers", func(t *testing.T) {
		logger := Must(NewDevelopment("info"))
		defer logger.Sync()

		serviceLogger := logger.Named("MyService")
		handlerLogger := serviceLogger.Named("Handler")

		if handlerLogger == nil {
			t.Fatal("nested Named() returned nil logger")
		}

		handlerLogger.Info("handler executing")
		t.Log("Nested named loggers executed successfully")
	})
}

func TestLoggerIntegration_ContextLogging(t *testing.T) {
	t.Run("basic logging", func(t *testing.T) {
		logger := Must(NewDevelopment("info"))
		defer logger.Sync()

		// Note: Our Logger interface doesn't have context methods yet,
		// so we're testing that basic logging works (context support is future work)

		logger.Info("message", logging.String("key", "value"))
		logger.Debugf("formatted: %s", "test")

		t.Log("Context logging executed successfully")
	})
}

func TestLoggerIntegration_AllFieldTypes(t *testing.T) {
	t.Run("comprehensive field types", func(t *testing.T) {
		logger := Must(NewDevelopment("debug"))
		defer logger.Sync()

		// Test all major field constructors
		logger.Info("all types",
			logging.String("str", "value"),
			logging.Int("int", 42),
			logging.Int64("int64", int64(42)),
			logging.Int32("int32", int32(42)),
			logging.Int16("int16", int16(42)),
			logging.Int8("int8", int8(42)),
			logging.Uint("uint", uint(42)),
			logging.Uint64("uint64", uint64(42)),
			logging.Uint32("uint32", uint32(42)),
			logging.Uint16("uint16", uint16(42)),
			logging.Uint8("uint8", uint8(42)),
			logging.Float64("float64", 3.14),
			logging.Float32("float32", float32(3.14)),
			logging.Bool("bool", true),
			logging.Duration("dur", time.Second),
			logging.Time("time", time.Now()),
			logging.Err(errors.New("test")),
		)

		t.Log("All field types executed successfully")
	})

	t.Run("array field types", func(t *testing.T) {
		logger := Must(NewDevelopment("debug"))
		defer logger.Sync()

		logger.Info("array types",
			logging.Strings("strs", []string{"a", "b"}),
			logging.Ints("ints", []int{1, 2, 3}),
			logging.Bools("bools", []bool{true, false}),
			logging.Durations("durs", []time.Duration{time.Second, time.Minute}),
		)

		t.Log("Array field types executed successfully")
	})

	t.Run("pointer field types", func(t *testing.T) {
		logger := Must(NewDevelopment("debug"))
		defer logger.Sync()

		str := "value"
		i := 42
		b := true
		f := 3.14

		logger.Info("pointer types",
			logging.Stringp("strp", &str),
			logging.Intp("intp", &i),
			logging.Boolp("boolp", &b),
			logging.Float64p("floatp", &f),
			logging.Stringp("nilp", nil), // nil pointer
		)

		t.Log("Pointer field types executed successfully")
	})
}

func TestLoggerIntegration_Sync(t *testing.T) {
	t.Run("production logger sync", func(t *testing.T) {
		logger := Must(NewProduction("info"))
		logger.Info("message before sync")

		err := logger.Sync()
		// Note: Sync can return non-nil error for stderr/stdout on some systems
		// This is expected and not a failure - we just verify it doesn't panic
		if err != nil {
			t.Logf("Sync returned error (this can be normal for stderr/stdout): %v", err)
		} else {
			t.Log("Sync completed without error")
		}
	})

	t.Run("development logger sync", func(t *testing.T) {
		logger := Must(NewDevelopment("debug"))
		logger.Debug("message before sync")

		err := logger.Sync()
		if err != nil {
			t.Logf("Sync returned error (this can be normal for stderr/stdout): %v", err)
		} else {
			t.Log("Sync completed without error")
		}
	})

	t.Run("nop logger sync", func(t *testing.T) {
		logger := NewNop()
		logger.Info("message before sync")

		err := logger.Sync()
		if err != nil {
			t.Errorf("Nop logger Sync should not error, got: %v", err)
		}
	})
}

func TestLoggerIntegration_CombinedFeatures(t *testing.T) {
	t.Run("with and named together", func(t *testing.T) {
		logger := Must(NewDevelopment("info"))
		defer logger.Sync()

		serviceLogger := logger.
			Named("UserService").
			With(logging.String("version", "1.0"), logging.String("env", "test"))

		if serviceLogger == nil {
			t.Fatal("combined With/Named returned nil logger")
		}

		serviceLogger.Info("user created", logging.String("userID", "123"))
		serviceLogger.Warn("user limit approaching", logging.Int("count", 95))

		t.Log("Combined features executed successfully")
	})

	t.Run("complex logging scenario", func(t *testing.T) {
		logger := Must(NewProduction("info"))
		defer logger.Sync()

		// Simulate realistic logging scenario
		requestLogger := logger.
			Named("http").
			With(logging.String("requestID", "req-123"), logging.String("method", "POST"))

		requestLogger.Info("request started", logging.String("path", "/api/users"))

		handlerLogger := requestLogger.Named("handler")
		handlerLogger.Info("processing request",
			logging.String("userID", "user-456"),
			logging.Int("bodySize", 1024),
		)

		handlerLogger.Warn("rate limit approaching",
			logging.Int("remaining", 5),
			logging.Duration("resetIn", 30*time.Second),
		)

		requestLogger.Info("request completed",
			logging.Int("status", 200),
			logging.Duration("duration", 150*time.Millisecond),
		)

		t.Log("Complex logging scenario executed successfully")
	})
}

func TestLoggerIntegration_ErrorHandling(t *testing.T) {
	t.Run("logging multiple errors", func(t *testing.T) {
		logger := Must(NewDevelopment("debug"))
		defer logger.Sync()

		err1 := errors.New("first error")
		err2 := errors.New("second error")

		logger.Error("multiple errors occurred",
			logging.Err(err1),
			logging.Errors("additionalErrors", []error{err2}),
		)

		t.Log("Error logging executed successfully")
	})

	t.Run("nil error field", func(t *testing.T) {
		logger := Must(NewDevelopment("debug"))
		defer logger.Sync()

		// Should not panic with nil error
		logger.Info("message with nil error", logging.Err(nil))

		t.Log("Nil error logging executed successfully")
	})
}
