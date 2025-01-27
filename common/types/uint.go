// Copyright 2018 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package types

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/common/types/traits"

	anypb "google.golang.org/protobuf/types/known/anypb"
	structpb "google.golang.org/protobuf/types/known/structpb"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

// Uint type implementation which supports comparison and math operators.
type Uint uint64

var (
	// UintType singleton.
	UintType = NewTypeValue("uint",
		traits.AdderType,
		traits.ComparerType,
		traits.DividerType,
		traits.ModderType,
		traits.MultiplierType,
		traits.SubtractorType)

	uint32WrapperType = reflect.TypeOf(&wrapperspb.UInt32Value{})

	uint64WrapperType = reflect.TypeOf(&wrapperspb.UInt64Value{})
)

// Uint constants
const (
	uintZero = Uint(0)
)

// Add implements traits.Adder.Add.
func (i Uint) Add(other ref.Val) ref.Val {
	otherUint, ok := other.(Uint)
	if !ok {
		return MaybeNoSuchOverloadErr(other)
	}
	if val, err := addUint64Checked(uint64(i), uint64(otherUint)); err != nil {
		return wrapErr(err)
	} else {
		return Uint(val)
	}
}

// Compare implements traits.Comparer.Compare.
func (i Uint) Compare(other ref.Val) ref.Val {
	otherUint, ok := other.(Uint)
	if !ok {
		return MaybeNoSuchOverloadErr(other)
	}
	if i < otherUint {
		return IntNegOne
	}
	if i > otherUint {
		return IntOne
	}
	return IntZero
}

// ConvertToNative implements ref.Val.ConvertToNative.
func (i Uint) ConvertToNative(typeDesc reflect.Type) (interface{}, error) {
	switch typeDesc.Kind() {
	case reflect.Uint, reflect.Uint32:
		v, err := uint64ToUint32Checked(uint64(i))
		if err != nil {
			return 0, err
		}
		return reflect.ValueOf(v).Convert(typeDesc).Interface(), nil
	case reflect.Uint64:
		return reflect.ValueOf(i).Convert(typeDesc).Interface(), nil
	case reflect.Ptr:
		switch typeDesc {
		case anyValueType:
			// Primitives must be wrapped before being set on an Any field.
			return anypb.New(wrapperspb.UInt64(uint64(i)))
		case jsonValueType:
			// JSON can accurately represent 32-bit uints as floating point values.
			if i.isJSONSafe() {
				return structpb.NewNumberValue(float64(i)), nil
			}
			// Proto3 to JSON conversion requires string-formatted uint64 values
			// since the conversion to floating point would result in truncation.
			return structpb.NewStringValue(strconv.FormatUint(uint64(i), 10)), nil
		case uint32WrapperType:
			// Convert the value to a wrapperspb.UInt32Value, error on overflow.
			v, err := uint64ToUint32Checked(uint64(i))
			if err != nil {
				return 0, err
			}
			return wrapperspb.UInt32(v), nil
		case uint64WrapperType:
			// Convert the value to a wrapperspb.UInt64Value.
			return wrapperspb.UInt64(uint64(i)), nil
		}
		switch typeDesc.Elem().Kind() {
		case reflect.Uint32:
			v, err := uint64ToUint32Checked(uint64(i))
			if err != nil {
				return 0, err
			}
			p := reflect.New(typeDesc.Elem())
			p.Elem().Set(reflect.ValueOf(v).Convert(typeDesc.Elem()))
			return p.Interface(), nil
		case reflect.Uint64:
			v := uint64(i)
			p := reflect.New(typeDesc.Elem())
			p.Elem().Set(reflect.ValueOf(v).Convert(typeDesc.Elem()))
			return p.Interface(), nil
		}
	case reflect.Interface:
		iv := i.Value()
		if reflect.TypeOf(iv).Implements(typeDesc) {
			return iv, nil
		}
		if reflect.TypeOf(i).Implements(typeDesc) {
			return i, nil
		}
	}
	return nil, fmt.Errorf("unsupported type conversion from 'uint' to %v", typeDesc)
}

// ConvertToType implements ref.Val.ConvertToType.
func (i Uint) ConvertToType(typeVal ref.Type) ref.Val {
	switch typeVal {
	case IntType:
		if v, err := uint64ToInt64Checked(uint64(i)); err != nil {
			return wrapErr(err)
		} else {
			return Int(v)
		}
	case UintType:
		return i
	case DoubleType:
		return Double(i)
	case StringType:
		return String(fmt.Sprintf("%d", uint64(i)))
	case TypeType:
		return UintType
	}
	return NewErr("type conversion error from '%s' to '%s'", UintType, typeVal)
}

// Divide implements traits.Divider.Divide.
func (i Uint) Divide(other ref.Val) ref.Val {
	otherUint, ok := other.(Uint)
	if !ok {
		return MaybeNoSuchOverloadErr(other)
	}
	if div, err := divideUint64Checked(uint64(i), uint64(otherUint)); err != nil {
		return wrapErr(err)
	} else {
		return Uint(div)
	}
}

// Equal implements ref.Val.Equal.
func (i Uint) Equal(other ref.Val) ref.Val {
	otherUint, ok := other.(Uint)
	if !ok {
		return MaybeNoSuchOverloadErr(other)
	}
	return Bool(i == otherUint)
}

// Modulo implements traits.Modder.Modulo.
func (i Uint) Modulo(other ref.Val) ref.Val {
	otherUint, ok := other.(Uint)
	if !ok {
		return MaybeNoSuchOverloadErr(other)
	}
	if mod, err := moduloUint64Checked(uint64(i), uint64(otherUint)); err != nil {
		return wrapErr(err)
	} else {
		return Uint(mod)
	}
}

// Multiply implements traits.Multiplier.Multiply.
func (i Uint) Multiply(other ref.Val) ref.Val {
	otherUint, ok := other.(Uint)
	if !ok {
		return MaybeNoSuchOverloadErr(other)
	}
	if val, err := multiplyUint64Checked(uint64(i), uint64(otherUint)); err != nil {
		return wrapErr(err)
	} else {
		return Uint(val)
	}
}

// Subtract implements traits.Subtractor.Subtract.
func (i Uint) Subtract(subtrahend ref.Val) ref.Val {
	subtraUint, ok := subtrahend.(Uint)
	if !ok {
		return MaybeNoSuchOverloadErr(subtrahend)
	}
	if val, err := subtractUint64Checked(uint64(i), uint64(subtraUint)); err != nil {
		return wrapErr(err)
	} else {
		return Uint(val)
	}
}

// Type implements ref.Val.Type.
func (i Uint) Type() ref.Type {
	return UintType
}

// Value implements ref.Val.Value.
func (i Uint) Value() interface{} {
	return uint64(i)
}

// isJSONSafe indicates whether the uint is safely representable as a floating point value in JSON.
func (i Uint) isJSONSafe() bool {
	return i <= maxIntJSON
}
