package utils

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"time"
	"unsafe"
)

// TEnBuffer 开辟一个内存空间，减少GC
type TEnBuffer struct {
	buff     []byte //内存
	capacity uint64 //限制容量,不得超过该容量
	//start    uint64 //开始位置
	stop uint64 //结束位置 相当于slice capacity
	curr uint64 //指针位置 相当于slice length
}

func NewEnBuff(maxCap uint64) *TEnBuffer {
	buff := make([]byte, 0, maxCap)
	return &TEnBuffer{buff, maxCap, 0, 0}
}

// StartAppend 开始使用内存 star curr 移动至最后
func (eb *TEnBuffer) StartAppend() {
	//eb.start = eb.stop
	eb.curr = 0
	eb.buff = make([]byte, 0, cap(eb.buff))
}

// StopAppend stop=curr curr=start
func (eb *TEnBuffer) StopAppend() {
	if eb.curr > eb.stop {
		eb.stop = eb.curr
	}
}

// Reset 重置空间
func (eb *TEnBuffer) Reset() {
	eb.buff = make([]byte, 0, cap(eb.buff))
	eb.curr = 0
	eb.stop = 0
}

// Release 释放内存空间
func (eb *TEnBuffer) Release() {
	eb.buff = make([]byte, 0)
	eb.capacity = 0
	eb.curr = 0
	//eb.start = 0
	eb.stop = 0
}

// WriteBool true : 1 false : 0
func (eb *TEnBuffer) WriteBool(val bool) {
	var bVal uint8
	bVal = 0
	if val {
		bVal = 1
	}
	eb.buff = append(eb.buff, []byte{bVal}...)
	eb.curr++
}

// WriteUint8 写入1个字节 或写入 uint8
func (eb *TEnBuffer) WriteUint8(val byte) {
	eb.buff = append(eb.buff, []byte{val}...)
	eb.curr++
}

// WriteUint16 uint16
func (eb *TEnBuffer) WriteUint16(val uint16) {
	temBuff := make([]byte, 2)
	binary.BigEndian.PutUint16(temBuff, val)
	copy(eb.buff[eb.curr:eb.curr+2], temBuff)
	eb.buff = append(eb.buff, temBuff...)
	eb.curr += 2
}

// WriteUint32 uint32
func (eb *TEnBuffer) WriteUint32(val uint32) {
	temBuff := make([]byte, 4)
	binary.BigEndian.PutUint32(temBuff, val)
	eb.buff = append(eb.buff, temBuff...)
	eb.curr += 4
}

// WriteUint64 uint64
func (eb *TEnBuffer) WriteUint64(val uint64) {
	temBuff := make([]byte, 8)
	binary.BigEndian.PutUint64(temBuff, val)
	eb.buff = append(eb.buff, temBuff...)
	eb.curr += 8
}

// WriteInt int 根据实际的大小存储值
func (eb *TEnBuffer) WriteInt(val int) {
	switch {
	case val <= math.MaxUint8:
		eb.WriteUint8(uint8(1))
		eb.WriteUint8(uint8(val))
	case val > math.MaxUint8 && val <= math.MaxUint16:
		eb.WriteUint8(uint8(2))
		eb.WriteUint16(uint16(val))
	case val > math.MaxUint16 && val <= math.MaxUint32:
		eb.WriteUint8(uint8(4))
		eb.WriteUint32(uint32(val))
	default:
		eb.WriteUint8(uint8(8))
		eb.WriteUint64(uint64(val))
	}
}

// WriteFloat float
func (eb *TEnBuffer) WriteFloat(val float32) {
	eval := math.Float32bits(val)
	eb.WriteUint32(eval)
}

// WriteFloat64 float64
func (eb *TEnBuffer) WriteFloat64(val float64) {
	eval := math.Float64bits(val)
	eb.WriteUint64(eval)
}

// WriteTime uint64
func (eb *TEnBuffer) WriteTime(val time.Time) {
	eb.WriteUint64(uint64(val.Unix()))
}

// WriteBytes 写入长度，再写入数据
func (eb *TEnBuffer) WriteBytes(val []byte) {
	l := len(val)
	eb.WriteInt(l)
	eb.buff = append(eb.buff, val...)
	eb.curr += uint64(l)
}

// AppendBytes 不写入长度，直接追加数据
func (eb *TEnBuffer) AppendBytes(val []byte) {
	l := len(val)
	eb.buff = append(eb.buff, val...)
	eb.curr += uint64(l)
}

// LoadBytes 使用数据直接复写内存
func (eb *TEnBuffer) LoadBytes(val []byte) {
	l := len(val)
	eb.buff = val
	eb.curr = uint64(l)
	eb.stop = eb.curr
	eb.curr = 0
}

// WriteString 转换为bytes写入
func (eb *TEnBuffer) WriteString(val string) {
	p := (*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{val, len(val)},
	))
	eb.WriteBytes(*p)
}

// WriteStruct 写入struct
func (eb *TEnBuffer) WriteStruct(val any) error {
	value, err := json.Marshal(val)
	if err != nil {
		return err
	}
	eb.WriteBytes(value)
	return nil
}

// ReadStruct 调用前必须 调用 StopAppend()
func (eb *TEnBuffer) ReadStruct(val any) error {
	tempBuff := eb.buff[:eb.stop]
	return json.Unmarshal(tempBuff, val)
}
func (eb *TEnBuffer) ReadBool() (bool, error) {
	if (eb.curr + 1) > eb.stop {
		return false, fmt.Errorf("out of index")
	}
	result := eb.buff[eb.curr] == 1
	eb.curr++
	return result, nil
}
func (eb *TEnBuffer) ReadUint8() (uint8, error) {
	if (eb.curr + 1) > eb.stop {
		return 0, fmt.Errorf("out of index")
	}
	result := eb.buff[eb.curr]
	eb.curr++
	return result, nil
}
func (eb *TEnBuffer) ReadUint16() (uint16, error) {
	if (eb.curr + 2) > eb.stop {
		return 0, fmt.Errorf("out of index")
	}
	result := binary.BigEndian.Uint16(eb.buff[eb.curr : eb.curr+2])
	eb.curr += 2
	return result, nil
}
func (eb *TEnBuffer) ReadUint32() (uint32, error) {
	if (eb.curr + 4) > eb.stop {
		return 0, fmt.Errorf("out of index")
	}
	result := binary.BigEndian.Uint32(eb.buff[eb.curr : eb.curr+4])
	eb.curr += 4
	return result, nil
}

func (eb *TEnBuffer) ReadUint64() (uint64, error) {
	if (eb.curr + 8) > eb.stop {
		return 0, fmt.Errorf("out of index")
	}
	result := binary.BigEndian.Uint64(eb.buff[eb.curr : eb.curr+8])
	eb.curr += 8
	return result, nil
}

func (eb *TEnBuffer) ReadInt() (int, error) {
	var bits uint8
	var err error

	if bits, err = eb.ReadUint8(); err != nil {
		return 0, err
	}
	switch bits {
	case 1:
		var val uint8
		if val, err = eb.ReadUint8(); err != nil {
			return 0, err
		}
		result := int(val)
		return result, nil
	case 2:
		var val uint16
		if val, err = eb.ReadUint16(); err != nil {
			return 0, err
		}
		result := int(val)
		return result, nil
	case 4:
		var val uint32
		if val, err = eb.ReadUint32(); err != nil {
			return 0, err
		}
		result := int(val)
		return result, nil
	case 8:
		var val uint64
		if val, err = eb.ReadUint64(); err != nil {
			return 0, err
		}
		result := int(val)
		return result, nil
	default:
		return 0, fmt.Errorf("数据格式错误")
	}
}

func (eb *TEnBuffer) ReadFloat() (float32, error) {
	if (eb.curr + 4) > eb.stop {
		return 0, errors.New("out of index")
	}
	temBuff := eb.buff[eb.curr : eb.curr+4]

	val := math.Float32frombits(binary.BigEndian.Uint32(temBuff))
	eb.curr += 4
	return val, nil
}

func (eb *TEnBuffer) ReadFloat64() (float64, error) {
	if (eb.curr + 8) > eb.stop {
		return 0, errors.New("out of index")
	}
	temBuff := eb.buff[eb.curr : eb.curr+8]
	eb.curr += 8
	val := math.Float64frombits(binary.BigEndian.Uint64(temBuff))
	return val, nil
}

func (eb *TEnBuffer) ReadTime() (time.Time, error) {
	val, err := eb.ReadUint64()
	if err != nil {
		return time.Now(), err
	}
	result := time.Unix(int64(val), 0)
	return result, err
}

func (eb *TEnBuffer) ReadBytes() ([]byte, error) {
	l, err := eb.ReadInt()
	if err != nil {
		return nil, err
	}
	if (eb.curr + uint64(l)) > eb.stop {
		return nil, fmt.Errorf("out of index")
	}
	result := eb.buff[eb.curr : eb.curr+uint64(l)]
	eb.curr += uint64(l)
	return result, nil
}

// AllBytes 读取当前的数据，不考虑长度 与 appendBytes 配合使用
func (eb *TEnBuffer) AllBytes() ([]byte, error) {
	result := eb.buff[:eb.stop]
	return result, nil
}

func (eb *TEnBuffer) ReadString() (string, error) {
	result, err := eb.ReadBytes()
	if err != nil {
		return "", err
	}
	return *(*string)(unsafe.Pointer(&result)), nil
}

func (eb *TEnBuffer) Encrypt(key []byte) {
	iKey := 0
	for iIndex := uint64(0); iIndex < eb.stop; iIndex++ {
		eb.buff[iIndex] = eb.buff[iIndex] ^ key[iKey]
		iKey++
		if iKey >= len(key) {
			iKey = 0
		}
	}
}

func (eb *TEnBuffer) Decrypt(key []byte) {
	iKey := 0
	for iIndex := uint64(0); iIndex < eb.stop; iIndex++ {
		eb.buff[iIndex] = eb.buff[iIndex] ^ key[iKey]
		iKey++
		if iKey >= len(key) {
			iKey = 0
		}
	}
}

// Compress 调用前必须调用 StopAppend()
func (eb *TEnBuffer) Compress() error {
	compressed := new(bytes.Buffer)
	w, _ := gzip.NewWriterLevel(compressed, gzip.BestCompression)

	data := eb.buff[0:eb.stop]
	if _, err := w.Write(data); err != nil {
		return err
	}

	if err := w.Flush(); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
	bResult := compressed.Bytes()
	eb.LoadBytes(bResult)

	return nil
	/*
		var data []byte
		if err := eb.ReadAll(&data); err != nil {
			return err
		}
		var compressed []byte
		if err := snappy.Encode(compressed, data); err != nil {
			return err
		}
		eb.buff = compressed
		return nil

	*/

}

// DeCompress 调用前必须调用 StopAppend()
func (eb *TEnBuffer) DeCompress() error {
	var data []byte
	var err error
	var reader *gzip.Reader
	var result []byte

	if data, err = eb.AllBytes(); err != nil {
		return err
	}
	if reader, err = gzip.NewReader(bytes.NewReader(data)); err != nil {
		return err
	}
	if result, err = io.ReadAll(reader); err != nil {
		return err
	}
	eb.LoadBytes(result)
	return nil
	/*
		var compressed []byte
		if err := eb.ReadAll(&compressed); err != nil {
			return err
		}
		var data []byte
		if err := snappy.Decode(data, compressed); err != nil {
			return err
		}
		eb.buff = data
		return nil

	*/
}

// EncodeBase64 转为base64编码 调用前必须调用 StopAppend()
func (eb *TEnBuffer) EncodeBase64() error {
	var data []byte
	var err error
	if data, err = eb.AllBytes(); err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	encoder := base64.NewEncoder(base64.StdEncoding, buf)
	if _, err = encoder.Write(data); err != nil {
		_ = encoder.Close()
		return err
	}
	//close后才写入数据
	_ = encoder.Close()
	eb.LoadBytes(buf.Bytes())
	return nil
}

// DecodeBase64 转为base64编码 调用前必须调用 StopAppend()
func (eb *TEnBuffer) DecodeBase64() error {
	var data []byte
	var err error
	var result []byte
	if data, err = eb.AllBytes(); err != nil {
		return err
	}
	if result, err = io.ReadAll(base64.NewDecoder(base64.StdEncoding, bytes.NewReader(data))); err != nil {
		return err
	}
	eb.LoadBytes(result)
	return nil
}
