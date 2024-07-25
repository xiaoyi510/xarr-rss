package torrent

import (
	"bytes"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"github.com/jackpal/bencode-go"
)

// TorrentHash 生成info TorrentHash
func TorrentHash(i interface{}) (string, error) {
	//  name, size, and piece hashes
	var buf bytes.Buffer
	err := bencode.Marshal(&buf, i)
	if err != nil {
		a := [20]byte{}
		return hex.EncodeToString([]byte(string(a[:]))), err
	}

	// info TorrentHash，在与跟踪器和对等设备对话时，它唯一地标识文件
	h := sha1.Sum(buf.Bytes())
	return hex.EncodeToString([]byte(string(h[:]))), nil
}

// TorrentHash 生成info TorrentHash
func TorrentHashV2(i interface{}) (string, error) {
	//  name, size, and piece hashes
	var buf bytes.Buffer
	err := bencode.Marshal(&buf, i)
	if err != nil {
		a := [20]byte{}
		return hex.EncodeToString([]byte(string(a[:]))), err
	}

	// info TorrentHash，在与跟踪器和对等设备对话时，它唯一地标识文件
	h := sha256.Sum256(buf.Bytes())
	return hex.EncodeToString([]byte(string(h[:]))), nil
}
