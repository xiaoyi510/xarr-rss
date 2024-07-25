package opencc_ext

import (
	"github.com/xiaoyi510/opencc"
	"log"
)

var T2s *opencc.OpenCC

func init() {
	var err error
	T2s, err = opencc.New("t2s")
	if err != nil {
		log.Fatal(err)
	}
}
