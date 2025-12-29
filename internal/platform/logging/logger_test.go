package logging

import (
	"errors"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestField_StringType(t *testing.T) {
	f := Field{
		Key:    "name",
		Type:   StringType,
		String: "value",
	}

	if f.Key != "name" {
		t.Errorf("expected key 'name', got %q", f.Key)
	}
	if f.Type != StringType {
		t.Errorf("expected StringType, got %v", f.Type)
	}
	if f.String != "value" {
		t.Errorf("expected string 'value', got %q", f.String)
	}
}

func TestField_Int64Type(t *testing.T) {
	f := Field{
		Key:     "count",
		Type:    Int64Type,
		Integer: 42,
	}

	if f.Key != "count" {
		t.Errorf("expected key 'count', got %q", f.Key)
	}
	if f.Type != Int64Type {
		t.Errorf("expected Int64Type, got %v", f.Type)
	}
	if f.Integer != 42 {
		t.Errorf("expected integer 42, got %d", f.Integer)
	}
}

func TestZapLogger_ImplementsInterface(t *testing.T) {
	var _ Logger = (*zapLogger)(nil)
}

func TestToZapField_String(t *testing.T) {
	f := Field{
		Key:    "name",
		Type:   StringType,
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
	f := Field{
		Key:     "count",
		Type:    Int64Type,
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
	f := Field{
		Key:     "enabled",
		Type:    BoolType,
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

func TestString(t *testing.T) {
	f := String("name", "value")

	if f.Key != "name" {
		t.Errorf("expected key 'name', got %q", f.Key)
	}
	if f.Type != StringType {
		t.Errorf("expected StringType, got %v", f.Type)
	}
	if f.String != "value" {
		t.Errorf("expected string 'value', got %q", f.String)
	}
}

func TestBool(t *testing.T) {
	fTrue := Bool("enabled", true)
	if fTrue.Integer != 1 {
		t.Errorf("expected integer 1 for true, got %d", fTrue.Integer)
	}

	fFalse := Bool("disabled", false)
	if fFalse.Integer != 0 {
		t.Errorf("expected integer 0 for false, got %d", fFalse.Integer)
	}
}

func TestInt(t *testing.T) {
	f := Int("count", 42)

	if f.Type != Int64Type {
		t.Errorf("expected Int64Type, got %v", f.Type)
	}
	if f.Integer != 42 {
		t.Errorf("expected integer 42, got %d", f.Integer)
	}
}

// Pointer variant tests

func TestStringp(t *testing.T) {
	t.Run("non-nil", func(t *testing.T) {
		val := "test"
		f := Stringp("key", &val)
		if f.Type != StringType {
			t.Errorf("expected StringType, got %v", f.Type)
		}
		if f.String != "test" {
			t.Errorf("expected string 'test', got %q", f.String)
		}
	})

	t.Run("nil", func(t *testing.T) {
		f := Stringp("key", nil)
		if f.Type != ReflectType {
			t.Errorf("expected ReflectType for nil, got %v", f.Type)
		}
		if f.Interface != nil {
			t.Errorf("expected nil interface, got %v", f.Interface)
		}
	})
}

func TestBoolp(t *testing.T) {
	t.Run("non-nil true", func(t *testing.T) {
		val := true
		f := Boolp("key", &val)
		if f.Type != BoolType {
			t.Errorf("expected BoolType, got %v", f.Type)
		}
		if f.Integer != 1 {
			t.Errorf("expected integer 1, got %d", f.Integer)
		}
	})

	t.Run("non-nil false", func(t *testing.T) {
		val := false
		f := Boolp("key", &val)
		if f.Type != BoolType {
			t.Errorf("expected BoolType, got %v", f.Type)
		}
		if f.Integer != 0 {
			t.Errorf("expected integer 0, got %d", f.Integer)
		}
	})

	t.Run("nil", func(t *testing.T) {
		f := Boolp("key", nil)
		if f.Type != ReflectType {
			t.Errorf("expected ReflectType for nil, got %v", f.Type)
		}
	})
}

func TestIntp(t *testing.T) {
	t.Run("non-nil", func(t *testing.T) {
		val := 42
		f := Intp("key", &val)
		if f.Type != Int64Type {
			t.Errorf("expected Int64Type, got %v", f.Type)
		}
		if f.Integer != 42 {
			t.Errorf("expected integer 42, got %d", f.Integer)
		}
	})

	t.Run("nil", func(t *testing.T) {
		f := Intp("key", nil)
		if f.Type != ReflectType {
			t.Errorf("expected ReflectType for nil, got %v", f.Type)
		}
	})
}

func TestInt64p(t *testing.T) {
	t.Run("non-nil", func(t *testing.T) {
		val := int64(42)
		f := Int64p("key", &val)
		if f.Type != Int64Type {
			t.Errorf("expected Int64Type, got %v", f.Type)
		}
		if f.Integer != 42 {
			t.Errorf("expected integer 42, got %d", f.Integer)
		}
	})

	t.Run("nil", func(t *testing.T) {
		f := Int64p("key", nil)
		if f.Type != ReflectType {
			t.Errorf("expected ReflectType for nil, got %v", f.Type)
		}
	})
}

func TestInt32p(t *testing.T) {
	t.Run("non-nil", func(t *testing.T) {
		val := int32(42)
		f := Int32p("key", &val)
		if f.Type != Int32Type {
			t.Errorf("expected Int32Type, got %v", f.Type)
		}
		if f.Integer != 42 {
			t.Errorf("expected integer 42, got %d", f.Integer)
		}
	})

	t.Run("nil", func(t *testing.T) {
		f := Int32p("key", nil)
		if f.Type != ReflectType {
			t.Errorf("expected ReflectType for nil, got %v", f.Type)
		}
	})
}

func TestInt16p(t *testing.T) {
	t.Run("non-nil", func(t *testing.T) {
		val := int16(42)
		f := Int16p("key", &val)
		if f.Type != Int16Type {
			t.Errorf("expected Int16Type, got %v", f.Type)
		}
		if f.Integer != 42 {
			t.Errorf("expected integer 42, got %d", f.Integer)
		}
	})

	t.Run("nil", func(t *testing.T) {
		f := Int16p("key", nil)
		if f.Type != ReflectType {
			t.Errorf("expected ReflectType for nil, got %v", f.Type)
		}
	})
}

func TestInt8p(t *testing.T) {
	t.Run("non-nil", func(t *testing.T) {
		val := int8(42)
		f := Int8p("key", &val)
		if f.Type != Int8Type {
			t.Errorf("expected Int8Type, got %v", f.Type)
		}
		if f.Integer != 42 {
			t.Errorf("expected integer 42, got %d", f.Integer)
		}
	})

	t.Run("nil", func(t *testing.T) {
		f := Int8p("key", nil)
		if f.Type != ReflectType {
			t.Errorf("expected ReflectType for nil, got %v", f.Type)
		}
	})
}

func TestUintp(t *testing.T) {
	t.Run("non-nil", func(t *testing.T) {
		val := uint(42)
		f := Uintp("key", &val)
		if f.Type != Uint64Type {
			t.Errorf("expected Uint64Type, got %v", f.Type)
		}
		if f.Integer != 42 {
			t.Errorf("expected integer 42, got %d", f.Integer)
		}
	})

	t.Run("nil", func(t *testing.T) {
		f := Uintp("key", nil)
		if f.Type != ReflectType {
			t.Errorf("expected ReflectType for nil, got %v", f.Type)
		}
	})
}

func TestUint64p(t *testing.T) {
	t.Run("non-nil", func(t *testing.T) {
		val := uint64(42)
		f := Uint64p("key", &val)
		if f.Type != Uint64Type {
			t.Errorf("expected Uint64Type, got %v", f.Type)
		}
		if f.Integer != 42 {
			t.Errorf("expected integer 42, got %d", f.Integer)
		}
	})

	t.Run("nil", func(t *testing.T) {
		f := Uint64p("key", nil)
		if f.Type != ReflectType {
			t.Errorf("expected ReflectType for nil, got %v", f.Type)
		}
	})
}

func TestUint32p(t *testing.T) {
	t.Run("non-nil", func(t *testing.T) {
		val := uint32(42)
		f := Uint32p("key", &val)
		if f.Type != Uint32Type {
			t.Errorf("expected Uint32Type, got %v", f.Type)
		}
		if f.Integer != 42 {
			t.Errorf("expected integer 42, got %d", f.Integer)
		}
	})

	t.Run("nil", func(t *testing.T) {
		f := Uint32p("key", nil)
		if f.Type != ReflectType {
			t.Errorf("expected ReflectType for nil, got %v", f.Type)
		}
	})
}

func TestUint16p(t *testing.T) {
	t.Run("non-nil", func(t *testing.T) {
		val := uint16(42)
		f := Uint16p("key", &val)
		if f.Type != Uint16Type {
			t.Errorf("expected Uint16Type, got %v", f.Type)
		}
		if f.Integer != 42 {
			t.Errorf("expected integer 42, got %d", f.Integer)
		}
	})

	t.Run("nil", func(t *testing.T) {
		f := Uint16p("key", nil)
		if f.Type != ReflectType {
			t.Errorf("expected ReflectType for nil, got %v", f.Type)
		}
	})
}

func TestUint8p(t *testing.T) {
	t.Run("non-nil", func(t *testing.T) {
		val := uint8(42)
		f := Uint8p("key", &val)
		if f.Type != Uint8Type {
			t.Errorf("expected Uint8Type, got %v", f.Type)
		}
		if f.Integer != 42 {
			t.Errorf("expected integer 42, got %d", f.Integer)
		}
	})

	t.Run("nil", func(t *testing.T) {
		f := Uint8p("key", nil)
		if f.Type != ReflectType {
			t.Errorf("expected ReflectType for nil, got %v", f.Type)
		}
	})
}

func TestUintptrp(t *testing.T) {
	t.Run("non-nil", func(t *testing.T) {
		val := uintptr(42)
		f := Uintptrp("key", &val)
		if f.Type != UintptrType {
			t.Errorf("expected UintptrType, got %v", f.Type)
		}
		if f.Integer != 42 {
			t.Errorf("expected integer 42, got %d", f.Integer)
		}
	})

	t.Run("nil", func(t *testing.T) {
		f := Uintptrp("key", nil)
		if f.Type != ReflectType {
			t.Errorf("expected ReflectType for nil, got %v", f.Type)
		}
	})
}

func TestFloat64p(t *testing.T) {
	t.Run("non-nil", func(t *testing.T) {
		val := 3.14
		f := Float64p("key", &val)
		if f.Type != Float64Type {
			t.Errorf("expected Float64Type, got %v", f.Type)
		}
	})

	t.Run("nil", func(t *testing.T) {
		f := Float64p("key", nil)
		if f.Type != ReflectType {
			t.Errorf("expected ReflectType for nil, got %v", f.Type)
		}
	})
}

func TestFloat32p(t *testing.T) {
	t.Run("non-nil", func(t *testing.T) {
		val := float32(3.14)
		f := Float32p("key", &val)
		if f.Type != Float32Type {
			t.Errorf("expected Float32Type, got %v", f.Type)
		}
	})

	t.Run("nil", func(t *testing.T) {
		f := Float32p("key", nil)
		if f.Type != ReflectType {
			t.Errorf("expected ReflectType for nil, got %v", f.Type)
		}
	})
}

func TestComplex128p(t *testing.T) {
	t.Run("non-nil", func(t *testing.T) {
		val := complex(1.0, 2.0)
		f := Complex128p("key", &val)
		if f.Type != Complex128Type {
			t.Errorf("expected Complex128Type, got %v", f.Type)
		}
	})

	t.Run("nil", func(t *testing.T) {
		f := Complex128p("key", nil)
		if f.Type != ReflectType {
			t.Errorf("expected ReflectType for nil, got %v", f.Type)
		}
	})
}

func TestComplex64p(t *testing.T) {
	t.Run("non-nil", func(t *testing.T) {
		val := complex64(complex(1.0, 2.0))
		f := Complex64p("key", &val)
		if f.Type != Complex64Type {
			t.Errorf("expected Complex64Type, got %v", f.Type)
		}
	})

	t.Run("nil", func(t *testing.T) {
		f := Complex64p("key", nil)
		if f.Type != ReflectType {
			t.Errorf("expected ReflectType for nil, got %v", f.Type)
		}
	})
}

func TestTimep(t *testing.T) {
	t.Run("non-nil", func(t *testing.T) {
		now := time.Now()
		f := Timep("key", &now)
		if f.Type != TimeType {
			t.Errorf("expected TimeType, got %v", f.Type)
		}
	})

	t.Run("nil", func(t *testing.T) {
		f := Timep("key", nil)
		if f.Type != ReflectType {
			t.Errorf("expected ReflectType for nil, got %v", f.Type)
		}
	})
}

func TestDurationp(t *testing.T) {
	t.Run("non-nil", func(t *testing.T) {
		dur := time.Duration(42)
		f := Durationp("key", &dur)
		if f.Type != DurationType {
			t.Errorf("expected DurationType, got %v", f.Type)
		}
		if f.Integer != 42 {
			t.Errorf("expected integer 42, got %d", f.Integer)
		}
	})

	t.Run("nil", func(t *testing.T) {
		f := Durationp("key", nil)
		if f.Type != ReflectType {
			t.Errorf("expected ReflectType for nil, got %v", f.Type)
		}
	})
}

// Array/Slice constructor tests

func TestBools(t *testing.T) {
	t.Run("non-empty", func(t *testing.T) {
		vals := []bool{true, false, true}
		f := Bools("flags", vals)
		if f.Type != BoolsType {
			t.Errorf("expected BoolsType, got %v", f.Type)
		}
		if f.Key != "flags" {
			t.Errorf("expected key 'flags', got %q", f.Key)
		}
		slice, ok := f.Interface.([]bool)
		if !ok {
			t.Errorf("expected []bool in Interface, got %T", f.Interface)
		}
		if len(slice) != 3 {
			t.Errorf("expected slice length 3, got %d", len(slice))
		}
	})

	t.Run("empty", func(t *testing.T) {
		f := Bools("flags", []bool{})
		if f.Type != BoolsType {
			t.Errorf("expected BoolsType, got %v", f.Type)
		}
	})

	t.Run("nil", func(t *testing.T) {
		f := Bools("flags", nil)
		if f.Type != BoolsType {
			t.Errorf("expected BoolsType, got %v", f.Type)
		}
		// Note: nil slice stored in interface{} is not == nil,
		// it's a typed nil ([]bool(nil)). This matches zap's behavior.
		slice, ok := f.Interface.([]bool)
		if !ok {
			t.Errorf("expected []bool type in Interface, got %T", f.Interface)
		}
		if slice != nil {
			t.Errorf("expected nil slice, got %v", slice)
		}
	})
}

func TestInts(t *testing.T) {
	t.Run("non-empty", func(t *testing.T) {
		vals := []int{1, 2, 3}
		f := Ints("numbers", vals)
		if f.Type != IntsType {
			t.Errorf("expected IntsType, got %v", f.Type)
		}
		slice, ok := f.Interface.([]int)
		if !ok {
			t.Errorf("expected []int in Interface, got %T", f.Interface)
		}
		if len(slice) != 3 {
			t.Errorf("expected slice length 3, got %d", len(slice))
		}
	})

	t.Run("empty", func(t *testing.T) {
		f := Ints("numbers", []int{})
		if f.Type != IntsType {
			t.Errorf("expected IntsType, got %v", f.Type)
		}
	})

	t.Run("nil", func(t *testing.T) {
		f := Ints("numbers", nil)
		if f.Type != IntsType {
			t.Errorf("expected IntsType, got %v", f.Type)
		}
	})
}

func TestInt64s(t *testing.T) {
	t.Run("non-empty", func(t *testing.T) {
		vals := []int64{1, 2, 3}
		f := Int64s("numbers", vals)
		if f.Type != Int64sType {
			t.Errorf("expected Int64sType, got %v", f.Type)
		}
		slice, ok := f.Interface.([]int64)
		if !ok {
			t.Errorf("expected []int64 in Interface, got %T", f.Interface)
		}
		if len(slice) != 3 {
			t.Errorf("expected slice length 3, got %d", len(slice))
		}
	})

	t.Run("empty", func(t *testing.T) {
		f := Int64s("numbers", []int64{})
		if f.Type != Int64sType {
			t.Errorf("expected Int64sType, got %v", f.Type)
		}
	})

	t.Run("nil", func(t *testing.T) {
		f := Int64s("numbers", nil)
		if f.Type != Int64sType {
			t.Errorf("expected Int64sType, got %v", f.Type)
		}
	})
}

func TestInt32s(t *testing.T) {
	t.Run("non-empty", func(t *testing.T) {
		vals := []int32{1, 2, 3}
		f := Int32s("numbers", vals)
		if f.Type != Int32sType {
			t.Errorf("expected Int32sType, got %v", f.Type)
		}
		slice, ok := f.Interface.([]int32)
		if !ok {
			t.Errorf("expected []int32 in Interface, got %T", f.Interface)
		}
		if len(slice) != 3 {
			t.Errorf("expected slice length 3, got %d", len(slice))
		}
	})

	t.Run("empty", func(t *testing.T) {
		f := Int32s("numbers", []int32{})
		if f.Type != Int32sType {
			t.Errorf("expected Int32sType, got %v", f.Type)
		}
	})

	t.Run("nil", func(t *testing.T) {
		f := Int32s("numbers", nil)
		if f.Type != Int32sType {
			t.Errorf("expected Int32sType, got %v", f.Type)
		}
	})
}

func TestInt16s(t *testing.T) {
	t.Run("non-empty", func(t *testing.T) {
		vals := []int16{1, 2, 3}
		f := Int16s("numbers", vals)
		if f.Type != Int16sType {
			t.Errorf("expected Int16sType, got %v", f.Type)
		}
		slice, ok := f.Interface.([]int16)
		if !ok {
			t.Errorf("expected []int16 in Interface, got %T", f.Interface)
		}
		if len(slice) != 3 {
			t.Errorf("expected slice length 3, got %d", len(slice))
		}
	})

	t.Run("empty", func(t *testing.T) {
		f := Int16s("numbers", []int16{})
		if f.Type != Int16sType {
			t.Errorf("expected Int16sType, got %v", f.Type)
		}
	})

	t.Run("nil", func(t *testing.T) {
		f := Int16s("numbers", nil)
		if f.Type != Int16sType {
			t.Errorf("expected Int16sType, got %v", f.Type)
		}
	})
}

func TestInt8s(t *testing.T) {
	t.Run("non-empty", func(t *testing.T) {
		vals := []int8{1, 2, 3}
		f := Int8s("numbers", vals)
		if f.Type != Int8sType {
			t.Errorf("expected Int8sType, got %v", f.Type)
		}
		slice, ok := f.Interface.([]int8)
		if !ok {
			t.Errorf("expected []int8 in Interface, got %T", f.Interface)
		}
		if len(slice) != 3 {
			t.Errorf("expected slice length 3, got %d", len(slice))
		}
	})

	t.Run("empty", func(t *testing.T) {
		f := Int8s("numbers", []int8{})
		if f.Type != Int8sType {
			t.Errorf("expected Int8sType, got %v", f.Type)
		}
	})

	t.Run("nil", func(t *testing.T) {
		f := Int8s("numbers", nil)
		if f.Type != Int8sType {
			t.Errorf("expected Int8sType, got %v", f.Type)
		}
	})
}

func TestUints(t *testing.T) {
	t.Run("non-empty", func(t *testing.T) {
		vals := []uint{1, 2, 3}
		f := Uints("numbers", vals)
		if f.Type != UintsType {
			t.Errorf("expected UintsType, got %v", f.Type)
		}
		slice, ok := f.Interface.([]uint)
		if !ok {
			t.Errorf("expected []uint in Interface, got %T", f.Interface)
		}
		if len(slice) != 3 {
			t.Errorf("expected slice length 3, got %d", len(slice))
		}
	})

	t.Run("empty", func(t *testing.T) {
		f := Uints("numbers", []uint{})
		if f.Type != UintsType {
			t.Errorf("expected UintsType, got %v", f.Type)
		}
	})

	t.Run("nil", func(t *testing.T) {
		f := Uints("numbers", nil)
		if f.Type != UintsType {
			t.Errorf("expected UintsType, got %v", f.Type)
		}
	})
}

func TestUint64s(t *testing.T) {
	t.Run("non-empty", func(t *testing.T) {
		vals := []uint64{1, 2, 3}
		f := Uint64s("numbers", vals)
		if f.Type != Uint64sType {
			t.Errorf("expected Uint64sType, got %v", f.Type)
		}
		slice, ok := f.Interface.([]uint64)
		if !ok {
			t.Errorf("expected []uint64 in Interface, got %T", f.Interface)
		}
		if len(slice) != 3 {
			t.Errorf("expected slice length 3, got %d", len(slice))
		}
	})

	t.Run("empty", func(t *testing.T) {
		f := Uint64s("numbers", []uint64{})
		if f.Type != Uint64sType {
			t.Errorf("expected Uint64sType, got %v", f.Type)
		}
	})

	t.Run("nil", func(t *testing.T) {
		f := Uint64s("numbers", nil)
		if f.Type != Uint64sType {
			t.Errorf("expected Uint64sType, got %v", f.Type)
		}
	})
}

func TestUint32s(t *testing.T) {
	t.Run("non-empty", func(t *testing.T) {
		vals := []uint32{1, 2, 3}
		f := Uint32s("numbers", vals)
		if f.Type != Uint32sType {
			t.Errorf("expected Uint32sType, got %v", f.Type)
		}
		slice, ok := f.Interface.([]uint32)
		if !ok {
			t.Errorf("expected []uint32 in Interface, got %T", f.Interface)
		}
		if len(slice) != 3 {
			t.Errorf("expected slice length 3, got %d", len(slice))
		}
	})

	t.Run("empty", func(t *testing.T) {
		f := Uint32s("numbers", []uint32{})
		if f.Type != Uint32sType {
			t.Errorf("expected Uint32sType, got %v", f.Type)
		}
	})

	t.Run("nil", func(t *testing.T) {
		f := Uint32s("numbers", nil)
		if f.Type != Uint32sType {
			t.Errorf("expected Uint32sType, got %v", f.Type)
		}
	})
}

func TestUint16s(t *testing.T) {
	t.Run("non-empty", func(t *testing.T) {
		vals := []uint16{1, 2, 3}
		f := Uint16s("numbers", vals)
		if f.Type != Uint16sType {
			t.Errorf("expected Uint16sType, got %v", f.Type)
		}
		slice, ok := f.Interface.([]uint16)
		if !ok {
			t.Errorf("expected []uint16 in Interface, got %T", f.Interface)
		}
		if len(slice) != 3 {
			t.Errorf("expected slice length 3, got %d", len(slice))
		}
	})

	t.Run("empty", func(t *testing.T) {
		f := Uint16s("numbers", []uint16{})
		if f.Type != Uint16sType {
			t.Errorf("expected Uint16sType, got %v", f.Type)
		}
	})

	t.Run("nil", func(t *testing.T) {
		f := Uint16s("numbers", nil)
		if f.Type != Uint16sType {
			t.Errorf("expected Uint16sType, got %v", f.Type)
		}
	})
}

func TestUint8s(t *testing.T) {
	t.Run("non-empty", func(t *testing.T) {
		vals := []uint8{1, 2, 3}
		f := Uint8s("numbers", vals)
		if f.Type != Uint8sType {
			t.Errorf("expected Uint8sType, got %v", f.Type)
		}
		slice, ok := f.Interface.([]uint8)
		if !ok {
			t.Errorf("expected []uint8 in Interface, got %T", f.Interface)
		}
		if len(slice) != 3 {
			t.Errorf("expected slice length 3, got %d", len(slice))
		}
	})

	t.Run("empty", func(t *testing.T) {
		f := Uint8s("numbers", []uint8{})
		if f.Type != Uint8sType {
			t.Errorf("expected Uint8sType, got %v", f.Type)
		}
	})

	t.Run("nil", func(t *testing.T) {
		f := Uint8s("numbers", nil)
		if f.Type != Uint8sType {
			t.Errorf("expected Uint8sType, got %v", f.Type)
		}
	})
}

func TestUintptrs(t *testing.T) {
	t.Run("non-empty", func(t *testing.T) {
		vals := []uintptr{1, 2, 3}
		f := Uintptrs("pointers", vals)
		if f.Type != UintptrsType {
			t.Errorf("expected UintptrsType, got %v", f.Type)
		}
		slice, ok := f.Interface.([]uintptr)
		if !ok {
			t.Errorf("expected []uintptr in Interface, got %T", f.Interface)
		}
		if len(slice) != 3 {
			t.Errorf("expected slice length 3, got %d", len(slice))
		}
	})

	t.Run("empty", func(t *testing.T) {
		f := Uintptrs("pointers", []uintptr{})
		if f.Type != UintptrsType {
			t.Errorf("expected UintptrsType, got %v", f.Type)
		}
	})

	t.Run("nil", func(t *testing.T) {
		f := Uintptrs("pointers", nil)
		if f.Type != UintptrsType {
			t.Errorf("expected UintptrsType, got %v", f.Type)
		}
	})
}

func TestFloat64s(t *testing.T) {
	t.Run("non-empty", func(t *testing.T) {
		vals := []float64{1.1, 2.2, 3.3}
		f := Float64s("values", vals)
		if f.Type != Float64sType {
			t.Errorf("expected Float64sType, got %v", f.Type)
		}
		slice, ok := f.Interface.([]float64)
		if !ok {
			t.Errorf("expected []float64 in Interface, got %T", f.Interface)
		}
		if len(slice) != 3 {
			t.Errorf("expected slice length 3, got %d", len(slice))
		}
	})

	t.Run("empty", func(t *testing.T) {
		f := Float64s("values", []float64{})
		if f.Type != Float64sType {
			t.Errorf("expected Float64sType, got %v", f.Type)
		}
	})

	t.Run("nil", func(t *testing.T) {
		f := Float64s("values", nil)
		if f.Type != Float64sType {
			t.Errorf("expected Float64sType, got %v", f.Type)
		}
	})
}

func TestFloat32s(t *testing.T) {
	t.Run("non-empty", func(t *testing.T) {
		vals := []float32{1.1, 2.2, 3.3}
		f := Float32s("values", vals)
		if f.Type != Float32sType {
			t.Errorf("expected Float32sType, got %v", f.Type)
		}
		slice, ok := f.Interface.([]float32)
		if !ok {
			t.Errorf("expected []float32 in Interface, got %T", f.Interface)
		}
		if len(slice) != 3 {
			t.Errorf("expected slice length 3, got %d", len(slice))
		}
	})

	t.Run("empty", func(t *testing.T) {
		f := Float32s("values", []float32{})
		if f.Type != Float32sType {
			t.Errorf("expected Float32sType, got %v", f.Type)
		}
	})

	t.Run("nil", func(t *testing.T) {
		f := Float32s("values", nil)
		if f.Type != Float32sType {
			t.Errorf("expected Float32sType, got %v", f.Type)
		}
	})
}

func TestComplex128s(t *testing.T) {
	t.Run("non-empty", func(t *testing.T) {
		vals := []complex128{complex(1, 2), complex(3, 4), complex(5, 6)}
		f := Complex128s("values", vals)
		if f.Type != Complex128sType {
			t.Errorf("expected Complex128sType, got %v", f.Type)
		}
		slice, ok := f.Interface.([]complex128)
		if !ok {
			t.Errorf("expected []complex128 in Interface, got %T", f.Interface)
		}
		if len(slice) != 3 {
			t.Errorf("expected slice length 3, got %d", len(slice))
		}
	})

	t.Run("empty", func(t *testing.T) {
		f := Complex128s("values", []complex128{})
		if f.Type != Complex128sType {
			t.Errorf("expected Complex128sType, got %v", f.Type)
		}
	})

	t.Run("nil", func(t *testing.T) {
		f := Complex128s("values", nil)
		if f.Type != Complex128sType {
			t.Errorf("expected Complex128sType, got %v", f.Type)
		}
	})
}

func TestComplex64s(t *testing.T) {
	t.Run("non-empty", func(t *testing.T) {
		vals := []complex64{complex(1, 2), complex(3, 4), complex(5, 6)}
		f := Complex64s("values", vals)
		if f.Type != Complex64sType {
			t.Errorf("expected Complex64sType, got %v", f.Type)
		}
		slice, ok := f.Interface.([]complex64)
		if !ok {
			t.Errorf("expected []complex64 in Interface, got %T", f.Interface)
		}
		if len(slice) != 3 {
			t.Errorf("expected slice length 3, got %d", len(slice))
		}
	})

	t.Run("empty", func(t *testing.T) {
		f := Complex64s("values", []complex64{})
		if f.Type != Complex64sType {
			t.Errorf("expected Complex64sType, got %v", f.Type)
		}
	})

	t.Run("nil", func(t *testing.T) {
		f := Complex64s("values", nil)
		if f.Type != Complex64sType {
			t.Errorf("expected Complex64sType, got %v", f.Type)
		}
	})
}

func TestDurations(t *testing.T) {
	t.Run("non-empty", func(t *testing.T) {
		vals := []time.Duration{time.Second, time.Minute, time.Hour}
		f := Durations("times", vals)
		if f.Type != DurationsType {
			t.Errorf("expected DurationsType, got %v", f.Type)
		}
		slice, ok := f.Interface.([]time.Duration)
		if !ok {
			t.Errorf("expected []time.Duration in Interface, got %T", f.Interface)
		}
		if len(slice) != 3 {
			t.Errorf("expected slice length 3, got %d", len(slice))
		}
	})

	t.Run("empty", func(t *testing.T) {
		f := Durations("times", []time.Duration{})
		if f.Type != DurationsType {
			t.Errorf("expected DurationsType, got %v", f.Type)
		}
	})

	t.Run("nil", func(t *testing.T) {
		f := Durations("times", nil)
		if f.Type != DurationsType {
			t.Errorf("expected DurationsType, got %v", f.Type)
		}
	})
}

func TestStrings(t *testing.T) {
	t.Run("non-empty", func(t *testing.T) {
		vals := []string{"a", "b", "c"}
		f := Strings("names", vals)
		if f.Type != StringsType {
			t.Errorf("expected StringsType, got %v", f.Type)
		}
		slice, ok := f.Interface.([]string)
		if !ok {
			t.Errorf("expected []string in Interface, got %T", f.Interface)
		}
		if len(slice) != 3 {
			t.Errorf("expected slice length 3, got %d", len(slice))
		}
	})

	t.Run("empty", func(t *testing.T) {
		f := Strings("names", []string{})
		if f.Type != StringsType {
			t.Errorf("expected StringsType, got %v", f.Type)
		}
	})

	t.Run("nil", func(t *testing.T) {
		f := Strings("names", nil)
		if f.Type != StringsType {
			t.Errorf("expected StringsType, got %v", f.Type)
		}
	})
}

func TestTimes(t *testing.T) {
	t.Run("non-empty", func(t *testing.T) {
		now := time.Now()
		vals := []time.Time{now, now, now}
		f := Times("timestamps", vals)
		if f.Type != TimesType {
			t.Errorf("expected TimesType, got %v", f.Type)
		}
		slice, ok := f.Interface.([]time.Time)
		if !ok {
			t.Errorf("expected []time.Time in Interface, got %T", f.Interface)
		}
		if len(slice) != 3 {
			t.Errorf("expected slice length 3, got %d", len(slice))
		}
	})

	t.Run("empty", func(t *testing.T) {
		f := Times("timestamps", []time.Time{})
		if f.Type != TimesType {
			t.Errorf("expected TimesType, got %v", f.Type)
		}
	})

	t.Run("nil", func(t *testing.T) {
		f := Times("timestamps", nil)
		if f.Type != TimesType {
			t.Errorf("expected TimesType, got %v", f.Type)
		}
	})
}

func TestErrors(t *testing.T) {
	t.Run("non-empty", func(t *testing.T) {
		err1 := errors.New("error 1")
		err2 := errors.New("error 2")
		err3 := errors.New("error 3")
		vals := []error{err1, err2, err3}
		f := Errors("failures", vals)
		if f.Type != ErrorsType {
			t.Errorf("expected ErrorsType, got %v", f.Type)
		}
		slice, ok := f.Interface.([]error)
		if !ok {
			t.Errorf("expected []error in Interface, got %T", f.Interface)
		}
		if len(slice) != 3 {
			t.Errorf("expected slice length 3, got %d", len(slice))
		}
	})

	t.Run("empty", func(t *testing.T) {
		f := Errors("failures", []error{})
		if f.Type != ErrorsType {
			t.Errorf("expected ErrorsType, got %v", f.Type)
		}
	})

	t.Run("nil", func(t *testing.T) {
		f := Errors("failures", nil)
		if f.Type != ErrorsType {
			t.Errorf("expected ErrorsType, got %v", f.Type)
		}
	})
}

// Structured type tests

func TestStack(t *testing.T) {
	f := Stack("stack")
	if f.Key != "stack" {
		t.Errorf("expected key 'stack', got %q", f.Key)
	}
	if f.Type != StringerType {
		t.Errorf("expected StringerType, got %v", f.Type)
	}
	if f.Interface == nil {
		t.Error("expected non-nil Interface for stack trace")
	}
}

func TestStackSkip(t *testing.T) {
	f := StackSkip("stack", 2)
	if f.Key != "stack" {
		t.Errorf("expected key 'stack', got %q", f.Key)
	}
	if f.Type != StringerType {
		t.Errorf("expected StringerType, got %v", f.Type)
	}
	if f.Interface == nil {
		t.Error("expected non-nil Interface for stack trace")
	}
}

// mockObjectMarshaler is a test implementation of ObjectMarshaler
type mockObjectMarshaler struct {
	called bool
}

func (m *mockObjectMarshaler) MarshalLogObject(enc ObjectEncoder) error {
	m.called = true
	enc.AddString("test", "value")
	enc.AddInt("count", 42)
	return nil
}

func TestObject(t *testing.T) {
	mock := &mockObjectMarshaler{}
	f := Object("obj", mock)
	if f.Key != "obj" {
		t.Errorf("expected key 'obj', got %q", f.Key)
	}
	if f.Type != ObjectMarshalerType {
		t.Errorf("expected ObjectMarshalerType, got %v", f.Type)
	}
	if f.Interface == nil {
		t.Error("expected non-nil Interface")
	}
}

func TestInline(t *testing.T) {
	mock := &mockObjectMarshaler{}
	f := Inline(mock)
	if f.Key != "" {
		t.Errorf("expected empty key for inline, got %q", f.Key)
	}
	if f.Type != InlineMarshalerType {
		t.Errorf("expected InlineMarshalerType, got %v", f.Type)
	}
	if f.Interface == nil {
		t.Error("expected non-nil Interface")
	}
}

func TestDict(t *testing.T) {
	f := Dict("dict",
		String("name", "value"),
		Int("count", 42),
		Bool("enabled", true),
	)
	if f.Key != "dict" {
		t.Errorf("expected key 'dict', got %q", f.Key)
	}
	if f.Type != ObjectMarshalerType {
		t.Errorf("expected ObjectMarshalerType, got %v", f.Type)
	}
	if f.Interface == nil {
		t.Error("expected non-nil Interface")
	}

	// Verify the dictObject can be used as ObjectMarshaler
	om, ok := f.Interface.(ObjectMarshaler)
	if !ok {
		t.Error("expected Interface to implement ObjectMarshaler")
	}
	if om == nil {
		t.Error("expected non-nil ObjectMarshaler")
	}
}

func TestDictObject(t *testing.T) {
	om := DictObject(
		String("name", "value"),
		Int("count", 42),
	)
	if om == nil {
		t.Error("expected non-nil ObjectMarshaler")
	}

	// Verify it implements the interface
	_, ok := om.(ObjectMarshaler)
	if !ok {
		t.Error("expected DictObject to implement ObjectMarshaler")
	}
}

// Test adapter implementations
func TestObjectMarshalerAdapter(t *testing.T) {
	mock := &mockObjectMarshaler{}
	adapter := objectMarshalerAdapter{om: mock}

	// Create a zap logger to get a real ObjectEncoder
	logger, err := zap.NewDevelopment()
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
		var _ Logger = logger
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
		var _ Logger = logger
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
		var _ Logger = logger
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
		var _ Logger = logger
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
	var _ Logger = logger

	// Verify we can call methods without panicking
	logger.Info("test message", String("key", "value"))
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
	newLogger := logger.With(String("field", "value"))
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
		logger.Debug("debug message", String("key", "value"))
		logger.Info("info message", Int("count", 42))
		logger.Warn("warn message", Bool("flag", true))
		logger.Error("error message", Err(errors.New("test error")))

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

		childLogger := logger.With(String("service", "test"), Int("pid", 12345))
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
			With(String("service", "test")).
			With(Int("pid", 12345)).
			With(String("version", "1.0"))

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

		logger.Info("message", String("key", "value"))
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
			String("str", "value"),
			Int("int", 42),
			Int64("int64", int64(42)),
			Int32("int32", int32(42)),
			Int16("int16", int16(42)),
			Int8("int8", int8(42)),
			Uint("uint", uint(42)),
			Uint64("uint64", uint64(42)),
			Uint32("uint32", uint32(42)),
			Uint16("uint16", uint16(42)),
			Uint8("uint8", uint8(42)),
			Float64("float64", 3.14),
			Float32("float32", float32(3.14)),
			Bool("bool", true),
			Duration("dur", time.Second),
			Time("time", time.Now()),
			Err(errors.New("test")),
		)

		t.Log("All field types executed successfully")
	})

	t.Run("array field types", func(t *testing.T) {
		logger := Must(NewDevelopment("debug"))
		defer logger.Sync()

		logger.Info("array types",
			Strings("strs", []string{"a", "b"}),
			Ints("ints", []int{1, 2, 3}),
			Bools("bools", []bool{true, false}),
			Durations("durs", []time.Duration{time.Second, time.Minute}),
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
			Stringp("strp", &str),
			Intp("intp", &i),
			Boolp("boolp", &b),
			Float64p("floatp", &f),
			Stringp("nilp", nil), // nil pointer
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
			With(String("version", "1.0"), String("env", "test"))

		if serviceLogger == nil {
			t.Fatal("combined With/Named returned nil logger")
		}

		serviceLogger.Info("user created", String("userID", "123"))
		serviceLogger.Warn("user limit approaching", Int("count", 95))

		t.Log("Combined features executed successfully")
	})

	t.Run("complex logging scenario", func(t *testing.T) {
		logger := Must(NewProduction("info"))
		defer logger.Sync()

		// Simulate realistic logging scenario
		requestLogger := logger.
			Named("http").
			With(String("requestID", "req-123"), String("method", "POST"))

		requestLogger.Info("request started", String("path", "/api/users"))

		handlerLogger := requestLogger.Named("handler")
		handlerLogger.Info("processing request",
			String("userID", "user-456"),
			Int("bodySize", 1024),
		)

		handlerLogger.Warn("rate limit approaching",
			Int("remaining", 5),
			Duration("resetIn", 30*time.Second),
		)

		requestLogger.Info("request completed",
			Int("status", 200),
			Duration("duration", 150*time.Millisecond),
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
			Err(err1),
			Errors("additionalErrors", []error{err2}),
		)

		t.Log("Error logging executed successfully")
	})

	t.Run("nil error field", func(t *testing.T) {
		logger := Must(NewDevelopment("debug"))
		defer logger.Sync()

		// Should not panic with nil error
		logger.Info("message with nil error", Err(nil))

		t.Log("Nil error logging executed successfully")
	})
}
