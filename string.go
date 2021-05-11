package kgo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/json-iterator/go"
	xhtml "golang.org/x/net/html"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"math/rand"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

// Md5 获取字节切片md5值.
// length指定结果长度,默认32.
func (ks *LkkString) Md5Byte(str []byte, length ...uint8) []byte {
	var l uint8 = 32
	if len(length) > 0 {
		l = length[0]
	}

	return md5Byte(str, l)
}

// Md5 获取字符串md5值.
// length指定结果长度,默认32.
func (ks *LkkString) Md5(str string, length ...uint8) string {
	var l uint8 = 32
	if len(length) > 0 {
		l = length[0]
	}

	return string(md5Byte([]byte(str), l))
}

// IsMd5 是否md5值.
func (ks *LkkString) IsMd5(str string) bool {
	return str != "" && RegMd5.MatchString(str)
}

// ShaXByte 计算字节切片的 shaX 散列值,x为1/256/512 .
func (ks *LkkString) ShaXByte(str []byte, x uint16) []byte {
	return shaXByte(str, x)
}

// ShaX 计算字符串的 shaX 散列值,x为1/256/512 .
func (ks *LkkString) ShaX(str string, x uint16) string {
	return string(shaXByte([]byte(str), x))
}

// IsSha1 是否Sha1值.
func (ks *LkkString) IsSha1(str string) bool {
	return str != "" && RegSha1.MatchString(str)
}

// IsSha256 是否Sha256值.
func (ks *LkkString) IsSha256(str string) bool {
	return str != "" && RegSha256.MatchString(str)
}

// IsSha512 是否Sha512值.
func (ks *LkkString) IsSha512(str string) bool {
	return str != "" && RegSha512.MatchString(str)
}

// Index 查找子串sub在字符串str中第一次出现的位置,不存在则返回-1;
// ignoreCase为是否忽略大小写.
func (ks *LkkString) Index(str, sub string, ignoreCase bool) int {
	if str == "" || sub == "" {
		return -1
	}

	if ignoreCase {
		str = strings.ToLower(str)
		sub = strings.ToLower(sub)
	}

	return strings.Index(str, sub)
}

// LastIndex 查找子串sub在字符串str中最后一次出现的位置,不存在则返回-1;
// ignoreCase为是否忽略大小写.
func (ks *LkkString) LastIndex(str, sub string, ignoreCase bool) int {
	if str == "" || sub == "" {
		return -1
	}

	if ignoreCase {
		str = strings.ToLower(str)
		sub = strings.ToLower(sub)
	}

	return strings.LastIndex(str, sub)
}

// Addslashes 使用反斜线引用字符串.
func (ks *LkkString) Addslashes(str string) string {
	var buf bytes.Buffer
	for _, char := range str {
		switch char {
		case '\'', '"', '\\':
			buf.WriteRune('\\')
		}
		buf.WriteRune(char)
	}
	return buf.String()
}

// Stripslashes 反引用一个引用字符串.
func (ks *LkkString) Stripslashes(str string) string {
	var buf bytes.Buffer
	l, skip := len(str), false
	for i, char := range str {
		if skip {
			skip = false
		} else if char == '\\' {
			if i+1 < l && str[i+1] == '\\' {
				skip = true
			}
			continue
		}
		buf.WriteRune(char)
	}
	return buf.String()
}

// JsonEncode 对val变量进行 JSON 编码.
// 依赖库github.com/json-iterator/go.
func (ks *LkkString) JsonEncode(val interface{}) ([]byte, error) {
	var jsons = jsoniter.ConfigCompatibleWithStandardLibrary
	return jsons.Marshal(val)
}

// JsonDecode 对 JSON 格式的str字符串进行解码,注意res使用指针.
// 依赖库github.com/json-iterator/go.
func (ks *LkkString) JsonDecode(str []byte, res interface{}) error {
	var jsons = jsoniter.ConfigCompatibleWithStandardLibrary
	return jsons.Unmarshal(str, res)
}

// Utf8ToGbk UTF-8转GBK编码.
func (ks *LkkString) Utf8ToGbk(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := io.ReadAll(reader)
	return d, e
}

// GbkToUtf8 GBK转UTF-8编码.
func (ks *LkkString) GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := io.ReadAll(reader)
	return d, e
}

// IsUtf8 字符串是否UTF-8编码.
func (ks *LkkString) IsUtf8(s []byte) bool {
	return utf8.Valid(s)
}

// IsGbk 字符串是否GBK编码.
func (ks *LkkString) IsGbk(s []byte) (res bool) {
	length := len(s)
	var i, j int
	for i < length {
		j = i + 1
		//大于127的使用双字节编码,且落在gbk编码范围内的字符
		//GBK中每个汉字包含两个字节，第一个字节(首字节)的范围是0x81-0xFE(即129-254),第二个字节(尾字节)的范围是0x40-0xFE(即64-254)
		if s[i] > 0x7f && j < length {
			if s[i] >= 0x81 &&
				s[i] <= 0xfe &&
				s[j] >= 0x40 &&
				s[j] <= 0xfe {
				i += 2
				res = true
			} else {
				res = false
				break
			}
		} else {
			i++
		}
	}

	return
}

// Img2Base64 将图片字节转换为base64字符串.ext为图片扩展名,默认jpg.
func (ks *LkkString) Img2Base64(content []byte, ext ...string) string {
	var imgType string = "jpg"
	if len(ext) > 0 {
		imgType = strings.ToLower(ext[0])
	}

	return img2Base64(content, imgType)
}

// StartsWith 字符串str是否以sub开头.
func (ks *LkkString) StartsWith(str, sub string, ignoreCase bool) bool {
	if str != "" && sub != "" {
		i := ks.Index(str, sub, ignoreCase)
		return i == 0
	}

	return false
}

// EndsWith 字符串str是否以sub结尾.
func (ks *LkkString) EndsWith(str, sub string, ignoreCase bool) bool {
	if str != "" && sub != "" {
		i := ks.LastIndex(str, sub, ignoreCase)
		return i != -1 && (len(str)-len(sub)) == i
	}

	return false
}

// Trim 去除字符串首尾处的空白字符（或者其他字符）.
func (ks *LkkString) Trim(str string, characterMask ...string) string {
	mask := getTrimMask(characterMask)
	return strings.Trim(str, mask)
}

// IsEmpty 字符串是否为空(包括空格).
func (ks *LkkString) IsEmpty(str string) bool {
	if str == "" || len(ks.Trim(str)) == 0 {
		return true
	}

	return false
}

// IsLetters 字符串是否全(英文)字母组成.
func (ks *LkkString) IsLetters(str string) bool {
	for _, r := range str {
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') {
			return false
		}
	}
	return str != ""
}

// IsUpper 字符串是否全部大写.
func (ks *LkkString) IsUpper(str string) bool {
	for _, r := range str {
		if !unicode.IsUpper(r) && unicode.IsLetter(r) {
			return false
		}
	}
	return str != ""
}

// IsLower 字符串是否全部小写.
func (ks *LkkString) IsLower(str string) bool {
	for _, r := range str {
		if !unicode.IsLower(r) && unicode.IsLetter(r) {
			return false
		}
	}
	return str != ""
}

// HasLetter 字符串是否含有(英文)字母.
func (ks *LkkString) HasLetter(str string) bool {
	for _, r := range str {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			return true
		}
	}
	return false
}

// IsASCII 是否ASCII字符串.
func (ks *LkkString) IsASCII(str string) bool {
	//return str != "" && RegAscii.MatchString(str)
	n := len(str)
	for i := 0; i < n; i++ {
		if str[i] > 127 {
			return false
		}
	}

	return str != ""
}

// IsMultibyte 字符串是否含有多字节字符.
func (ks *LkkString) IsMultibyte(str string) bool {
	return str != "" && RegMultiByte.MatchString(str)
}

// HasFullWidth 是否含有全角字符.
func (ks *LkkString) HasFullWidth(str string) bool {
	return str != "" && RegFullWidth.MatchString(str)
}

// HasHalfWidth 是否含有半角字符.
func (ks *LkkString) HasHalfWidth(str string) bool {
	return str != "" && RegHalfWidth.MatchString(str)
}

// IsEnglish 字符串是否纯英文.letterCase是否检查大小写,枚举值(CASE_NONE,CASE_LOWER,CASE_UPPER).
func (ks *LkkString) IsEnglish(str string, letterCase LkkCaseSwitch) bool {
	switch letterCase {
	case CASE_LOWER:
		return str != "" && RegAlphaLower.MatchString(str)
	case CASE_UPPER:
		return str != "" && RegAlphaUpper.MatchString(str)
	default: //CASE_NONE
		return ks.IsLetters(str)
	}
}

// HasEnglish 是否含有英文字符,HasLetter的别名.
func (ks *LkkString) HasEnglish(str string) bool {
	return ks.HasLetter(str)
}

// HasChinese 字符串是否含有中文.
func (ks *LkkString) HasChinese(str string) bool {
	for _, r := range str {
		if unicode.Is(unicode.Scripts["Han"], r) {
			return true
		}
	}

	return false
}

// IsChinese 字符串是否全部中文.
func (ks *LkkString) IsChinese(str string) bool {
	return str != "" && RegChineseAll.MatchString(str)
}

// IsChineseName 字符串是否中文名称.
func (ks *LkkString) IsChineseName(str string) bool {
	return str != "" && RegChineseName.MatchString(str)
}

// IsWord 是否词语(不以下划线开头的中文、英文、数字、下划线).
func (ks *LkkString) IsWord(str string) bool {
	return str != "" && RegWord.MatchString(str)
}

// IsNumeric 字符串是否数值(不包含复数和科学计数法).
func (ks *LkkString) IsNumeric(str string) bool {
	if str == "" {
		return false
	}
	_, err := strconv.ParseFloat(str, 64)
	return err == nil
}

// IsAlphaNumeric 是否字母或数字.
func (ks *LkkString) IsAlphaNumeric(str string) bool {
	return str != "" && RegAlphaNumeric.MatchString(str)
}

// HasSpecialChar 字符串是否含有特殊字符.
func (ks *LkkString) HasSpecialChar(str string) bool {
	for _, r := range str {
		// IsPunct 判断 r 是否为一个标点字符 (类别 P)
		// IsSymbol 判断 r 是否为一个符号字符
		// IsMark 判断 r 是否为一个 mark 字符 (类别 M)
		if unicode.IsPunct(r) || unicode.IsSymbol(r) || unicode.IsMark(r) {
			return true
		}
	}

	return false
}

// IsJSON 字符串是否合法的json格式.
func (ks *LkkString) IsJSON(str string) bool {
	length := len(str)
	if length == 0 {
		return false
	} else if (str[0] != '{' || str[length-1] != '}') && (str[0] != '[' || str[length-1] != ']') {
		return false
	}

	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}

// IsIP 检查字符串是否IP地址.
func (ks *LkkString) IsIP(str string) bool {
	return str != "" && net.ParseIP(str) != nil
}

// IsIPv4 检查字符串是否IPv4地址.
func (ks *LkkString) IsIPv4(str string) bool {
	ipAddr := net.ParseIP(str)
	// 不是合法的IP地址
	if ipAddr == nil {
		return false
	}

	return ipAddr.To4() != nil && strings.ContainsRune(str, '.')
}

// IsIPv6 检查字符串是否IPv6地址.
func (ks *LkkString) IsIPv6(str string) bool {
	ipAddr := net.ParseIP(str)
	return ipAddr != nil && strings.ContainsRune(str, ':')
}

// IsDNSName 是否DNS名称.
func (ks *LkkString) IsDNSName(str string) bool {
	if str == "" || len(strings.Replace(str, ".", "", -1)) > 255 {
		// constraints already violated
		return false
	}
	return !ks.IsIP(str) && RegDNSname.MatchString(str)
}

// IsHost 字符串是否主机名(IP或DNS名称).
func (ks *LkkString) IsHost(str string) bool {
	return ks.IsIP(str) || ks.IsDNSName(str)
}

// IsDialAddr 是否网络拨号地址(形如127.0.0.1:80),用于net.Dial()检查.
func (ks *LkkString) IsDialAddr(str string) bool {
	h, p, err := net.SplitHostPort(str)
	if err == nil && h != "" && p != "" && (ks.IsDNSName(h) || ks.IsIP(h)) && isPort(p) {
		return true
	}

	return false
}

// IsMACAddr 是否MAC物理网卡地址.
func (ks *LkkString) IsMACAddr(str string) bool {
	_, err := net.ParseMAC(str)
	return err == nil
}

// IsEmail 检查字符串是否邮箱.参数validateTrue,是否验证邮箱主机的真实性.
func (ks *LkkString) IsEmail(email string, validateHost bool) (bool, error) {
	//长度检查
	length := len(email)
	at := strings.LastIndexByte(email, '@')
	if (length < 6 || length > 254) || (at <= 0 || at > length-3) {
		return false, fmt.Errorf("[IsEmail] invalid email length")
	}

	//验证邮箱格式
	chkFormat := RegEmail.MatchString(email)
	if !chkFormat {
		return false, fmt.Errorf("[IsEmail] invalid email format")
	}

	//验证主机
	if validateHost {
		host := email[at+1:]
		if _, err := net.LookupMX(host); err != nil {
			//因无法确定mx主机的smtp端口,所以去掉Hello/Mail/Rcpt检查邮箱是否存在
			//仅检查主机是否有效
			if _, err := net.LookupIP(host); err != nil {
				return false, err
			}
		}
	}

	return true, nil
}

// IsMobilecn 检查字符串是否中国大陆手机号.
func (ks *LkkString) IsMobilecn(str string) bool {
	return str != "" && RegMobilecn.MatchString(str)
}

// IsTel 是否固定电话或400/800电话.
func (ks *LkkString) IsTel(str string) bool {
	return str != "" && RegTelephone.MatchString(str)
}

// Jsonp2Json 将jsonp转为json串.
// Example: forbar({a:"1",b:2}) to {"a":"1","b":2}
func (ks *LkkString) Jsonp2Json(str string) (string, error) {
	start := strings.Index(str, "(")
	end := strings.LastIndex(str, ")")

	if start == -1 || end == -1 {
		return "", errors.New("[Jsonp2Json] invalid jsonp.")
	}

	start += 1
	if start >= end {
		return "", errors.New("[Jsonp2Json] invalid jsonp.")
	}

	res := str[start:end]

	return res, nil
}

// Strpos 查找字符串首次出现的位置,找不到时返回-1.
// haystack在该字符串中进行查找,needle要查找的字符串;
// offset起始位置,为负数时时,搜索会从字符串结尾指定字符数开始.
func (ks *LkkString) Strpos(haystack, needle string, offset int) int {
	length := len(haystack)
	if length == 0 || offset > length || -offset > length {
		return -1
	}

	if offset < 0 {
		offset += length
	}
	pos := strings.Index(haystack[offset:], needle)
	if pos == -1 {
		return -1
	}
	return pos + offset
}

// Stripos  查找字符串首次出现的位置（不区分大小写）,找不到时返回-1.
// haystack在该字符串中进行查找,needle要查找的字符串;
// offset起始位置,为负数时时,搜索会从字符串结尾指定字符数开始.
func (ks *LkkString) Stripos(haystack, needle string, offset int) int {
	length := len(haystack)
	if length == 0 || offset > length || -offset > length {
		return -1
	}

	if offset < 0 {
		offset += length
	}
	pos := ks.Index(haystack[offset:], needle, true)
	if pos == -1 {
		return -1
	}
	return pos + offset
}

// Dstrpos 检查字符串str是否包含数组arr的元素之一,返回检查结果和匹配的字符串.
// chkCase为是否检查大小写.
func (ks *LkkString) Dstrpos(str string, arr []string, chkCase bool) (bool, string) {
	if len(str) == 0 || len(arr) == 0 {
		return false, ""
	}

	for _, v := range arr {
		if (chkCase && ks.Strpos(str, v, 0) != -1) || (!chkCase && ks.Stripos(str, v, 0) != -1) {
			return true, v
		}
	}

	return false, ""
}

// Nl2br 将换行符转换为br标签.
func (ks *LkkString) Nl2br(str string) string {
	return strings.Replace(str, "\n", "<br />", -1)
}

// Br2nl 将br标签转换为换行符.
func (ks *LkkString) Br2nl(str string) string {
	// <br> , <br /> , <br/>
	// <BR> , <BR /> , <BR/>

	l := len(str)
	buf := make([]byte, 0, l) //prealloca

	for i := 0; i < l; i++ {
		switch str[i] {
		case 60: //<
			if l >= i+3 {
				/*
					b = 98
					B = 66
					r = 82
					R = 114
					SPACE = 32
					/ = 47
					> = 62
				*/

				if l >= i+3 && ((str[i+1] == 98 || str[i+1] == 66) && (str[i+2] == 82 || str[i+2] == 114) && str[i+3] == 62) { // <br> || <BR>
					buf = append(buf, bytLinefeed...)
					i += 3
					continue
				}

				if l >= i+4 && ((str[i+1] == 98 || str[i+1] == 66) && (str[i+2] == 82 || str[i+2] == 114) && str[i+3] == 47 && str[i+4] == 62) { // <br/> || <BR/>
					buf = append(buf, bytLinefeed...)
					i += 4
					continue
				}

				if l >= i+5 && ((str[i+1] == 98 || str[i+1] == 66) && (str[i+2] == 82 || str[i+2] == 114) && str[i+3] == 32 && str[i+4] == 47 && str[i+5] == 62) { // <br /> || <BR />
					buf = append(buf, bytLinefeed...)
					i += 5
					continue
				}
			}
		default:
			buf = append(buf, str[i])
		}
	}

	return string(buf)
}

// RemoveSpace 移除字符串中的空白字符.
// all为true时移除全部空白,为false时只替换连续的空白字符为一个空格.
func (ks *LkkString) RemoveSpace(str string, all bool) string {
	if all && str != "" {
		return strings.Join(strings.Fields(str), "")
	} else if str != "" {
		//先将2个以上的连续空白符转为空格
		str = RegWhitespaceDuplicate.ReplaceAllString(str, " ")
		//再将[\t\n\f\r]等转为空格
		str = RegWhitespace.ReplaceAllString(str, " ")
	}

	return strings.TrimSpace(str)
}

// StripTags 过滤html标签.
func (ks *LkkString) StripTags(str string) string {
	return RegHtmlTag.ReplaceAllString(str, "")
}

// Html2Text 将html转换为纯文本.
func (ks *LkkString) Html2Text(str string) string {
	domDoc := xhtml.NewTokenizer(strings.NewReader(str))
	previousStartToken := domDoc.Token()
	var text string
loopDom:
	for {
		nx := domDoc.Next()
		switch {
		case nx == xhtml.ErrorToken:
			break loopDom // End of the document
		case nx == xhtml.StartTagToken:
			previousStartToken = domDoc.Token()
		case nx == xhtml.TextToken:
			if chk, _ := ks.Dstrpos(previousStartToken.Data, TextHtmlExcludeTags, false); chk {
				continue
			}

			text += " " + strings.TrimSpace(xhtml.UnescapeString(string(domDoc.Text())))
		}
	}

	return ks.RemoveSpace(text, false)
}

// ParseStr 将URI查询字符串转换为字典.
func (ks *LkkString) ParseStr(encodedString string, result map[string]interface{}) error {
	// split encodedString.
	if encodedString[0] == '?' {
		encodedString = strings.TrimLeft(encodedString, "?")
	}

	parts := strings.Split(encodedString, "&")
	for _, part := range parts {
		pos := strings.Index(part, "=")
		if pos <= 0 {
			continue
		}
		key, err := url.QueryUnescape(part[:pos])
		if err != nil {
			return err
		}
		for key[0] == ' ' && key[1:] != "" {
			key = key[1:]
		}
		if key == "" || key[0] == '[' {
			continue
		}
		value, err := url.QueryUnescape(part[pos+1:])
		if err != nil {
			return err
		}

		// split into multiple keys
		var keys []string
		left := 0
		for i, k := range key {
			if k == '[' && left == 0 {
				left = i
			} else if k == ']' {
				if left > 0 {
					if len(keys) == 0 {
						keys = append(keys, key[:left])
					}
					keys = append(keys, key[left+1:i])
					left = 0
					if i+1 < len(key) && key[i+1] != '[' {
						break
					}
				}
			}
		}
		if len(keys) == 0 {
			keys = append(keys, key)
		}
		// first key
		first := ""
		for i, chr := range keys[0] {
			if chr == ' ' || chr == '.' || chr == '[' {
				first += "_"
			} else {
				first += string(chr)
			}
			if chr == '[' {
				first += keys[0][i+1:]
				break
			}
		}
		keys[0] = first

		// build nested map
		if err := buildQueryMap(result, keys, value); err != nil {
			return err
		}
	}

	return nil
}

// ParseUrl 解析URL,返回其组成部分.
// component为需要返回的组成;
// -1: all; 1: scheme; 2: host; 4: port; 8: user; 16: pass; 32: path; 64: query; 128: fragment .
func (ks *LkkString) ParseUrl(str string, component int16) (map[string]string, error) {
	u, err := url.Parse(str)
	if err != nil {
		return nil, err
	}
	if component == -1 {
		component = 1 | 2 | 4 | 8 | 16 | 32 | 64 | 128
	}
	var res = make(map[string]string)
	if (component & 1) == 1 {
		res["scheme"] = u.Scheme
	}
	if (component & 2) == 2 {
		res["host"] = u.Hostname()
	}
	if (component & 4) == 4 {
		res["port"] = u.Port()
	}
	if (component & 8) == 8 {
		res["user"] = u.User.Username()
	}
	if (component & 16) == 16 {
		res["pass"], _ = u.User.Password()
	}
	if (component & 32) == 32 {
		res["path"] = u.Path
	}
	if (component & 64) == 64 {
		res["query"] = u.RawQuery
	}
	if (component & 128) == 128 {
		res["fragment"] = u.Fragment
	}
	return res, nil
}

// UrlEncode 编码 URL 字符串.
func (ks *LkkString) UrlEncode(str string) string {
	return url.QueryEscape(str)
}

// UrlDecode 解码已编码的 URL 字符串.
func (ks *LkkString) UrlDecode(str string) (string, error) {
	return url.QueryUnescape(str)
}

// RawUrlEncode 按照 RFC 3986 对 URL 进行编码.
func (ks *LkkString) RawUrlEncode(str string) string {
	return strings.Replace(url.QueryEscape(str), "+", "%20", -1)
}

// RawUrlDecode 对已编码的 URL 字符串进行解码.
func (ks *LkkString) RawUrlDecode(str string) (string, error) {
	return url.QueryUnescape(strings.Replace(str, "%20", "+", -1))
}

// HttpBuildQuery 根据参数生成 URL-encode 之后的请求字符串.
func (ks *LkkString) HttpBuildQuery(queryData url.Values) string {
	return queryData.Encode()
}

// FormatUrl 格式化URL.
func (ks *LkkString) FormatUrl(str string) string {
	if str != "" {
		if ks.Strpos(str, "://", 0) == -1 {
			str = "http://" + str
		}

		// 将"\"替换为"/"
		str = strings.ReplaceAll(str, "\\", "/")

		// 将连续的"//"或"\\"或"\/",替换为"/"
		str = RegUrlBackslashDuplicate.ReplaceAllString(str, "$1/")
	}

	return str
}

// GetDomain 从URL字符串中获取域名.
// 可选参数isMain,默认为false,取完整域名;为true时,取主域名(如abc.test.com取test.com).
func (ks *LkkString) GetDomain(str string, isMain ...bool) string {
	str = ks.FormatUrl(str)
	u, err := url.Parse(str)
	main := false
	if len(isMain) > 0 {
		main = isMain[0]
	}

	if err != nil || !strings.Contains(str, ".") {
		return ""
	} else if !main {
		return u.Hostname()
	}

	parts := strings.Split(u.Hostname(), ".")
	domain := parts[len(parts)-2] + "." + parts[len(parts)-1]

	return domain
}

// ClearUrlPrefix 清除URL的前缀;
// str为URL字符串,prefix为前缀,默认"/".
func (ks *LkkString) ClearUrlPrefix(str string, prefix ...string) string {
	var p string = "/"
	if len(prefix) > 0 {
		p = prefix[0]
	}

	for p != "" && strings.HasPrefix(str, p) {
		str = str[len(p):]
	}

	return str
}

// ClearUrlSuffix 清除URL的后缀;
// str为URL字符串,suffix为后缀,默认"/".
func (ks *LkkString) ClearUrlSuffix(str string, suffix ...string) string {
	var s string = "/"
	if len(suffix) > 0 {
		s = suffix[0]
	}

	for s != "" && strings.HasSuffix(str, s) {
		str = str[0 : len(str)-len(s)]
	}

	return str
}

// Random 生成随机字符串.
// length为长度,rtype为枚举:
// RAND_STRING_ALPHA 字母;
// RAND_STRING_NUMERIC 数值;
// RAND_STRING_ALPHANUM 字母+数值;
// RAND_STRING_SPECIAL 字母+数值+特殊字符;
// RAND_STRING_CHINESE 仅中文.
func (ks *LkkString) Random(length uint8, rtype LkkRandString) string {
	if length == 0 {
		return ""
	}

	var letter []rune
	alphas := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbers := "0123456789"
	specials := "~!@#$%^&*()_+{}:|<>?`-=;,."

	rand.Seed(time.Now().UTC().UnixNano())

	switch rtype {
	case RAND_STRING_ALPHA:
		letter = []rune(alphas)
	case RAND_STRING_NUMERIC:
		letter = []rune(numbers)
	case RAND_STRING_ALPHANUM:
		letter = []rune(alphas + numbers)
	case RAND_STRING_SPECIAL:
		letter = []rune(alphas + numbers + specials)
	case RAND_STRING_CHINESE:
		letter = CommonChinese
	default:
		letter = []rune(alphas)
	}

	res := make([]rune, length)
	for i := range res {
		res[i] = letter[rand.Intn(len(letter))]
	}

	return string(res)
}

// DetectEncoding 匹配字符编码,TODO.
func (ks *LkkString) DetectEncoding() {
	//TODO 检查字符编码
}
