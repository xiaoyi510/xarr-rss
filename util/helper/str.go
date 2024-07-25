package helper

import (
	"bytes"
	"github.com/dlclark/regexp2"
	"sort"
	"strconv"
	"strings"
)

var replaceReg *regexp2.Regexp
var replaceList [][]byte

func init() {
	replaceReg = regexp2.MustCompile("[\\s`~!@#$%^&*()－ー_\\-+=～<>?:\"{}|,.\\/;'\\\\[\\]·~！@#￥%……&*（）——\\-+={}|《》？：“”【】、；‘'，。、×]", regexp2.IgnoreCase|regexp2.Compiled|regexp2.IgnorePatternWhitespace)
	replaceList = [][]byte{
		[]byte(" "),
		[]byte("["),
		[]byte(" "),
		[]byte("\t"),
		[]byte("`"),
		[]byte("~"),
		[]byte("!"),
		[]byte("@"),
		[]byte("#"),
		[]byte("$"),
		[]byte("%"),
		[]byte("^"),
		[]byte("&"),
		[]byte("*"),
		[]byte("("),
		[]byte(")"),
		[]byte("－"),
		[]byte("ー"),
		[]byte("_"),
		[]byte("-"),
		[]byte("+"),
		[]byte("="),
		[]byte("～"),
		[]byte("<"),
		[]byte(">"),
		[]byte("?"),
		[]byte(":"),
		[]byte("\""),
		[]byte("{"),
		[]byte("}"),
		[]byte("|"),
		[]byte(","),
		[]byte("."),
		[]byte("/"),
		[]byte(";"),
		[]byte("'"),
		[]byte("["),
		[]byte("]"),
		[]byte("·"),
		[]byte("~"),
		[]byte("！"),
		[]byte("@"),
		[]byte("#"),
		[]byte("￥"),
		[]byte("%"),
		[]byte("…"),
		[]byte("…"),
		[]byte("&"),
		[]byte("*"),
		[]byte("（"),
		[]byte("）"),
		[]byte("—"),
		[]byte("—"),
		[]byte("\\"),
		[]byte("-"),
		[]byte("+"),
		[]byte("="),
		[]byte("{"),
		[]byte("}"),
		[]byte("|"),
		[]byte("《"),
		[]byte("》"),
		[]byte("？"),
		[]byte("："),
		[]byte("“"),
		[]byte("”"),
		[]byte("【"),
		[]byte("】"),
		[]byte("、"),
		[]byte("；"),
		[]byte("‘"),
		[]byte("'"),
		[]byte("，"),
		[]byte("。"),
		[]byte("、"),
		[]byte("×"),
		[]byte("]"),
	}

}

// 获取字符串数字连续值
func GetStringNumberContinue(numbers []string) []string {
	// 先将数据排序
	sort.Slice(numbers, func(i, j int) bool {
		a, _ := strconv.Atoi(numbers[i])
		b, _ := strconv.Atoi(numbers[j])
		return a < b
	})
	ses := []string{}
	// 判断是否为连续值
	lenNumber := len(numbers)
	if lenNumber > 1 {
		var min, max int
		for index, number := range numbers {
			now, _ := strconv.Atoi(number)
			if min == 0 {
				max = now
				min = now
				continue
			}
			// 如果有最小值 判断是否为连续值
			if max+1 == now {
				max = now
			} else {
				// 不是连续值了
				// 保存上一组数据
				if min == max {
					ses = append(ses, strconv.Itoa(min))
				} else {
					ses = append(ses, strconv.Itoa(min)+"~"+strconv.Itoa(max))
				}
				min = now
				max = now
			}

			// 判断是否为最后一个
			if index+1 == lenNumber {
				if min == max {
					ses = append(ses, strconv.Itoa(min))
				} else {
					ses = append(ses, strconv.Itoa(min)+"~"+strconv.Itoa(max))
				}
			}
		}
	} else if lenNumber == 1 {
		ses = append(ses, numbers...)
	}

	return ses
}

func GetBetweenStr(str, start, end string) string {
	n := strings.Index(str, start)
	if n == -1 {
		n = 0
	} else {
		n = n + len(start) // 增加了else，不加的会把start带上
	}
	str = string([]byte(str)[n:])
	m := strings.Index(str, end)
	if m == -1 {
		m = len(str)
	}
	str = string([]byte(str)[:m])
	return str
}

func GetStrRight(str, start string, length int) string {
	n := strings.Index(str, start)
	if n == -1 {
		return ""
	} else {
		n = n + len(start) // 增加了else，不加的会把start带上
	}
	m := []rune(str)[n:]
	if length > len(m) {
		return string(m[:length])
	}
	return ""
}

func StrToInt(str string) int {
	atoi, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return atoi
}
func IntToStr(num int) string {
	return strconv.Itoa(num)
}

func StrToInt64(str string) int64 {
	atoi, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return int64(atoi)
}

func ReplaceRegString(str string) string {
	//reg, err := regexp.Compile("[^\u4e00-\u9fa5\u3040-\u309f\u30A0-\u30FF\u31F0-\u31FF\uAC00-\uD7AF\u1100-\u11ff\u3130-\u318f\x{0400}-\x{052f}a-zA-Z0-9]")
	//reg, err := regexp.Compile("[\\s`~!@#$%^&*()－ー_\\-+=～<>?:\"{}|,.\\/;'\\\\[\\]·~！@#￥%……&*（）——\\-+={}|《》？：“”【】、；‘'，。、×]")
	//if err != nil {
	//	panic(err)
	//}

	//return StrReplace(str, []string{
	//	" ",
	//	"[",
	//	" ",
	//	"\t",
	//	"`",
	//	"~",
	//	"!",
	//	"@",
	//	"#",
	//	"$",
	//	"%",
	//	"^",
	//	"&",
	//	"*",
	//	"(",
	//	")",
	//	"－",
	//	"ー",
	//	"_",
	//	"-",
	//	"+",
	//	"=",
	//	"～",
	//	"<",
	//	">",
	//	"?",
	//	":",
	//	"\"",
	//	"{",
	//	"}",
	//	"|",
	//	",",
	//	".",
	//	"/",
	//	";",
	//	"'",
	//	"[",
	//	"]",
	//	"·",
	//	"~",
	//	"！",
	//	"@",
	//	"#",
	//	"￥",
	//	"%",
	//	"…",
	//	"…",
	//	"&",
	//	"*",
	//	"（",
	//	"）",
	//	"—",
	//	"—",
	//	"\\",
	//	"-",
	//	"+",
	//	"=",
	//	"{",
	//	"}",
	//	"|",
	//	"《",
	//	"》",
	//	"？",
	//	"：",
	//	"“",
	//	"”",
	//	"【",
	//	"】",
	//	"、",
	//	"；",
	//	"‘",
	//	"'",
	//	"，",
	//	"。",
	//	"、",
	//	"×",
	//	"]",
	//}, []string{""})
	tmp := []byte(str)
	for _, it := range replaceList {
		tmp = bytes.ReplaceAll(tmp, it, []byte(""))
	}
	return string(tmp)
	//ret, err := replaceReg.Replace(str, "", -1, -1)
	//if err != nil {
	//	log.Panic(err)
	//}
	//return ret
}

func ReplaceRegStringSpace(str string) string {
	return strings.ReplaceAll(strings.ReplaceAll(str, "-", ""), " ", "")
}

// 批量字符串替换
func StrReplace(str string, search []string, replace []string) string {
	for k, v := range search {
		// 判断替换值是否只有一个
		re := ""
		if len(replace) == 1 {
			re = replace[0]
		} else if k <= len(replace)-1 {
			re = replace[k]
		}

		str = strings.Replace(str, v, re, -1)
	}
	return str
}
