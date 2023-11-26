package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"os"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

type TEnString struct {
	String string
}
type TChrPos = uint8

var DefaultKey = "BBx0T9WGTzrsAbiTb3HO5HLi031SyyVX"

const (
	PosKeyLeft TChrPos = iota
	PosKey
	PosValLeft
	PosVal
)

// GetDefaultKey for encript key
func GetDefaultKey() string {
	result := os.Getenv("DEFAULT_KEY")
	if result == "" {
		result = "BBx0T9WGTzrsAbiTb3HO5HLi031SyyVX"
	}
	return result
}

// IndexOfByte 列出指定字节的全部位置
func (str *TEnString) IndexOfByte(target byte) []int {
	result := make([]int, 0, len(str.String))
	iIndex := strings.IndexByte(str.String, target)
	if iIndex < 0 {
		return result
	}
	result = append(result, iIndex)
	for iIndex+1 < len(str.String) {
		i := strings.IndexByte(str.String[iIndex+1:], target)
		if i < 0 {
			return result
		}
		iIndex = i + iIndex + 1
		result = append(result, iIndex)
	}
	return result
}
func (str *TEnString) Load(source ...string) {
	var sb StringBuffer
	for iIndex := 0; iIndex < len(source); iIndex++ {
		sb.AppendStr(source[iIndex])
	}
	str.String = sb.String()
}

// LoadWithSplit 使用split拼接多个字符
func (str *TEnString) LoadWithSplit(split string, source ...string) {
	var sb StringBuffer
	isFirst := true
	for iIndex := 0; iIndex < len(source); iIndex++ {
		if isFirst {
			sb.AppendStr(source[iIndex])
			isFirst = false
		} else {
			sb.AppendStr(split).AppendStr(source[iIndex])
		}

	}
	str.String = sb.String()
}

// Pos 一次性列出指定多个字符串的位置，位置顺序与指定字符串的顺序一致
func (str *TEnString) Pos(targets ...string) []int {
	var result []int
	if len(targets) == 0 {
		return result
	}
	result = make([]int, len(targets))
	iIndex := strings.Index(str.String, targets[0])
	result[0] = iIndex
	if len(targets) == 1 {
		return result
	}
	iSkip := 0
	for idx := 1; idx < len(targets); idx++ {
		iSkip = iIndex + len(targets[idx-1])
		i := strings.Index(str.String[iSkip:], targets[idx])
		if i < 0 {
			result[idx] = -1
			continue
		}
		iIndex = i + iSkip
		result[idx] = iIndex
	}
	return result
}

// IndexOfString 列出指定字符串的全部位置
func (str *TEnString) IndexOfString(target string) []int {
	var result []int
	result = make([]int, 0, len(str.String))
	iIndex := strings.Index(str.String, target)
	if iIndex < 0 {
		return result
	}
	result = append(result, iIndex)
	iSkip := 0
	for {
		iSkip = iIndex + len(target)
		iIndex = strings.Index(str.String[iSkip:], target)
		if iIndex < 0 {
			return result
		}
		iIndex += iSkip
		result = append(result, iIndex)
	}

}

// SubStr 截取字符串,不包括开始和结束字符
func (str *TEnString) SubStr(start, stop string) string {
	arrPos := str.Pos(start, stop)
	if arrPos[0] < 0 && arrPos[1] < 0 {
		return ""
	}
	if arrPos[0] < 0 && arrPos[1] >= 0 {
		return str.String[:arrPos[1]]
	}
	if arrPos[0] >= 0 && arrPos[1] < 0 {
		return str.String[arrPos[0]+len(start):]
	}
	return str.String[arrPos[0]+len(start) : arrPos[1]]
}

func (str *TEnString) SubStrTrim(leftStr, rightStr string) string {
	source := []rune(str.String)
	runeleft := []rune(leftStr)[0]
	runeright := []rune(rightStr)[0]
	ileft := strings.IndexRune(str.String, runeleft)
	if ileft < 0 {
		return "error"
	}
	iright := len(source) - 1
	bl := false
	br := false
	for ileft < iright {
		if (unicode.IsSpace(source[ileft]) || source[ileft] == runeleft) && !bl {
			ileft++
		} else {
			bl = true
		}
		if (unicode.IsSpace(source[iright]) || source[iright] == runeright) && !br {
			iright--
		} else {
			br = true
		}
		if bl && br {
			return string(source[ileft : iright+1])
		}
	}
	return ""
}

// SubStrSkipQuote 截取指定字符串位置之间的字符，跳过括弧内的相同字符 待测试
func (str *TEnString) SubStrSkipQuote(leftStr, rightStr string, quote string) string {
	source := []rune(str.String)
	runeleft := []rune(leftStr)[0]
	runeright := []rune(rightStr)[0]
	runeQuote := []rune(quote)[0]
	ileft := 0
	iright := len(source) - 1
	blQuote := false
	bl := false
	brQuote := false
	br := false
	for ileft < iright {
		if source[ileft] == runeQuote {
			blQuote = !blQuote
		}
		if source[iright] == runeQuote {
			brQuote = !brQuote
		}
		if source[ileft] == runeleft && !blQuote && !bl {
			ileft++
			bl = true
		}
		if source[iright] == runeright && !brQuote && !br {
			br = true
		}
		if !bl {
			ileft++
		}
		if !br {
			iright--
		}
		if bl && br {
			return string(source[ileft:iright])
		}

	}
	return ""
}

// RemoveSubstr 移除开始和结束的字符，返回剩余的字符和移除的字符，均不包括开始和结束字符
func (str *TEnString) RemoveSubstr(start, stop string) (*string, *string) {
	var remaining string
	var subStr string
	arrPos := str.Pos(start, stop)
	if arrPos[0] < 0 && arrPos[1] < 0 {
		remaining = str.String
		return &remaining, nil
	}
	if arrPos[0] > 0 && arrPos[1] < 0 {
		remaining = str.String[:arrPos[0]]
		subStr = str.String[:arrPos[0]+len(start)]
		return &remaining, &subStr
	}
	if arrPos[0] < 0 && arrPos[1] > 0 {
		remaining = str.String[arrPos[1]+len(stop):]
		subStr = str.String[:arrPos[1]]
		return &remaining, &subStr
	}
	remaining = str.String[:arrPos[0]] + str.String[arrPos[1]+len(stop):]
	subStr = str.String[arrPos[0]+len(start) : arrPos[1]]
	return &remaining, &subStr
}
func (str *TEnString) ContainsWithoutQuote(substr, quote string) bool {
	bQuote := false
	for iIndex := 0; iIndex < len(str.String); iIndex++ {
		if iIndex+len(quote) <= len(str.String) && str.String[iIndex:iIndex+len(quote)] == quote {
			bQuote = !bQuote
		}
		if iIndex+len(substr) <= len(str.String) && str.String[iIndex:iIndex+len(substr)] == substr {
			if !bQuote {
				return true
			}
		}
	}
	return false
}

// ToMap 将字符串转换为map,split,quote 必须为单个字符,quote 可以为空
// a=b,c=d | a = b ,c = d
// a=b c=d    e=f
// a=b and   c = d
// a = 'bf'  c="af we"
func (str *TEnString) ToMap(split, compare, quote string) *map[string]string {
	vPos := PosKeyLeft            //当前的位置
	runeSplit := []rune(split)[0] //单字符Split rune值
	strTmp := str.String
	if compare != "=" {
		strTmp = strings.ReplaceAll(str.String, compare, "=")
	}
	runeSource := []rune(strTmp) //原始rune值

	runeCompare := rune(61) //比较符rune值
	runeQuote := rune(0)    //包裹符rune值
	if quote != "" {
		runeQuote = []rune(quote)[0]
	}
	runeEmpty := rune(32) //  空字符串
	var key string
	bQuote := false
	var sb StringBuffer
	result := make(map[string]string)
	for iIndex := 0; iIndex < len(runeSource); iIndex++ {
		switch runeSource[iIndex] {
		case runeEmpty:
			if bQuote {
				sb.AppendRune(runeSource[iIndex])
			}
		case runeCompare:
			if !bQuote {
				key = sb.String()
				sb.Reset()
				vPos = PosValLeft
			} else {
				sb.AppendStr(compare)
			}
		case runeSplit:
			if !bQuote {
				result[key] = sb.String()
				sb.Reset()
				vPos = PosKeyLeft
			} else {
				sb.AppendRune(runeSource[iIndex])
			}
		case runeQuote:
			if !bQuote && runeSource[iIndex-1] == runeQuote {
				sb.AppendRune(runeSource[iIndex])
			}
			bQuote = !bQuote

		default:
			sb.AppendRune(runeSource[iIndex])
			switch vPos {
			case PosKeyLeft:
				vPos = PosKey
			case PosValLeft:
				vPos = PosVal
			}
		}
	}
	result[key] = sb.String()
	sb.Reset()
	return &result
}

// ToMapAny 使用不同的compare进行拆解,split必须为单字符，quote 可空，并且必须为单字符
// "a>=b,c>=d , e =”' f' ,g = 'o=h'" -->map[a:b c:d e:' f g:o=h]
func (str *TEnString) ToMapAny(split, quote string, compare []string) *map[string]string {
	runeSource := []rune(str.String)
	runeSplit := []rune(split)[0]
	runeQuote := rune(0) //包裹符rune值
	if quote != "" {
		runeQuote = []rune(quote)[0]
	}
	if compare == nil {
		return nil
	}
	if len(compare) == 0 {
		return nil
	}
	sort.Slice(compare, func(i, j int) bool {
		return len(compare[i]) < len(compare[j])
	})

	var tmpArray []string //临时存放拆解的key=val
	bQuote := false
	var sb StringBuffer
	result := make(map[string]string)
	for iIndex := 0; iIndex < len(runeSource); iIndex++ {
		if runeSource[iIndex] == runeQuote {
			bQuote = !bQuote
			sb.AppendRune(runeSource[iIndex])
			continue
		}
		if runeSource[iIndex] == runeSplit {
			if bQuote {
				sb.AppendRune(runeSource[iIndex])
			} else {
				tmpArray = append(tmpArray, sb.String())
				sb.Reset()
			}
			continue
		}
		sb.AppendRune(runeSource[iIndex])
	}
	if sb.Len() > 0 {
		tmpArray = append(tmpArray, sb.String())
	}
	for idx := 0; idx < len(tmpArray); idx++ {
		tmpStr := TEnString{tmpArray[idx]}
	innerLoop:
		for iIndex := len(compare) - 1; iIndex >= 0; iIndex-- {
			if tmpStr.ContainsWithoutQuote(compare[iIndex], quote) {
				dmap := tmpStr.ToMap(split, compare[iIndex], quote)
				for k, v := range *dmap {
					result[k] = v
				}
				break innerLoop
			}
		}
	}
	return &result
}

func (str *TEnString) Encrypt(key string) string {
	// 转成字节数组
	origData := []byte(str.String)
	k := []byte(key)
	// 分组秘钥
	// NewCipher该函数限制了输入k的长度必须为16, 24或者32
	block, _ := aes.NewCipher(k)
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 补全码
	origData = func(ciphertext []byte, blocksize int) []byte {
		padding := blocksize - len(ciphertext)%blocksize
		padtext := bytes.Repeat([]byte{byte(padding)}, padding)
		return append(ciphertext, padtext...)
	}(origData, blockSize) // PKCS7Padding(origData, blockSize)
	// 加密模式
	blockMode := cipher.NewCBCEncrypter(block, k[:blockSize])
	// 创建数组
	cryted := make([]byte, len(origData))
	// 加密
	blockMode.CryptBlocks(cryted, origData)
	return base64.StdEncoding.EncodeToString(cryted)
}
func (str *TEnString) Decrypt(key string) string {
	// 转成字节数组
	crytedByte, _ := base64.StdEncoding.DecodeString(str.String)
	k := []byte(key)
	// 分组秘钥
	block, _ := aes.NewCipher(k)
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 加密模式
	blockMode := cipher.NewCBCDecrypter(block, k[:blockSize])
	// 创建数组
	orig := make([]byte, len(crytedByte))
	// 解密
	blockMode.CryptBlocks(orig, crytedByte)
	// 去补全码
	orig = func(origData []byte) []byte {
		length := len(origData)
		unpadding := int(origData[length-1])
		return origData[:(length - unpadding)]
	}(orig) //PKCS7UnPadding(orig)
	return string(orig)
}

/*
前端解密：

	async function decrypt(encrypted, key) {
	    const encryptedArray = Uint8Array.from(atob(encrypted), c => c.charCodeAt(0));
	    const iv = encryptedArray.slice(0, 12);
	    const data = encryptedArray.slice(12);
	    const keyData = await window.crypto.subtle.importKey(
	        "raw",
	        new TextEncoder().encode(key),
	        "AES-GCM",
	        false,
	        ["decrypt"]
	    );
	    const decrypted = await window.crypto.subtle.decrypt(
	        {
	            name: "AES-GCM",
	            iv: iv
	        },
	        keyData,
	        data
	    );
	    const decoder = new TextDecoder();
	    return decoder.decode(decrypted);
	}
*/
func (str *TEnString) Encrypt4Frontend(key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	data := []byte(str.String)

	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], data)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

/*
前端加密

	async function encrypt(message, key) {
	    const keyData = await window.crypto.subtle.importKey(
	        "raw",
	        new TextEncoder().encode(key),
	        "AES-GCM",
	        false,
	        ["encrypt"]
	    );
	    const iv = window.crypto.getRandomValues(new Uint8Array(12));
	    const data = new TextEncoder().encode(message);
	    const encrypted = await window.crypto.subtle.encrypt(
	        {
	            name: "AES-GCM",
	            iv: iv
	        },
	        keyData,
	        data
	    );
	    const encryptedArray = new Uint8Array(encrypted);
	    const result = new Uint8Array(iv.length + encryptedArray.length);
	    result.set(iv);
	    result.set(encryptedArray, iv.length);
	    return btoa(String.fromCharCode.apply(null, result));
	}

	async function main() {
	    const key = "example key 1234";
	    const message = "Hello, World!";
	    const encrypted = await encrypt(message, key);
	    console.log(encrypted);
	}
*/
func (str *TEnString) Decrypt4Frontend(key []byte) (string, error) {
	data, err := base64.StdEncoding.DecodeString(str.String)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(data) < 12 {
		return "", errors.New("带解密字符不得少于12个字符")
	}
	nonce := data[:12]
	data = data[12:]

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}

	plaintext, err := aesgcm.Open(nil, nonce, data, nil)
	if err != nil {
		panic(err)
	}

	return string(plaintext), nil
}

func (str *TEnString) LowerCase(quote string) string {
	var iQuote uint8
	iQuote = quote[0] //必须是ascii字符
	isASCII, hasUpper, skip := true, false, false
	lowerNum := uint8('a' - 'A')
	for i := 0; i < len(str.String); i++ {
		c := str.String[i]
		if c >= utf8.RuneSelf {
			isASCII = false
			break
		}
		hasUpper = hasUpper || ('A' <= c && c <= 'Z')
	}

	if isASCII { // optimize for ASCII-only strings.
		if !hasUpper {
			return str.String
		}
		var b strings.Builder
		b.Grow(len(str.String))
		for i := 0; i < len(str.String); i++ {
			c := str.String[i]
			if c == iQuote {
				skip = !skip
			}
			if skip {
				b.WriteByte(c)
				continue
			}
			if 'A' <= c && c <= 'Z' {
				c += lowerNum
			}
			b.WriteByte(c)
		}
		return b.String()
	}
	return strings.Map(unicode.ToLower, str.String)

}

// CutFromFirst 从左侧第一个指定字符开始去除右侧的字符，并移除改字符左侧的空字符
// CutFromFirst("abcdef (123 (5689","(")  ==> "abcdef"
func (str *TEnString) CutFromFirst(start string) string {
	iEmpty := uint8(32) //空字符
	iStart := []uint8(start)[0]

	iIndex := strings.Index(str.String, start)
	if iIndex < 0 {
		return str.String
	}
	for ; iIndex >= 0; iIndex-- {
		if str.String[iIndex] != iEmpty && str.String[iIndex] != iStart {
			return str.String[:iIndex+1]
		}
	}
	return ""
}

// CutFromLast ("abcdef (123 (5689","(")  ==> "abcdef (123"
func (str *TEnString) CutFromLast(start string) string {
	iEmpty := uint8(32) //空字符
	iStart := []uint8(start)[0]

	iIndex := strings.LastIndex(str.String, start)
	if iIndex < 0 {
		return str.String
	}
	for ; iIndex >= 0; iIndex-- {
		if str.String[iIndex] != iEmpty && str.String[iIndex] != iStart {
			return str.String[:iIndex+1]
		}
	}
	return ""
}

// TrimFromRight 从右侧开始移除空字符或指定的字符，直至出现其他的字符
// TrimFromRight("abcdef ghijkl  ;   ",";")  ==> "abcdef ghijkl"
func (str *TEnString) TrimFromRight(tailFlag string) string {
	iEmpty := uint8(32) //空字符
	iFlag := []uint8(tailFlag)[0]

	for iIndex := len(str.String) - 1; iIndex >= 0; iIndex-- {
		if str.String[iIndex] != iEmpty && str.String[iIndex] != iFlag {
			return str.String[:iIndex+1]
		}
	}
	return ""
}

// TrimFromLeft 从左侧开始移除空字符或指定的字符，直至出现其他的字符
func (str *TEnString) TrimFromLeft(headFlag string) string {
	iEmpty := uint8(32) //空字符
	iFlag := []uint8(headFlag)[0]

	for iIndex := 0; iIndex < len(str.String); iIndex++ {
		if str.String[iIndex] != iEmpty && str.String[iIndex] != iFlag {
			return str.String[iIndex:]
		}
	}
	return ""
}

/*
func (str *TEnString) ConvertToCharset(from, to string) string {
	srcCoder := mahonia.NewDecoder(from)
	srcResult := srcCoder.ConvertString(str.String)
	tagCoder := mahonia.NewDecoder(to)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}
*/

func (str *TEnString) Repeat(source, split string, count int) string {
	if count <= 0 {
		return ""
	}
	s := source + split

	// Since we cannot return an error on overflow,
	// we should panic if the repeat will generate
	// an overflow.
	// See Issue golang.org/issue/16237
	if len(s)*count/count != len(s) {
		panic("strings: Repeat count causes overflow")
	}

	n := len(s) * (count - 1)
	var b StringBuffer
	b.Grow(n)
	_, _ = b.WriteString(s)
	for b.Len() < n {
		if b.Len() <= n/2 {
			_, _ = b.WriteString(b.String())
		} else {
			_, _ = b.WriteString(b.String()[:n-b.Len()])
			break
		}
	}
	_, _ = b.WriteString(source)
	return b.String()
}
func (str *TEnString) EmptyVal(source, target string) string {
	if source == "" {
		return target
	}
	if len(strings.TrimSpace(source)) == 0 {
		return target
	}

	strings.TrimSpace(source)
	return source
}
func (str *TEnString) Split(split string) []string {
	return strings.Split(str.String, split)

}
