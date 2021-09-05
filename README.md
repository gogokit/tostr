tostr包用于生成指定golang对象对应的字符串, 该字符串以易于阅读的方式呈现指定对象的内容.   
tostr包中相关函数的实现基于一种明确定义的有向图, 该有向图中的结点和边如下:   
对象obj为有向图的源点, 任意结点的后继结点如下:    
&emsp;(1) 如果v.Kind()为reflect.Slice, 则切片v中包含的所有元素为v的后继结点   
&emsp;(2) 如果v.Kind()为reflect.Interface, 则v.Elem()为v的后继结点    
&emsp;(3) 如果v.Kind()为reflect.Array, 则v中包含的所有元素为v的后继结点   
&emsp;(4) 如果v.Kind()为reflect.Struct, 则v的所有字段对应的值为v的后继结点    
&emsp;(5) 如果v.Kind()为reflect.Map, 则v中所有Key, Value为v的后继结点        
对于任意结点v, 当且仅当v满足如下条件时, v不存在后继结点:      
&emsp;(1) v.Kind()与后面给出的某个值相等: reflect.Invalid, reflect.String, reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128, reflect.Chan, reflect.Func, reflect.UnsafePointer    
&emsp;(2) v.Kind()为reflect.Interface且v.IsNil()为true    
&emsp;(3) v.Kind()为reflect.Array且v.Len()为0    
&emsp;(4) v.Kind()为reflect.Slice且v.IsNil()为true或v.Len()为0    
&emsp;(5) v.Kind()为reflect.Struct且cur.NumField()为0    
&emsp;(6) v.Kind()的值为reflect.Map且v.IsNil()为true或v.Len()为0    
&emsp;(7) v.Kind()的值为reflect.Ptr且v.IsNil()为true     
所有指向指针的指针类型结点均视为两两不同的结点    
将[]T类型的切片中直接包含的元素视为*T, 也即指向T类型的指针, 该指针指向的对象即为切片真实包含的元素.    
如果两个对象满足如下条件之一, 则这两个对象在有向图中对应同一个结点:      
&emsp;(1) 为同类型的指针, 值不为nil且存储的地址值相等      
&emsp;(2) 为同类型的map, 值不为nil, map不为空且两个对象(map类型的变量实际对应一个指针值)的值相等    
为使生成字符串的长度尽量少, 同时保证的生成字符串的直观性, 在遍历过程中遇到的重复的指针类型和map类型的结点时第1次遍历时完整    
的呈现对象, 同时给对象编号, 编号形如Obj1, Obj2..., 之后再遇到该对象时使用占位符替换, 占位符形如$Obj1, $Obj2...,     
&emsp;使用tostr.String不会展示对象的类型信息, 如果需要展示对象的类型信息也可以使用tostr.StringByConf, 通过传入的conf进行控制      
例如: 下面的程序:    
```golang
package main

import (
	 "fmt"
	 "github.com/58kg/tostr"
)

func main() {
	 type Struct struct {
		 str string
	 }
	 s := Struct{
		 str: "str",
	 }
	 arr2 := []*Struct{&s, &s}
	 num := 1
	 arr1 := []*int{&num, &num}
	 fmt.Printf("arr2:%v\n", tostr.String(arr2))
	 fmt.Printf("arr1:%v\n", tostr.String(arr1))
	 arr2Str := tostr.StringByConf(arr2, tostr.Config{
		 InformationLevel: tostr.AllTypesInfo,
	 })
	 fmt.Printf("arr2Str:%v\n", arr2Str)
	 arr1Str := tostr.StringByConf(arr1, tostr.Config{
		 InformationLevel: tostr.AllTypesInfo,
	 })
	 fmt.Printf("arr1Str:%v\n", arr1Str)
}
```
输出内容如下:    
arr2:[<Obj1>{str:"str"}, $Obj1]    
arr1:[1, 1]     
arr2Str:([]*main.Struct)[(*<Obj1>main.Struct){str:(string)"str"}, $Obj1]    
arr1Str:([]*int)[(*int)1, (*int)1]    
      
为了进一步降低生成字符串的长度, 对于同类型的多个切片, 如果这些切片引用的底层数组之间存在共用的部分将对切片中的元素进行聚合, 保证该类型切片中
统一地址上的对象仅被打印一遍, 例如: 对于下面的程序:    
```golang
package main

import (
	 "fmt"
     "github.com/58kg/tostr"
)

func main() {
	 type Struct struct {
		 slice1 []int
		 slice2 []int
	 }
	 slice := []int{1, 2, 3, 4, 5}
	 s := Struct{
		 slice1: slice,
		 slice2: slice[1:3],
	 }
	 fmt.Printf("s:%v\n", tostr.String(s))
}
```
输出内容如下:    
{slice1:$Obj1(0-4), slice2:$Obj1(1-2)}, {<Obj1>:[1(0), 2(1), 3(2), 4, 5(4)]}    
     
上面的输出中slice1对应$Obj1(0-4)表示slice1中元素依次为Obj1的下标0至4之间的元素, slice2对应$Obj1(1-2)表示slice2中的元素依    
次为Obj1的下标1至2之间的元素, 切片Obj1中的每个元素后面括号中的数字表示元素在Obj1中的下标      
    
最后, 如果生成的字符串过大导致不便于阅读时可以使用本包下的Fmt函数格式化一下, 格式化之后层次分明, 阅读方便.    
