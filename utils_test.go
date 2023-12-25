package ddgo

import "testing"

func TestSubStrAndTrim(t *testing.T) {
	str := "/home/page/user1"
	step := "/home/page"
	if s := SubStrAndTrim(str, step); s != "user1" {
		t.Error(s)
	}
	step = "/home"
	if s := SubStrAndTrim(str, step); s != "page/user1" {
		t.Error(s)
	}
	step = "/home/page/"
	if s := SubStrAndTrim(str, step); s != "user1" {
		t.Error(s)
	}
	step = "/home/page/user1"
	if s := SubStrAndTrim(str, step); s != "" {
		t.Error(s)
	}
}
