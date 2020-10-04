package hippo

import "strings"
import "regexp"

//单个比对
func singleEquals(str string, col *AuthColumn) bool {
	return str == col.Value
}

func isInSet(str string, col *AuthColumn) bool {
	values := strings.Split(col.Value, ",")
	for _, value := range values {
		if str == value {
			return true
		}
	}
	return false
}

func isInRange() bool {
	return true
}

func isStartWith(str string, col *AuthColumn) bool {
	return strings.HasPrefix(str, col.Value)
}

func isEndWith(str string, col *AuthColumn) bool {
	return strings.HasSuffix(str, col.Value)
}

func isMatchRegex(str string, col *AuthColumn) bool {
	rep := regexp.MustCompile(col.Value)
	return rep.MatchString(str)
}

func isEqualsMicro() bool {
	return true
}
