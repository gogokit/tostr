package tostr

import (
	"bytes"
	"fmt"
)

// 返回tostr.String生成的字符串格式化之后的结果, indentSpacesCount指定缩进使用的空格数. 例如: s为:"{slice:[1, 2, 3], m:{1:"string", 2:"var"}}", indentSpacesCount为4, 则格式化后的结果如下:
// {
//     slice:[
//         1,
//         2,
//         3
//     ],
//     m:{
//         1:"string",
//         2:"var"
//     }
// }
func Fmt(s string, indentSpacesCount uint8) (string, error) {
	f := &formatter{s: s, indentSpacesCount: indentSpacesCount}

	for f.curIdx < len(f.s) {
		if err := f.format(0); err != nil {
			return "", err
		}
	}

	return f.buf.String(), nil
}

type formatter struct {
	curIdx int
	s      string
	buf    bytes.Buffer

	// 层次缩进使用的空格符的个数
	indentSpacesCount uint8
}

// PS: 要求当前f.buf中最后一个字符的下一个位置即为当前对象的起始位置
func (f *formatter) format(indentColumns int) error {
	if !f.skipBlanks() {
		return nil
	}

	cc := f.curChar()
	f.buf.WriteByte(cc)

	curIdx := f.curIdx
	f.curIdx++

	switch cc {
	case ':':
		return nil
	case '{', '[':
		rightChar := map[byte]byte{
			'{': '}',
			'[': ']',
		}[cc]

		elemSn := 1
		for {
			if !f.skipBlanks() {
				return fmt.Errorf("syntax error: no pired %c for %c, leftCharIdx:%v", rightChar, cc, curIdx)
			}

			if f.curChar() == rightChar {
				if elemSn > 1 {
					f.buf.WriteByte('\n')
					for i := 1; i <= indentColumns; i++ {
						f.buf.WriteByte(' ')
					}
				}
				f.buf.WriteByte(f.curChar())
				f.curIdx++
				return nil
			}

			if f.curChar() == ',' {
				f.buf.WriteByte(f.curChar())
				f.buf.WriteByte('\n')
				for i := 1; i <= indentColumns+int(f.indentSpacesCount); i++ {
					f.buf.WriteByte(' ')
				}
				f.curIdx++
				continue
			}

			if elemSn == 1 {
				f.buf.WriteByte('\n')
				for i := 1; i <= indentColumns+int(f.indentSpacesCount); i++ {
					f.buf.WriteByte(' ')
				}
			}

			elemSn++

			// 使用f.indentSpacesCount个空格作为缩进
			if err := f.format(indentColumns + int(f.indentSpacesCount)); err != nil {
				return err
			}
		}
	case '(', '<':
		rightChar := map[byte]byte{
			'(': ')',
			'<': '>',
		}[cc]

		if !f.forward(func(c byte) bool {
			return c == rightChar
		}) {
			return fmt.Errorf("syntax error: no pired %c for %c, leftCharIdx:%v", rightChar, cc, curIdx)
		}
		f.buf.WriteByte(f.curChar())
		f.curIdx++
		if cc == '<' && f.curIdx < len(f.s) && f.curChar() == '[' {
			if err := f.format(indentColumns); err != nil {
				return err
			}
		}
	case ',':
		f.buf.WriteByte('\n')
		for i := 1; i <= indentColumns; i++ {
			f.buf.WriteByte(' ')
		}
	case ptrChar:
		for {
			f.forward(func(c byte) bool {
				return c != ptrChar
			})
			if f.curChar() != '<' {
				break
			}
			if err := f.format(indentColumns); err != nil {
				return err
			}
		}
		return f.format(indentColumns)
	case '"':
		var lastIsEscaped bool
		// 读取到第一个不是转义的"
		for ; f.curIdx < len(f.s); f.curIdx++ {
			// 每次循环头检测之前lastIsEscaped为true当且仅当扫描过的下标最大的转义符('\')具有转义的语意
			if f.preChar() == '\\' && !lastIsEscaped {
				f.buf.WriteByte('\\')
			}

			if f.curChar() != '\\' {
				f.buf.WriteByte(f.curChar())
			}

			if f.curChar() == '\\' {
				if f.preChar() == '\\' {
					lastIsEscaped = !lastIsEscaped
					continue
				}

				lastIsEscaped = true
				continue
			}

			if f.curChar() == '"' && (f.preChar() != '\\' || !lastIsEscaped) {
				f.curIdx++
				return nil
			}
		}
		return fmt.Errorf("syntax error: no pired %c for %c, leftCharIdx:%v", '"', cc, curIdx)
	default:
		f.forward(func(c byte) bool {
			return f.isDelimiter(c)
		})
	}
	return nil
}

func (f *formatter) skipBlanks() (hasMore bool) {
	isBlank := func(c byte) bool {
		return c == ' ' || c == '\t' || c == '\r' || c == '\n'
	}
	for ; f.curIdx < len(f.s) && isBlank(f.s[f.curIdx]); f.curIdx++ {
	}
	return f.curIdx < len(f.s)
}

func (f *formatter) curChar() byte {
	return f.s[f.curIdx]
}

func (f *formatter) preChar() byte {
	return f.s[f.curIdx-1]
}

func (f *formatter) isDelimiter(c byte) bool {
	switch c {
	case '(', ')', '[', ']', '{', '}', '<', '>', ',', ':', ptrChar:
		return true
	default:
		return false
	}
}

func (f *formatter) forward(isOk func(byte) bool) (find bool) {
	for ; f.curIdx < len(f.s) && !isOk(f.s[f.curIdx]); f.curIdx++ {
		f.buf.WriteByte(f.s[f.curIdx])
	}
	return f.curIdx < len(f.s)
}
