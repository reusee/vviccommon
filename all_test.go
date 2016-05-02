package vviccommon

import "testing"

func TestShuffleTitle(t *testing.T) {
	s := "印花韩版T恤女短袖雪纺上衣 打底衫大码女装宽松烫钻"
	out, err := ShuffleTitle(s)
	if err != nil {
		t.Fatal(err)
	}
	_ = out
}
