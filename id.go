package boba

import (
	"fmt"
	"reflect"
)

func generateBlockID(compositeID string, cursor Cursor) string {
	c := fmt.Sprintf("%v", cursor)
	hash := fmt.Sprintf("%03x", fnv32(compositeID+c))

	return fmt.Sprintf("bl-%d,%d-%s", cursor.Row, cursor.Col, hash)
}

func generateCompositeID(screenID string) string {
	hash := fmt.Sprintf("%03x", fnv32(screenID))
	return fmt.Sprintf("sc-%s-0", hash)
}

func generateScreenID(s Screen) string {
	name := screenName(s)

	hash := fmt.Sprintf("%03x", fnv32(name))
	return fmt.Sprintf("sc-%s", hash)
}

func fnv32(s string) uint32 {
	h := uint32(2166136261)

	for i := range s {
		h ^= uint32(s[i])
		h *= 16777619
	}

	return h & 0xFFF // 3 hex digits
}

func screenName(s Screen) string {
	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	return t.Name()
}
