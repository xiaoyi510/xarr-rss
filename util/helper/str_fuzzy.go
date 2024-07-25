package helper

import (
	"golang.org/x/text/transform"
)

type nopTransformer struct{ transform.NopResetter }

func (nopTransformer) Transform(dst []byte, src []byte, atEOF bool) (int, int, error) {
	return 0, len(src), nil
}

func noopTransformer() transform.Transformer {
	return nopTransformer{}
}
func stringTransform(s string, t transform.Transformer) (transformed string) {
	// Fast path for the nop transformer to prevent unnecessary allocations.
	if _, ok := t.(nopTransformer); ok {
		return s
	}

	var err error
	transformed, _, err = transform.String(t, s)
	if err != nil {
		transformed = s
	}

	return
}
func MatchString(source, target string) (bool, string) {
	//transformer := noopTransformer()
	sourceNew := []rune((source))
	targetNew := []rune((target))

	lenDiff := len(targetNew) - len(sourceNew)

	if lenDiff < 0 {
		return false, ""
	}

	if lenDiff == 0 && source == target {
		return true, source
	}
	start := []rune{}
	startIndex := 0
	cureentIndex := 0
	index := 0

	// 循环查询
Goto_SourceReload:
	for index < len(sourceNew) {
		if string(start) == source {
			return true, string(start)
		}
		r1 := (sourceNew[index])

		// 判断还有内容没
		if len(targetNew) == 0 && len(string(start)) < len(source) {
			return false, ""
		}
		for i, r2 := range targetNew {
			//log.Println(string(r1), string(r2))
			// 判断两个字符是否相等
			if r1 == r2 {
				// 找到相同字符
				start = append(start, targetNew[i:i+1]...)
				//start += targetNew[i : i+1]
				// 记录查询开始位置
				startIndex += 1
				cureentIndex = startIndex - 1

				// 删除掉匹配的内容
				targetNew = targetNew[i+1:]

				// 找到了 让他重新找
				index++
				continue Goto_SourceReload
			}
			cureentIndex += 1
			// 判断是否当前错误次数过多
			if cureentIndex-startIndex > 2 {
				// 错误了 重新开始找
				if len(start) != 0 {
					index--
				}
				// 清空start
				start = []rune{}
				startIndex = 0
				cureentIndex = 0
				// 删除掉未匹配到的内容
				targetNew = targetNew[i+1:]

				// 重新去筛选
				continue Goto_SourceReload

			} else if len(start) > 0 {
				// 将中间间隔内容放进来
				start = append(start, targetNew[i:i+1]...)
			}
		}
		index++
	}
	if len(start) == 0 || len(string(start)) < len(source) {
		return false, ""
	}
	return true, string(start)

}
