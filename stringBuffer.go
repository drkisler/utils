package utils

import "strings"

type StringBuffer struct {
	strings.Builder
}

func (sb *StringBuffer) AppendStr(source string) *StringBuffer {
	_, _ = sb.WriteString(source)
	return sb
}
func (sb *StringBuffer) AppendLine(source string) *StringBuffer {
	_, _ = sb.WriteString(source)
	_ = sb.WriteByte('\r')
	_ = sb.WriteByte('\n')
	return sb
}
func (sb *StringBuffer) AppendRune(source rune) *StringBuffer {
	_, _ = sb.WriteRune(source)
	return sb
}
func (sb *StringBuffer) AppendRunes(source []rune) *StringBuffer {
	_, _ = sb.WriteString(string(source))
	return sb
}
func (sb *StringBuffer) AppendByte(source byte) *StringBuffer {
	_ = sb.WriteByte(source)
	return sb
}

func (sb *StringBuffer) AppendBytes(source []byte) *StringBuffer {
	_, _ = sb.Write(source)
	return sb
}
