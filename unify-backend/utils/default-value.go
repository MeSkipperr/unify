package utils

func DefaultString(val, def string) string {
	if val == "" {
		return def
	}
	return val
}

func DefaultInt(val, def int) int {
	if val == 0 {
		return def
	}
	return val
}
