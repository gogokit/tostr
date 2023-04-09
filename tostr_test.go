package tostr

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"
	"unsafe"

	. "github.com/smartystreets/goconvey/convey"
)

// 测试普通结构体
func TestToString(t *testing.T) {
	Convey("ToString", t, func() {
		type CommonStruct struct {
			IntSlice1    []int
			IntSlice2    []int
			IntSlice3    []int
			ByteSlice    []byte
			StructSlice1 []CommonStruct
			StructSlice2 []CommonStruct
			StructSlice3 []CommonStruct
			IntPtr       ***int
			SlicePtr     *[]int
		}

		num := 1
		numP := &num
		numPP := &numP
		numPPP := &numPP
		cs := CommonStruct{
			IntPtr:    numPPP,
			IntSlice1: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			ByteSlice: []byte("this is byte slice!"),
		}
		cs.IntSlice3 = cs.IntSlice1[2:5]
		cs.SlicePtr = &cs.IntSlice1
		cs.StructSlice1 = []CommonStruct{cs, cs, cs}
		cs.StructSlice2 = cs.StructSlice1[1:]
		cs.StructSlice3 = []CommonStruct{}

		Convey("empty_config", func() {
			const loopCnt = 1000
			const expectStr = `{IntSlice1:$Obj1(0-9), IntSlice2:nil, IntSlice3:$Obj1(2-4), ByteSlice:$Obj2(0-18), StructSlice1:$Obj3(0-2), StructSlice2:$Obj3(1-2), StructSlice3:[], IntPtr:1, SlicePtr:<Obj4>$Obj1(0-9)}, {<Obj1>:[1(0), 2, 3(2), 4, 5(4), 6, 7, 8, 9, 10(9)], <Obj2>:[116(0), 104, 105, 115, 32, 105, 115, 32, 98, 121, 116, 101, 32, 115, 108, 105, 99, 101, 33(18)], <Obj3>:[{IntSlice1:$Obj1(0-9), IntSlice2:nil, IntSlice3:$Obj1(2-4), ByteSlice:$Obj2(0-18), StructSlice1:nil, StructSlice2:nil, StructSlice3:nil, IntPtr:1, SlicePtr:$Obj4}(0), {IntSlice1:$Obj1(0-9), IntSlice2:nil, IntSlice3:$Obj1(2-4), ByteSlice:$Obj2(0-18), StructSlice1:nil, StructSlice2:nil, StructSlice3:nil, IntPtr:1, SlicePtr:$Obj4}(1), {IntSlice1:$Obj1(0-9), IntSlice2:nil, IntSlice3:$Obj1(2-4), ByteSlice:$Obj2(0-18), StructSlice1:nil, StructSlice2:nil, StructSlice3:nil, IntPtr:1, SlicePtr:$Obj4}(2)]}`
			for i := 1; i <= loopCnt; i++ {
				So(StringByConf(cs, Config{}), ShouldEqual, expectStr)
			}
		})

		Convey("NoBaseKindsInfoOnly", func() {
			const loopCnt = 1000
			const expectStr = `(tostr.CommonStruct){IntSlice1:([]int)$Obj1(0-9), IntSlice2:([]int)nil, IntSlice3:([]int)$Obj1(2-4), ByteSlice:([]uint8)$Obj2(0-18), StructSlice1:([]tostr.CommonStruct)$Obj3(0-2), StructSlice2:([]tostr.CommonStruct)$Obj3(1-2), StructSlice3:([]tostr.CommonStruct)[], IntPtr:(***int)1, SlicePtr:(*[]int)<Obj4>$Obj1(0-9)}, {<Obj1>:([]int)[1(0), 2, 3(2), 4, 5(4), 6, 7, 8, 9, 10(9)], <Obj2>:([]uint8)[116(0), 104, 105, 115, 32, 105, 115, 32, 98, 121, 116, 101, 32, 115, 108, 105, 99, 101, 33(18)], <Obj3>:([]tostr.CommonStruct)[(tostr.CommonStruct){IntSlice1:([]int)$Obj1(0-9), IntSlice2:([]int)nil, IntSlice3:([]int)$Obj1(2-4), ByteSlice:([]uint8)$Obj2(0-18), StructSlice1:([]tostr.CommonStruct)nil, StructSlice2:([]tostr.CommonStruct)nil, StructSlice3:([]tostr.CommonStruct)nil, IntPtr:(***int)1, SlicePtr:(*[]int)$Obj4}(0), (tostr.CommonStruct){IntSlice1:([]int)$Obj1(0-9), IntSlice2:([]int)nil, IntSlice3:([]int)$Obj1(2-4), ByteSlice:([]uint8)$Obj2(0-18), StructSlice1:([]tostr.CommonStruct)nil, StructSlice2:([]tostr.CommonStruct)nil, StructSlice3:([]tostr.CommonStruct)nil, IntPtr:(***int)1, SlicePtr:(*[]int)$Obj4}(1), (tostr.CommonStruct){IntSlice1:([]int)$Obj1(0-9), IntSlice2:([]int)nil, IntSlice3:([]int)$Obj1(2-4), ByteSlice:([]uint8)$Obj2(0-18), StructSlice1:([]tostr.CommonStruct)nil, StructSlice2:([]tostr.CommonStruct)nil, StructSlice3:([]tostr.CommonStruct)nil, IntPtr:(***int)1, SlicePtr:(*[]int)$Obj4}(2)]}`
			for i := 1; i <= loopCnt; i++ {
				So(StringByConf(cs, Config{
					InformationLevel: NoBaseKindsInfoOnly,
				}), ShouldEqual, expectStr)
			}
		})

		Convey("AllTypesInfo", func() {
			const loopCnt = 1000
			const expectStr = `(tostr.CommonStruct){IntSlice1:([]int)$Obj1(0-9), IntSlice2:([]int)nil, IntSlice3:([]int)$Obj1(2-4), ByteSlice:([]uint8)$Obj2(0-18), StructSlice1:([]tostr.CommonStruct)$Obj3(0-2), StructSlice2:([]tostr.CommonStruct)$Obj3(1-2), StructSlice3:([]tostr.CommonStruct)[], IntPtr:(***int)1, SlicePtr:(*[]int)<Obj4>$Obj1(0-9)}, {<Obj1>:([]int)[(int)1(0), (int)2, (int)3(2), (int)4, (int)5(4), (int)6, (int)7, (int)8, (int)9, (int)10(9)], <Obj2>:([]uint8)[(uint8)116(0), (uint8)104, (uint8)105, (uint8)115, (uint8)32, (uint8)105, (uint8)115, (uint8)32, (uint8)98, (uint8)121, (uint8)116, (uint8)101, (uint8)32, (uint8)115, (uint8)108, (uint8)105, (uint8)99, (uint8)101, (uint8)33(18)], <Obj3>:([]tostr.CommonStruct)[(tostr.CommonStruct){IntSlice1:([]int)$Obj1(0-9), IntSlice2:([]int)nil, IntSlice3:([]int)$Obj1(2-4), ByteSlice:([]uint8)$Obj2(0-18), StructSlice1:([]tostr.CommonStruct)nil, StructSlice2:([]tostr.CommonStruct)nil, StructSlice3:([]tostr.CommonStruct)nil, IntPtr:(***int)1, SlicePtr:(*[]int)$Obj4}(0), (tostr.CommonStruct){IntSlice1:([]int)$Obj1(0-9), IntSlice2:([]int)nil, IntSlice3:([]int)$Obj1(2-4), ByteSlice:([]uint8)$Obj2(0-18), StructSlice1:([]tostr.CommonStruct)nil, StructSlice2:([]tostr.CommonStruct)nil, StructSlice3:([]tostr.CommonStruct)nil, IntPtr:(***int)1, SlicePtr:(*[]int)$Obj4}(1), (tostr.CommonStruct){IntSlice1:([]int)$Obj1(0-9), IntSlice2:([]int)nil, IntSlice3:([]int)$Obj1(2-4), ByteSlice:([]uint8)$Obj2(0-18), StructSlice1:([]tostr.CommonStruct)nil, StructSlice2:([]tostr.CommonStruct)nil, StructSlice3:([]tostr.CommonStruct)nil, IntPtr:(***int)1, SlicePtr:(*[]int)$Obj4}(2)]}`
			for i := 1; i <= loopCnt; i++ {
				So(StringByConf(cs, Config{
					InformationLevel: AllTypesInfo,
				}), ShouldEqual, expectStr)
			}
		})

		Convey("filte_all_slice", func() {
			const loopCnt = 1000
			const expectStr = `{IntPtr:1}`
			for i := 1; i <= loopCnt; i++ {
				So(StringByConf(cs, Config{
					FilterStructField: []func(reflect.Value, int) bool{
						func(obj reflect.Value, fieldIdx int) bool {
							return strings.Contains(obj.Type().Field(fieldIdx).Name, "Slice")
						},
					},
				}), ShouldEqual, expectStr)
			}
		})

		Convey("show_byte_slice_by_string", func() {
			const loopCnt = 1000
			const expectStr = `{IntSlice1:$Obj1(0-9), IntSlice2:nil, IntSlice3:$Obj1(2-4), ByteSlice:"this is byte slice!", StructSlice1:$Obj2(0-2), StructSlice2:$Obj2(1-2), StructSlice3:[], IntPtr:1, SlicePtr:<Obj3>$Obj1(0-9)}, {<Obj1>:[1(0), 2, 3(2), 4, 5(4), 6, 7, 8, 9, 10(9)], <Obj2>:[{IntSlice1:$Obj1(0-9), IntSlice2:nil, IntSlice3:$Obj1(2-4), ByteSlice:"this is byte slice!", StructSlice1:nil, StructSlice2:nil, StructSlice3:nil, IntPtr:1, SlicePtr:$Obj3}(0), {IntSlice1:$Obj1(0-9), IntSlice2:nil, IntSlice3:$Obj1(2-4), ByteSlice:"this is byte slice!", StructSlice1:nil, StructSlice2:nil, StructSlice3:nil, IntPtr:1, SlicePtr:$Obj3}(1), {IntSlice1:$Obj1(0-9), IntSlice2:nil, IntSlice3:$Obj1(2-4), ByteSlice:"this is byte slice!", StructSlice1:nil, StructSlice2:nil, StructSlice3:nil, IntPtr:1, SlicePtr:$Obj3}(2)]}`
			for i := 1; i <= loopCnt; i++ {
				So(StringByConf(cs, Config{
					ToString: func(obj reflect.Value) (objStr string) {
						if obj.Type() == reflect.TypeOf([]byte{}) {
							return strconv.Quote(string(obj.Interface().([]byte)))
						}
						return ""
					},
					FastSpecifyToStringProbe: func(obj reflect.Value) (hasSpecifyToString bool) {
						return obj.Type() == reflect.TypeOf([]byte{})
					},
				}), ShouldEqual, expectStr)
			}
		})

		Convey("warn_size_callback", func() {
			const loopCnt = 1000
			const expectStr = `Warn: len(string) is more than 100, [Str]={IntSlice1:$Obj1(0-9), IntSlice2:nil, IntSlice3:$Obj1(2-4), ByteSlice:"this is byte slice!", StructSlice1:`
			for i := 1; i <= loopCnt; i++ {
				warnSize := 100
				So(StringByConf(cs, Config{
					ToString: func(obj reflect.Value) (objStr string) {
						if obj.Type() == reflect.TypeOf([]byte{}) {
							return strconv.Quote(string(obj.Interface().([]byte)))
						}
						return ""
					},
					FastSpecifyToStringProbe: func(obj reflect.Value) (hasSpecifyToString bool) {
						return obj.Type() == reflect.TypeOf([]byte{})
					},
					WarnSize: &warnSize,
					ResultSizeWarnCallback: func(str string) bool {
						return false
					},
				}), ShouldEqual, expectStr)
			}
		})

		Convey("map_key_sort", func() {
			type MapKey struct {
				Bool      bool
				Int       int
				UInt      uint
				Float     float64
				Complex   complex64
				Ptr       *MapKey
				Interface error
				Array     [3]int
			}

			mapKey1 := MapKey{
				Bool:    true,
				Int:     2,
				UInt:    2,
				Float:   2,
				Complex: 2 + 2i,
				Array:   [3]int{1, 2, 3},
			}
			mapKey1.Ptr = &mapKey1
			mapKey1.Interface = fmt.Errorf("test error")

			mapKey2 := MapKey{
				Bool:    true,
				Int:     2,
				UInt:    2,
				Float:   2,
				Complex: 2 + 2i,
				Array:   [3]int{1, 2, 4},
			}
			mapKey2.Ptr = &mapKey1
			mapKey2.Interface = mapKey1.Interface

			mapKey3 := MapKey{
				Bool:    true,
				Int:     2,
				UInt:    2,
				Float:   2.34,
				Complex: 2 + 2i,
				Array:   [3]int{1, 3, 4},
			}
			mapKey3.Ptr = &mapKey1
			mapKey3.Interface = mapKey1.Interface

			m := map[MapKey]string{
				mapKey1: "value1",
				mapKey2: "value2",
				mapKey3: "value3",
			}

			const loopCnt = 1000
			expectStr := `{{Bool:true, Int:2, UInt:2, Float:2, Complex:(2+2i), Ptr:<Obj1>{Bool:true, Int:2, UInt:2, Float:2, Complex:(2+2i), Ptr:$Obj1, Interface:<Obj2>{s:"test error"}, Array:[1, 2, 3]}, Interface:$Obj2, Array:[1, 2, 3]}:"value1", {Bool:true, Int:2, UInt:2, Float:2, Complex:(2+2i), Ptr:$Obj1, Interface:$Obj2, Array:[1, 2, 4]}:"value2", {Bool:true, Int:2, UInt:2, Float:2.34, Complex:(2+2i), Ptr:$Obj1, Interface:$Obj2, Array:[1, 3, 4]}:"value3"}`

			for i := 1; i <= loopCnt; i++ {
				So(StringByConf(m, Config{}), ShouldEqual, expectStr)
			}
		})

		Convey("IntroductionRecursion", func() {
			type MyStruct struct {
				S   []MyStruct
				Str string
			}

			s := []MyStruct{
				{
					Str: "结构体字段Str",
				},
			}
			s[0].S = s

			const loopCnt = 1000
			expectStr := `([]tostr.MyStruct)$Obj1(0-0), {<Obj1>:([]tostr.MyStruct)[(tostr.MyStruct){S:([]tostr.MyStruct)$Obj1(0-0), Str:(string)"结构体字段Str"}(0)]}`

			for i := 1; i <= loopCnt; i++ {
				So(StringByConf(s, Config{
					InformationLevel: AllTypesInfo,
				}), ShouldEqual, expectStr)
			}
		})

		Convey("struct_filed_filter", func() {
			type MyStruct struct {
				str string
			}

			s := MyStruct{
				str: "不应该展示此字段",
			}

			const loopCnt = 1000
			expectStr := `(tostr.MyStruct){}`

			for i := 1; i <= loopCnt; i++ {
				df := defaultConfig
				df.InformationLevel = AllTypesInfo
				So(StringByConf(s, GetDefaultConfig().SetInformationLevel(AllTypesInfo).SetFilters().SetFilters(GetDefaultConfig().FilterStructField...)), ShouldEqual, expectStr)
			}
		})

		Convey("specify_to_string", func() {
			type MyStruct struct {
				t         time.Time
				byteSlice []byte
				jsonMsg   json.RawMessage
				T         time.Time
				ByteSlice []byte
				JsonMsg   json.RawMessage
			}

			testTime1, _ := time.Parse("2006-01-02 15:04:05", "2022-03-31 12:12:12")
			testTime2, _ := time.Parse("2006-01-02 15:04:05", "2022-04-01 12:12:12")
			ms := MyStruct{
				t:         testTime1,
				byteSlice: []byte{'a', 'b', 'c'},
				jsonMsg:   json.RawMessage{'d', 'e', 'f'},
				T:         testTime2,
				ByteSlice: []byte{'a', 'b', 'c'},
				JsonMsg:   json.RawMessage{'d', 'e', 'f'},
			}

			var toStringMap = map[reflect.Type]func(obj reflect.Value) string{
				reflect.TypeOf(time.Time{}): func(obj reflect.Value) string {
					if obj.CanInterface() {
						return "{time.Time:\"" + obj.Interface().(time.Time).Format("2006-01-02 15:04:05.000") + "\"}"
					}
					ptr := (*time.Time)((unsafe.Pointer)(obj.UnsafeAddr()))
					return "{time.Time:\"" + ptr.Format("2006-01-02 15:04:05.000") + "\"}"
				},
				reflect.TypeOf([]byte{}): func(obj reflect.Value) string {
					if obj.CanInterface() {
						return strconv.Quote(string(obj.Interface().([]byte)))
					}
					ptr := (*[]byte)((unsafe.Pointer)(obj.UnsafeAddr()))
					return strconv.Quote(string(*ptr))
				},
				reflect.TypeOf(json.RawMessage{}): func(obj reflect.Value) string {
					if obj.CanInterface() {
						return strconv.Quote(string(obj.Interface().(json.RawMessage)))
					}
					ptr := (*json.RawMessage)((unsafe.Pointer)(obj.UnsafeAddr()))
					return strconv.Quote(string(*ptr))
				},
			}

			var df = Config{
				ToString: func(o reflect.Value) (objStr string) {
					if f, inMap := toStringMap[o.Type()]; inMap {
						return f(o)
					}

					return fmt.Sprintf("<Fatal error：format func for %v not register>", o.Type())
				},
				FastSpecifyToStringProbe: func(o reflect.Value) (hasSpecifyToString bool) {
					_, inMap := toStringMap[o.Type()]
					return inMap
				},
			}

			const loopCnt = 1000
			expectStr := `{t:{time.Time:"2022-03-31 12:12:12.000"}, byteSlice:"abc", jsonMsg:"def", T:{time.Time:"2022-04-01 12:12:12.000"}, ByteSlice:"abc", JsonMsg:"def"}`

			for i := 1; i <= loopCnt; i++ {
				So(StringByConf(&ms, df), ShouldEqual, expectStr)
			}
		})

		Convey("func", func() {
			type MyStruct struct {
				func1 func(interface{}, Config) string
				func2 interface{}
			}

			s := MyStruct{
				func1: StringByConf,
				func2: fmt.Println,
			}

			const loopCnt = 1000
			const expectStr = `(tostr.MyStruct){func1:(func(interface {}, tostr.Config) string)<github.com/gogokit/tostr.StringByConf>, func2:(interface {})(func(...interface {}) (int, error))<fmt.Println>}`

			for i := 1; i <= loopCnt; i++ {
				So(StringByConf(s, Config{
					InformationLevel: AllTypesInfo,
				}), ShouldEqual, expectStr)
			}
		})

		Convey("filter_struct_field", func() {
			type MyStruct struct {
				Field1 func(interface{}, Config) string
				Field2 interface{}
			}

			s := MyStruct{
				Field1: StringByConf,
				Field2: fmt.Println,
			}

			const loopCnt = 1000
			expectStr := `{Field2:<fmt.Println>}`

			for i := 1; i <= loopCnt; i++ {
				So(Stringer(s, "Field1", "Field3").String(), ShouldEqual, expectStr)
			}
		})

		Convey("get_default_config", func() {
			conf := GetDefaultConfig()

			const loopCnt = 1000

			for i := 1; i <= loopCnt; i++ {
				So(conf.WarnSize, ShouldNotEqual, defaultConfig.WarnSize)
				So(len(conf.FilterStructField), ShouldEqual, len(defaultConfig.FilterStructField))
				for j := 0; j < len(conf.FilterStructField); j++ {
					So(&conf.FilterStructField[0], ShouldNotEqual, &defaultConfig.FilterStructField[0])
				}
			}
		})

		Convey("set_warn_size", func() {
			s := struct {
				Field1 int64
				Field2 interface{}
			}{}

			const loopCnt = 1000

			for i := 1; i <= loopCnt; i++ {
				So(StringerByConf(s, GetDefaultConfig().Clone()).String(), ShouldEqual, `{Field1:0, Field2:nil}`)
				So(strings.Contains(StringerByConf(s, GetDefaultConfig().Clone().SetWarnSize(4)).String(), "Warn: len(string) is more than 4"), ShouldEqual, true)
			}
		})

		Convey("protobuf_filter", func() {
			s := struct {
				Field1 int64
				Field2 MessageImpl
			}{}

			const loopCnt = 1000

			for i := 1; i <= loopCnt; i++ {
				So(StringerByConf(s, Config{
					FilterStructField: []func(obj reflect.Value, fieldIdx int) (hitFilter bool){ProtobufFieldFilter},
				}).String(), ShouldEqual, `{Field1:0, Field2:{Num:0}}`)
			}
		})

		Convey("stringer_kvs", func() {
			k1 := struct {
				Field1 int64
				Field2 interface{}
			}{}

			v1 := struct {
				Field1 int64
				Field2 interface{}
			}{}

			k2 := struct {
				Field1 struct{}
				Field2 interface{}
			}{}

			v2 := struct {
				Field1 int64
				Field2 interface{}
			}{}

			const (
				loopCnt   = 1000
				expectStr = `{{Field1:0, Field2:nil}:{Field1:0, Field2:nil},{Field1:{}, Field2:nil}:{Field1:0, Field2:nil},{nil:1}:[1, 2, 3, 4],nil:nil,<github.com/gogokit/tostr.Stringer>:<github.com/gogokit/tostr.StringByConf>}`
			)
			for i := 1; i <= loopCnt; i++ {
				So(StringerKvs(k1, v1, k2, v2, map[interface{}]interface{}{nil: 1}, []int64{1, 2, 3, 4}, nil, nil, Stringer, StringByConf).String(), ShouldEqual, expectStr)
			}
		})
	})
}

type MessageImpl struct {
	Num                  int64
	XXX_NoUnkeyedLiteral int64
}

func (MessageImpl) Reset() {}
func (MessageImpl) String() string {
	return ""
}
func (MessageImpl) ProtoMessage() {}
