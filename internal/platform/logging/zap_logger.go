package logging

import (
	"context"
	"fmt"
	"math"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// zapLogger wraps zap logger to implement Logger interface
type zapLogger struct {
	logger        *zap.Logger
	sugaredLogger *zap.SugaredLogger
}

// Compile-time interface check
var _ Logger = (*zapLogger)(nil)

// Structured logging methods
func (l *zapLogger) Debug(msg string, fields ...Field) {
	l.logger.Debug(msg, toZapFields(fields)...)
}

func (l *zapLogger) Info(msg string, fields ...Field) {
	l.logger.Info(msg, toZapFields(fields)...)
}

func (l *zapLogger) Warn(msg string, fields ...Field) {
	l.logger.Warn(msg, toZapFields(fields)...)
}

func (l *zapLogger) Error(msg string, fields ...Field) {
	l.logger.Error(msg, toZapFields(fields)...)
}

func (l *zapLogger) Fatal(msg string, fields ...Field) {
	l.logger.Fatal(msg, toZapFields(fields)...)
}

func (l *zapLogger) Panic(msg string, fields ...Field) {
	l.logger.Panic(msg, toZapFields(fields)...)
}

// Sugared logging methods
func (l *zapLogger) Debugf(template string, args ...any) {
	l.sugaredLogger.Debugf(template, args...)
}

func (l *zapLogger) Infof(template string, args ...any) {
	l.sugaredLogger.Infof(template, args...)
}

func (l *zapLogger) Warnf(template string, args ...any) {
	l.sugaredLogger.Warnf(template, args...)
}

func (l *zapLogger) Errorf(template string, args ...any) {
	l.sugaredLogger.Errorf(template, args...)
}

func (l *zapLogger) Fatalf(template string, args ...any) {
	l.sugaredLogger.Fatalf(template, args...)
}

func (l *zapLogger) Panicf(template string, args ...any) {
	l.sugaredLogger.Panicf(template, args...)
}

// Context-aware structured logging
func (l *zapLogger) DebugContext(_ context.Context, msg string, fields ...Field) {
	l.logger.Debug(msg, toZapFields(fields)...)
}

func (l *zapLogger) InfoContext(_ context.Context, msg string, fields ...Field) {
	l.logger.Info(msg, toZapFields(fields)...)
}

func (l *zapLogger) WarnContext(_ context.Context, msg string, fields ...Field) {
	l.logger.Warn(msg, toZapFields(fields)...)
}

func (l *zapLogger) ErrorContext(_ context.Context, msg string, fields ...Field) {
	l.logger.Error(msg, toZapFields(fields)...)
}

func (l *zapLogger) FatalContext(_ context.Context, msg string, fields ...Field) {
	l.logger.Fatal(msg, toZapFields(fields)...)
}

func (l *zapLogger) PanicContext(_ context.Context, msg string, fields ...Field) {
	l.logger.Panic(msg, toZapFields(fields)...)
}

// Context-aware sugared logging
func (l *zapLogger) DebugfContext(_ context.Context, template string, args ...any) {
	l.sugaredLogger.Debugf(template, args...)
}

func (l *zapLogger) InfofContext(_ context.Context, template string, args ...any) {
	l.sugaredLogger.Infof(template, args...)
}

func (l *zapLogger) WarnfContext(_ context.Context, template string, args ...any) {
	l.sugaredLogger.Warnf(template, args...)
}

func (l *zapLogger) ErrorfContext(_ context.Context, template string, args ...any) {
	l.sugaredLogger.Errorf(template, args...)
}

func (l *zapLogger) FatalfContext(_ context.Context, template string, args ...any) {
	l.sugaredLogger.Fatalf(template, args...)
}

func (l *zapLogger) PanicfContext(_ context.Context, template string, args ...any) {
	l.sugaredLogger.Panicf(template, args...)
}

// Logger manipulation
func (l *zapLogger) With(fields ...Field) Logger {
	zapFields := toZapFields(fields)
	newLogger := l.logger.With(zapFields...)

	return &zapLogger{
		logger:        newLogger,
		sugaredLogger: newLogger.Sugar(),
	}
}

func (l *zapLogger) Named(name string) Logger {
	return &zapLogger{
		logger:        l.logger.Named(name),
		sugaredLogger: l.sugaredLogger.Named(name),
	}
}

func (l *zapLogger) Sync() error {
	return fmt.Errorf("failed to sync logger: %w", l.logger.Sync())
}

// toZapField converts our Field to zap.Field
//
//nolint:gocyclo,revive,funlen // because we need to handle all possible types
func toZapField(f Field) zap.Field {
	switch f.Type {
	case SkipType:
		return zap.Skip()
	case BoolType:
		return zap.Bool(f.Key, f.Integer == 1)
	case Int64Type:
		return zap.Int64(f.Key, f.Integer)
	case Int32Type:
		return zap.Int32(f.Key, int32(f.Integer))
	case Int16Type:
		return zap.Int16(f.Key, int16(f.Integer))
	case Int8Type:
		return zap.Int8(f.Key, int8(f.Integer))
	case Uint64Type:
		return zap.Uint64(f.Key, uint64(f.Integer))
	case Uint32Type:
		return zap.Uint32(f.Key, uint32(f.Integer))
	case Uint16Type:
		return zap.Uint16(f.Key, uint16(f.Integer))
	case Uint8Type:
		return zap.Uint8(f.Key, uint8(f.Integer))
	case UintptrType:
		return zap.Uintptr(f.Key, uintptr(f.Integer))
	case Float64Type:
		return zap.Float64(f.Key, math.Float64frombits(uint64(f.Integer)))
	case Float32Type:
		return zap.Float32(f.Key, math.Float32frombits(uint32(f.Integer)))
	case Complex64Type:
		return zap.Complex64(f.Key, f.Interface.(complex64))
	case Complex128Type:
		return zap.Complex128(f.Key, f.Interface.(complex128))
	case StringType:
		return zap.String(f.Key, f.String)
	case BinaryType:
		return zap.Binary(f.Key, f.Interface.([]byte))
	case ByteStringType:
		return zap.ByteString(f.Key, f.Interface.([]byte))
	case DurationType:
		return zap.Duration(f.Key, time.Duration(f.Integer))
	case TimeType:
		if f.Interface != nil {
			return zap.Time(f.Key, f.Interface.(time.Time))
		}

		return zap.Skip()
	case ErrorType:
		if f.Interface != nil {
			return zap.NamedError(f.Key, f.Interface.(error))
		}

		return zap.Skip()
	case ReflectType:
		return zap.Reflect(f.Key, f.Interface)
	case NamespaceType:
		return zap.Namespace(f.Key)
	case StringerType:
		return zap.Stringer(f.Key, f.Interface.(fmt.Stringer))
	case ObjectMarshalerType:
		// Convert our ObjectMarshaler to zapcore.ObjectMarshaler
		if om, ok := f.Interface.(ObjectMarshaler); ok {
			return zap.Object(f.Key, objectMarshalerAdapter{om})
		}
		// Fallback for direct zapcore.ObjectMarshaler
		return zap.Object(f.Key, f.Interface.(zapcore.ObjectMarshaler))
	case InlineMarshalerType:
		// Convert our ObjectMarshaler to zapcore.ObjectMarshaler
		if om, ok := f.Interface.(ObjectMarshaler); ok {
			return zap.Inline(objectMarshalerAdapter{om})
		}
		// Fallback for direct zapcore.ObjectMarshaler
		return zap.Inline(f.Interface.(zapcore.ObjectMarshaler))
	case ArrayMarshalerType:
		// Convert our ArrayMarshaler to zapcore.ArrayMarshaler
		if am, ok := f.Interface.(ArrayMarshaler); ok {
			return zap.Array(f.Key, arrayMarshalerAdapter{am})
		}
		// Fallback for direct zapcore.ArrayMarshaler
		return zap.Array(f.Key, f.Interface.(zapcore.ArrayMarshaler))
	// Array types
	case BoolsType:
		return zap.Bools(f.Key, f.Interface.([]bool))
	case IntsType:
		return zap.Ints(f.Key, f.Interface.([]int))
	case Int64sType:
		return zap.Int64s(f.Key, f.Interface.([]int64))
	case Int32sType:
		return zap.Int32s(f.Key, f.Interface.([]int32))
	case Int16sType:
		return zap.Int16s(f.Key, f.Interface.([]int16))
	case Int8sType:
		return zap.Int8s(f.Key, f.Interface.([]int8))
	case UintsType:
		return zap.Uints(f.Key, f.Interface.([]uint))
	case Uint64sType:
		return zap.Uint64s(f.Key, f.Interface.([]uint64))
	case Uint32sType:
		return zap.Uint32s(f.Key, f.Interface.([]uint32))
	case Uint16sType:
		return zap.Uint16s(f.Key, f.Interface.([]uint16))
	case Uint8sType:
		return zap.Uint8s(f.Key, f.Interface.([]uint8))
	case UintptrsType:
		return zap.Uintptrs(f.Key, f.Interface.([]uintptr))
	case Float64sType:
		return zap.Float64s(f.Key, f.Interface.([]float64))
	case Float32sType:
		return zap.Float32s(f.Key, f.Interface.([]float32))
	case Complex64sType:
		return zap.Complex64s(f.Key, f.Interface.([]complex64))
	case Complex128sType:
		return zap.Complex128s(f.Key, f.Interface.([]complex128))
	case DurationsType:
		return zap.Durations(f.Key, f.Interface.([]time.Duration))
	case StringsType:
		return zap.Strings(f.Key, f.Interface.([]string))
	case TimesType:
		return zap.Times(f.Key, f.Interface.([]time.Time))
	case ErrorsType:
		return zap.Errors(f.Key, f.Interface.([]error))
	default:
		return zap.Skip()
	}
}

// toZapFields converts slice of our Fields to zap.Fields
func toZapFields(fields []Field) []zap.Field {
	zapFields := make([]zap.Field, len(fields))
	for i, f := range fields {
		zapFields[i] = toZapField(f)
	}

	return zapFields
}

// objectMarshalerAdapter adapts our ObjectMarshaler to zapcore.ObjectMarshaler
type objectMarshalerAdapter struct {
	om ObjectMarshaler
}

func (a objectMarshalerAdapter) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	return fmt.Errorf("failed to marshal object field %w", a.om.MarshalLogObject(objectEncoderAdapter{enc}))
}

// objectEncoderAdapter adapts zapcore.ObjectEncoder to our ObjectEncoder
type objectEncoderAdapter struct {
	enc zapcore.ObjectEncoder
}

func (a objectEncoderAdapter) AddString(key, val string) {
	a.enc.AddString(key, val)
}

func (a objectEncoderAdapter) AddInt64(key string, val int64) {
	a.enc.AddInt64(key, val)
}

func (a objectEncoderAdapter) AddInt(key string, val int) {
	a.enc.AddInt(key, val)
}

func (a objectEncoderAdapter) AddBool(key string, val bool) {
	a.enc.AddBool(key, val)
}

func (a objectEncoderAdapter) AddFloat64(key string, val float64) {
	a.enc.AddFloat64(key, val)
}

func (a objectEncoderAdapter) AddDuration(key string, val time.Duration) {
	a.enc.AddDuration(key, val)
}

func (a objectEncoderAdapter) AddTime(key string, val time.Time) {
	a.enc.AddTime(key, val)
}

func (a objectEncoderAdapter) AddObject(key string, val ObjectMarshaler) error {
	return fmt.Errorf("failed to add object field %w", a.enc.AddObject(key, objectMarshalerAdapter{val}))
}

func (a objectEncoderAdapter) AddArray(key string, val ArrayMarshaler) error {
	return fmt.Errorf("failed to add array field %w", a.enc.AddArray(key, arrayMarshalerAdapter{val}))
}

// arrayMarshalerAdapter adapts our ArrayMarshaler to zapcore.ArrayMarshaler
type arrayMarshalerAdapter struct {
	am ArrayMarshaler
}

func (a arrayMarshalerAdapter) MarshalLogArray(enc zapcore.ArrayEncoder) error {
	return fmt.Errorf("failed to marshal array field %w", a.am.MarshalLogArray(arrayEncoderAdapter{enc}))
}

// arrayEncoderAdapter adapts zapcore.ArrayEncoder to our ArrayEncoder
type arrayEncoderAdapter struct {
	enc zapcore.ArrayEncoder
}

func (a arrayEncoderAdapter) AppendString(val string) {
	a.enc.AppendString(val)
}

func (a arrayEncoderAdapter) AppendInt64(val int64) {
	a.enc.AppendInt64(val)
}

func (a arrayEncoderAdapter) AppendInt(val int) {
	a.enc.AppendInt(val)
}

func (a arrayEncoderAdapter) AppendBool(val bool) {
	a.enc.AppendBool(val)
}

func (a arrayEncoderAdapter) AppendFloat64(val float64) {
	a.enc.AppendFloat64(val)
}

func (a arrayEncoderAdapter) AppendDuration(val time.Duration) {
	a.enc.AppendDuration(val)
}

func (a arrayEncoderAdapter) AppendTime(val time.Time) {
	a.enc.AppendTime(val)
}

func (a arrayEncoderAdapter) AppendObject(val ObjectMarshaler) error {
	return fmt.Errorf("failed to append object field %w", a.enc.AppendObject(objectMarshalerAdapter{val}))
}
