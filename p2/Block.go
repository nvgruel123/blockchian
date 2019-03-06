package p2

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"golang.org/x/crypto/sha3"
)

type Block struct {
	Header Header
	Value  MerklePatriciaTrie
}

type Header struct {
	Height     int32
	Timestamp  int64
	Hash       string
	ParentHash string
	Size       int32
}

type JsonBlock struct {
	Hash       string            `json:"hash"`
	Timestamp  int64             `json:"timeStamp"`
	Height     int32             `json:"height"`
	ParentHash string            `json:"parentHash"`
	Size       int32             `json:"size"`
	Mpt        map[string]string `json:"mpt"`
}

func (block *Block) Initial(height int32, parentHash string, value MerklePatriciaTrie) {

	time := int64(time.Now().Unix())
	mpt_byte := []byte(fmt.Sprintf("%v", value))
	size := int32(len(mpt_byte))
	hash := hash_block(height, time, parentHash, value, size)
	header := Header{height, time, hash, parentHash, size}
	block.Header = header
	block.Value = value
}

func DecodeFromJson(jsonString string) Block {
	byte_arr := []byte(jsonString)
	block := &Block{}
	var json_response JsonBlock
	json.Unmarshal(byte_arr, &json_response)
	header := &Header{}
	header.Hash = json_response.Hash
	header.Height = json_response.Height
	header.ParentHash = json_response.ParentHash
	header.Size = json_response.Size
	header.Timestamp = json_response.Timestamp
	mpt := MerklePatriciaTrie{}
	mpt.Initial()
	mpt_map := json_response.Mpt
	for k, v := range mpt_map {
		mpt.Insert(k, v)
	}
	block.Header = *header
	block.Value = mpt
	return *block
}

func (block *Block) EncodeToJson() string {
	header := &block.Header
	json_response := &JsonBlock{}
	json_response.Hash = header.Hash
	json_response.Height = header.Height
	json_response.ParentHash = header.ParentHash
	json_response.Size = header.Size
	json_response.Timestamp = header.Timestamp
	value := &block.Value
	json_response.Mpt = value.Data
	data, _ := json.Marshal(json_response)
	return string(data)
}

func hash_block(height int32, timestamp int64, parentHash string, mpt MerklePatriciaTrie, size int32) string {

	root := mpt.Root
	hash_str := fmt.Sprintf("%v%v%v%v%v", height, timestamp, parentHash, root, size)
	sum := sha3.Sum256([]byte(hash_str))

	return hex.EncodeToString(sum[:])
}