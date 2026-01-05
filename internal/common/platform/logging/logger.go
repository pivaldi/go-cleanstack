package logging

import (
	"context"
	"time"
)

// Field represents a structured logging field with complete independence from underlying logger
type Field struct {
	Key       string
	Type      FieldType
	Integer   int64  // For int types, bools, and encoded floats
	String    string // For string types
	Interface any    // For complex types (objects, arrays, etc.)
}

// FieldType identifies the type of field for encoding
type FieldType uint8

const (
	UnknownType FieldType = iota
	SkipType
	BoolType
	Int64Type
	Int32Type
	Int16Type
	Int8Type
	Uint64Type
	Uint32Type
	Uint16Type
	Uint8Type
	UintptrType
	Float64Type
	Float32Type
	Complex64Type
	Complex128Type
	StringType
	BinaryType
	ByteStringType
	DurationType
	TimeType
	ErrorType
	ReflectType
	NamespaceType
	StringerType
	ObjectMarshalerType
	InlineMarshalerType
	ArrayMarshalerType
	// Array types
	BoolsType
	IntsType
	Int64sType
	Int32sType
	Int16sType
	Int8sType
	UintsType
	Uint64sType
	Uint32sType
	Uint16sType
	Uint8sType
	UintptrsType
	Float64sType
	Float32sType
	Complex64sType
	Complex128sType
	DurationsType
	StringsType
	TimesType
	ErrorsType
)

// ObjectMarshaler allows custom types to control logging representation
type ObjectMarshaler interface {
	MarshalLogObject(enc ObjectEncoder) error
}

// ObjectEncoder provides methods for encoding object fields
type ObjectEncoder interface {
	AddString(key, val string)
	AddInt64(key string, val int64)
	AddInt(key string, val int)
	AddBool(key string, val bool)
	AddFloat64(key string, val float64)
	AddDuration(key string, val time.Duration)
	AddTime(key string, val time.Time)
	AddObject(key string, val ObjectMarshaler) error
	AddArray(key string, val ArrayMarshaler) error
}

// ArrayMarshaler allows arrays to control their logging representation
type ArrayMarshaler interface {
	MarshalLogArray(enc ArrayEncoder) error
}

// ArrayEncoder provides methods for encoding array elements
type ArrayEncoder interface {
	AppendString(val string)
	AppendInt64(val int64)
	AppendInt(val int)
	AppendBool(val bool)
	AppendFloat64(val float64)
	AppendDuration(val time.Duration)
	AppendTime(val time.Time)
	AppendObject(val ObjectMarshaler) error
}

// Logger is the interface for structured logging, independent of implementation
type Logger interface {
	// Structured logging methods
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
	Panic(msg string, fields ...Field)

	// Sugared logging methods (printf-style)
	Debugf(template string, args ...any)
	Infof(template string, args ...any)
	Warnf(template string, args ...any)
	Errorf(template string, args ...any)
	Fatalf(template string, args ...any)
	Panicf(template string, args ...any)

	// Context-aware structured logging
	DebugContext(ctx context.Context, msg string, fields ...Field)
	InfoContext(ctx context.Context, msg string, fields ...Field)
	WarnContext(ctx context.Context, msg string, fields ...Field)
	ErrorContext(ctx context.Context, msg string, fields ...Field)
	FatalContext(ctx context.Context, msg string, fields ...Field)
	PanicContext(ctx context.Context, msg string, fields ...Field)

	// Context-aware sugared logging
	DebugfContext(ctx context.Context, template string, args ...any)
	InfofContext(ctx context.Context, template string, args ...any)
	WarnfContext(ctx context.Context, template string, args ...any)
	ErrorfContext(ctx context.Context, template string, args ...any)
	FatalfContext(ctx context.Context, template string, args ...any)
	PanicfContext(ctx context.Context, template string, args ...any)

	// Logger manipulation
	With(fields ...Field) Logger
	Named(name string) Logger
	Sync() error
}
