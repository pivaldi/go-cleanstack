package logging

import (
	"fmt"
	"math"
	"runtime"
	"strings"
	"time"
)

// Skip constructs a no-op field
func Skip() Field {
	return Field{Type: SkipType}
}

// Binary constructs a field that carries an opaque binary blob
func Binary(key string, val []byte) Field {
	return Field{Key: key, Type: BinaryType, Interface: val}
}

// Bool constructs a field that carries a bool
func Bool(key string, val bool) Field {
	var ival int64
	if val {
		ival = 1
	}

	return Field{Key: key, Type: BoolType, Integer: ival}
}

// ByteString constructs a field that carries UTF-8 encoded text as a []byte
func ByteString(key string, val []byte) Field {
	return Field{Key: key, Type: ByteStringType, Interface: val}
}

// Complex128 constructs a field that carries a complex number
func Complex128(key string, val complex128) Field {
	return Field{Key: key, Type: Complex128Type, Interface: val}
}

// Complex64 constructs a field that carries a complex number
func Complex64(key string, val complex64) Field {
	return Field{Key: key, Type: Complex64Type, Interface: val}
}

// Float64 constructs a field that carries a float64
func Float64(key string, val float64) Field {
	return Field{Key: key, Type: Float64Type, Integer: int64(math.Float64bits(val))}
}

// Float32 constructs a field that carries a float32
func Float32(key string, val float32) Field {
	return Field{Key: key, Type: Float32Type, Integer: int64(math.Float32bits(val))}
}

// Int constructs a field with the given key and value
func Int(key string, val int) Field {
	return Int64(key, int64(val))
}

// Int64 constructs a field with the given key and value
func Int64(key string, val int64) Field {
	return Field{Key: key, Type: Int64Type, Integer: val}
}

// Int32 constructs a field with the given key and value
func Int32(key string, val int32) Field {
	return Field{Key: key, Type: Int32Type, Integer: int64(val)}
}

// Int16 constructs a field with the given key and value
func Int16(key string, val int16) Field {
	return Field{Key: key, Type: Int16Type, Integer: int64(val)}
}

// Int8 constructs a field with the given key and value
func Int8(key string, val int8) Field {
	return Field{Key: key, Type: Int8Type, Integer: int64(val)}
}

// String constructs a field with the given key and value
func String(key, val string) Field {
	return Field{Key: key, Type: StringType, String: val}
}

// Uint constructs a field with the given key and value
func Uint(key string, val uint) Field {
	return Uint64(key, uint64(val))
}

// Uint64 constructs a field with the given key and value
func Uint64(key string, val uint64) Field {
	return Field{Key: key, Type: Uint64Type, Integer: int64(val)}
}

// Uint32 constructs a field with the given key and value
func Uint32(key string, val uint32) Field {
	return Field{Key: key, Type: Uint32Type, Integer: int64(val)}
}

// Uint16 constructs a field with the given key and value
func Uint16(key string, val uint16) Field {
	return Field{Key: key, Type: Uint16Type, Integer: int64(val)}
}

// Uint8 constructs a field with the given key and value
func Uint8(key string, val uint8) Field {
	return Field{Key: key, Type: Uint8Type, Integer: int64(val)}
}

// Uintptr constructs a field with the given key and value
func Uintptr(key string, val uintptr) Field {
	return Field{Key: key, Type: UintptrType, Integer: int64(val)}
}

// Duration constructs a field with the given key and value
func Duration(key string, val time.Duration) Field {
	return Field{Key: key, Type: DurationType, Integer: int64(val)}
}

// Time constructs a field with the given key and value
func Time(key string, val time.Time) Field {
	return Field{Key: key, Type: TimeType, Interface: val}
}

// Err is a shorthand for NamedError with key "error"
func Err(err error) Field {
	return NamedError("error", err)
}

// NamedError constructs a field that carries an error with custom key
func NamedError(key string, err error) Field {
	return Field{Key: key, Type: ErrorType, Interface: err}
}

// Namespace creates a named, isolated scope within the logger's context
func Namespace(key string) Field {
	return Field{Key: key, Type: NamespaceType}
}

// Stringer constructs a field with the given key and the output of the value's String method
func Stringer(key string, val fmt.Stringer) Field {
	return Field{Key: key, Type: StringerType, Interface: val}
}

// Reflect constructs a field with the given key and an arbitrary object
func Reflect(key string, val any) Field {
	return Field{Key: key, Type: ReflectType, Interface: val}
}

// Any takes a key and an arbitrary value and chooses the best way to represent them
func Any(key string, value any) Field {
	return Field{Key: key, Type: ReflectType, Interface: value}
}

// nilField constructs a field that represents a nil pointer value
func nilField(key string) Field {
	return Reflect(key, nil)
}

// Boolp constructs a field that carries a *bool. The returned Field will safely
// and explicitly represent `nil` when appropriate.
func Boolp(key string, val *bool) Field {
	if val == nil {
		return nilField(key)
	}

	return Bool(key, *val)
}

// Complex128p constructs a field that carries a *complex128. The returned Field will safely
// and explicitly represent `nil` when appropriate.
func Complex128p(key string, val *complex128) Field {
	if val == nil {
		return nilField(key)
	}

	return Complex128(key, *val)
}

// Complex64p constructs a field that carries a *complex64. The returned Field will safely
// and explicitly represent `nil` when appropriate.
func Complex64p(key string, val *complex64) Field {
	if val == nil {
		return nilField(key)
	}

	return Complex64(key, *val)
}

// Durationp constructs a field that carries a *time.Duration. The returned Field will safely
// and explicitly represent `nil` when appropriate.
func Durationp(key string, val *time.Duration) Field {
	if val == nil {
		return nilField(key)
	}

	return Duration(key, *val)
}

// Float64p constructs a field that carries a *float64. The returned Field will safely
// and explicitly represent `nil` when appropriate.
func Float64p(key string, val *float64) Field {
	if val == nil {
		return nilField(key)
	}

	return Float64(key, *val)
}

// Float32p constructs a field that carries a *float32. The returned Field will safely
// and explicitly represent `nil` when appropriate.
func Float32p(key string, val *float32) Field {
	if val == nil {
		return nilField(key)
	}

	return Float32(key, *val)
}

// Intp constructs a field that carries a *int. The returned Field will safely
// and explicitly represent `nil` when appropriate.
func Intp(key string, val *int) Field {
	if val == nil {
		return nilField(key)
	}

	return Int(key, *val)
}

// Int64p constructs a field that carries a *int64. The returned Field will safely
// and explicitly represent `nil` when appropriate.
func Int64p(key string, val *int64) Field {
	if val == nil {
		return nilField(key)
	}

	return Int64(key, *val)
}

// Int32p constructs a field that carries a *int32. The returned Field will safely
// and explicitly represent `nil` when appropriate.
func Int32p(key string, val *int32) Field {
	if val == nil {
		return nilField(key)
	}

	return Int32(key, *val)
}

// Int16p constructs a field that carries a *int16. The returned Field will safely
// and explicitly represent `nil` when appropriate.
func Int16p(key string, val *int16) Field {
	if val == nil {
		return nilField(key)
	}

	return Int16(key, *val)
}

// Int8p constructs a field that carries a *int8. The returned Field will safely
// and explicitly represent `nil` when appropriate.
func Int8p(key string, val *int8) Field {
	if val == nil {
		return nilField(key)
	}

	return Int8(key, *val)
}

// Stringp constructs a field that carries a *string. The returned Field will safely
// and explicitly represent `nil` when appropriate.
func Stringp(key string, val *string) Field {
	if val == nil {
		return nilField(key)
	}

	return String(key, *val)
}

// Timep constructs a field that carries a *time.Time. The returned Field will safely
// and explicitly represent `nil` when appropriate.
func Timep(key string, val *time.Time) Field {
	if val == nil {
		return nilField(key)
	}

	return Time(key, *val)
}

// Uintp constructs a field that carries a *uint. The returned Field will safely
// and explicitly represent `nil` when appropriate.
func Uintp(key string, val *uint) Field {
	if val == nil {
		return nilField(key)
	}

	return Uint(key, *val)
}

// Uint64p constructs a field that carries a *uint64. The returned Field will safely
// and explicitly represent `nil` when appropriate.
func Uint64p(key string, val *uint64) Field {
	if val == nil {
		return nilField(key)
	}

	return Uint64(key, *val)
}

// Uint32p constructs a field that carries a *uint32. The returned Field will safely
// and explicitly represent `nil` when appropriate.
func Uint32p(key string, val *uint32) Field {
	if val == nil {
		return nilField(key)
	}

	return Uint32(key, *val)
}

// Uint16p constructs a field that carries a *uint16. The returned Field will safely
// and explicitly represent `nil` when appropriate.
func Uint16p(key string, val *uint16) Field {
	if val == nil {
		return nilField(key)
	}

	return Uint16(key, *val)
}

// Uint8p constructs a field that carries a *uint8. The returned Field will safely
// and explicitly represent `nil` when appropriate.
func Uint8p(key string, val *uint8) Field {
	if val == nil {
		return nilField(key)
	}

	return Uint8(key, *val)
}

// Uintptrp constructs a field that carries a *uintptr. The returned Field will safely
// and explicitly represent `nil` when appropriate.
func Uintptrp(key string, val *uintptr) Field {
	if val == nil {
		return nilField(key)
	}

	return Uintptr(key, *val)
}

// Bools constructs a field that carries a slice of bools
func Bools(key string, vals []bool) Field {
	return Field{Key: key, Type: BoolsType, Interface: vals}
}

// Ints constructs a field that carries a slice of ints
func Ints(key string, vals []int) Field {
	return Field{Key: key, Type: IntsType, Interface: vals}
}

// Int64s constructs a field that carries a slice of int64s
func Int64s(key string, vals []int64) Field {
	return Field{Key: key, Type: Int64sType, Interface: vals}
}

// Int32s constructs a field that carries a slice of int32s
func Int32s(key string, vals []int32) Field {
	return Field{Key: key, Type: Int32sType, Interface: vals}
}

// Int16s constructs a field that carries a slice of int16s
func Int16s(key string, vals []int16) Field {
	return Field{Key: key, Type: Int16sType, Interface: vals}
}

// Int8s constructs a field that carries a slice of int8s
func Int8s(key string, vals []int8) Field {
	return Field{Key: key, Type: Int8sType, Interface: vals}
}

// Uints constructs a field that carries a slice of uints
func Uints(key string, vals []uint) Field {
	return Field{Key: key, Type: UintsType, Interface: vals}
}

// Uint64s constructs a field that carries a slice of uint64s
func Uint64s(key string, vals []uint64) Field {
	return Field{Key: key, Type: Uint64sType, Interface: vals}
}

// Uint32s constructs a field that carries a slice of uint32s
func Uint32s(key string, vals []uint32) Field {
	return Field{Key: key, Type: Uint32sType, Interface: vals}
}

// Uint16s constructs a field that carries a slice of uint16s
func Uint16s(key string, vals []uint16) Field {
	return Field{Key: key, Type: Uint16sType, Interface: vals}
}

// Uint8s constructs a field that carries a slice of uint8s
func Uint8s(key string, vals []uint8) Field {
	return Field{Key: key, Type: Uint8sType, Interface: vals}
}

// Uintptrs constructs a field that carries a slice of uintptrs
func Uintptrs(key string, vals []uintptr) Field {
	return Field{Key: key, Type: UintptrsType, Interface: vals}
}

// Float64s constructs a field that carries a slice of float64s
func Float64s(key string, vals []float64) Field {
	return Field{Key: key, Type: Float64sType, Interface: vals}
}

// Float32s constructs a field that carries a slice of float32s
func Float32s(key string, vals []float32) Field {
	return Field{Key: key, Type: Float32sType, Interface: vals}
}

// Complex128s constructs a field that carries a slice of complex128s
func Complex128s(key string, vals []complex128) Field {
	return Field{Key: key, Type: Complex128sType, Interface: vals}
}

// Complex64s constructs a field that carries a slice of complex64s
func Complex64s(key string, vals []complex64) Field {
	return Field{Key: key, Type: Complex64sType, Interface: vals}
}

// Durations constructs a field that carries a slice of time.Durations
func Durations(key string, vals []time.Duration) Field {
	return Field{Key: key, Type: DurationsType, Interface: vals}
}

// Strings constructs a field that carries a slice of strings
func Strings(key string, vals []string) Field {
	return Field{Key: key, Type: StringsType, Interface: vals}
}

// Times constructs a field that carries a slice of time.Times
func Times(key string, vals []time.Time) Field {
	return Field{Key: key, Type: TimesType, Interface: vals}
}

// Errors constructs a field that carries a slice of errors
func Errors(key string, vals []error) Field {
	return Field{Key: key, Type: ErrorsType, Interface: vals}
}

// Stack constructs a field that stores a stacktrace of the current goroutine
// under provided key. Keep in mind that taking a stacktrace is eager and
// expensive (relatively speaking); this function both makes an allocation and
// takes about two microseconds.
func Stack(key string) Field {
	return StackSkip(key, 1) // Skip Stack function itself
}

// StackSkip constructs a field similarly to Stack, but also skips the given
// number of frames from the top of the stacktrace.
func StackSkip(key string, skip int) Field {
	// We store the skip count in the Integer field and let the encoder
	// handle the actual stacktrace generation
	return Field{Key: key, Type: StringerType, Interface: stacktrace(skip + 1)}
}

// stacktrace is an internal type that implements fmt.Stringer
// to lazily capture stack traces
type stacktrace int

func (s stacktrace) String() string {
	return captureStack(int(s))
}

// captureStack captures the current goroutine's stack trace, skipping the specified number of frames
func captureStack(skip int) string {
	const depth = 32
	var pcs [depth]uintptr
	// +2 to skip captureStack itself and its caller (stacktrace.String)
	n := runtime.Callers(skip+2, pcs[:])

	if n == 0 {
		return ""
	}

	frames := runtime.CallersFrames(pcs[:n])
	var buf strings.Builder

	for {
		frame, more := frames.Next()

		// Format similar to runtime.Stack output
		fmt.Fprintf(&buf, "%s\n\t%s:%d\n", frame.Function, frame.File, frame.Line)

		if !more {
			break
		}
	}

	return buf.String()
}

// Object constructs a field with the given key and ObjectMarshaler
func Object(key string, val ObjectMarshaler) Field {
	return Field{Key: key, Type: ObjectMarshalerType, Interface: val}
}

// Inline constructs a Field that is similar to Object, but it
// will add the elements of the provided ObjectMarshaler to the
// current namespace
func Inline(val ObjectMarshaler) Field {
	return Field{Type: InlineMarshalerType, Interface: val}
}

// Dict constructs a field containing the provided key-value pairs
func Dict(key string, fields ...Field) Field {
	return Field{Key: key, Type: ObjectMarshalerType, Interface: dictObject(fields)}
}

// DictObject constructs an ObjectMarshaler with the given list of fields
func DictObject(fields ...Field) ObjectMarshaler {
	return dictObject(fields)
}

// dictObject is an internal type that implements ObjectMarshaler
type dictObject []Field

//nolint:gocyclo,revive // because we need to handle all possible types
func (d dictObject) MarshalLogObject(enc ObjectEncoder) error {
	for _, f := range d {
		// For each field, we need to add it to the encoder
		// This is a simplified implementation - the actual encoding
		// will be done by the zap adapter
		switch f.Type {
		case SkipType:
			// Intentionally skip this field
			continue
		case StringType:
			enc.AddString(f.Key, f.String)
		case Int64Type, Int32Type, Int16Type, Int8Type:
			enc.AddInt64(f.Key, f.Integer)
		case UintptrType, Uint64Type, Uint32Type, Uint16Type, Uint8Type:
			enc.AddInt64(f.Key, f.Integer)
		case BoolType:
			enc.AddBool(f.Key, f.Integer == 1)
		case Float64Type:
			enc.AddFloat64(f.Key, math.Float64frombits(uint64(f.Integer)))
		case Float32Type:
			enc.AddFloat64(f.Key, float64(math.Float32frombits(uint32(f.Integer))))
		case DurationType:
			enc.AddDuration(f.Key, time.Duration(f.Integer))
		case TimeType:
			if t, ok := f.Interface.(time.Time); ok {
				enc.AddTime(f.Key, t)
			}
		case ObjectMarshalerType:
			if om, ok := f.Interface.(ObjectMarshaler); ok {
				if err := enc.AddObject(f.Key, om); err != nil {
					return fmt.Errorf("failed to add object field %s: %w", f.Key, err)
				}
			}
		case ArrayMarshalerType:
			if am, ok := f.Interface.(ArrayMarshaler); ok {
				if err := enc.AddArray(f.Key, am); err != nil {
					return fmt.Errorf("failed to add array field %s: %w", f.Key, err)
				}
			}
		case ErrorType:
			if err, ok := f.Interface.(error); ok {
				enc.AddString(f.Key, err.Error())
			}
		case StringerType:
			if stringer, ok := f.Interface.(fmt.Stringer); ok {
				enc.AddString(f.Key, stringer.String())
			}
		case BinaryType, ByteStringType:
			if bytes, ok := f.Interface.([]byte); ok {
				enc.AddString(f.Key, string(bytes))
			}
		default:
			// Handle remaining types: Complex*, ReflectType, array types, etc.
			// Convert to string representation as a fallback
			enc.AddString(f.Key, fmt.Sprintf("%+v", f.Interface))
		}
	}

	return nil
}
