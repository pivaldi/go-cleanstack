package logging

import (
	"errors"
	"testing"
	"time"
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
	//nolint:S1040 // Testing purpose
	_, ok := om.(ObjectMarshaler)
	if !ok {
		t.Error("expected DictObject to implement ObjectMarshaler")
	}
}
