package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

// UtxoItem the data structure of an UTXO item
type UtxoItem struct {
	Amount       *big.Int //utxo的面值
	FrozenHeight int64    //锁定until账本高度超过
}

func (item *UtxoItem) Loads(data []byte) error {
	decoder := json.NewDecoder(bytes.NewBuffer(data))
	return decoder.Decode(item)
}

var (
	refTxid = "f2173ac96c47a8c11575cc44fc810cd60e59c7098e8538a57970cc156ea3da85"
	addr    = "XgW8QLK8kMHx3YQ1GpqvQTsjWfKthXhQ3"
	memonic = "睛 瓷 景 雪 等 确 片 单 得 硅 奴 即"
	p       = "./data/blockchain/xuper/utxoVM/"
	//44a02d3dcc0f62158236547792243db89dc83e6342b042a5920e0edfbd2c2716
)

func main() {
	// fmt.Println(account.CreateAccount(1, 1))
	set()
	// get()
}

func set() {

	db, err := leveldb.OpenFile(p, &opt.Options{
		OpenFilesCacheCapacity: 16,
		BlockCacheCapacity:     16 / 2 * opt.MiB,
		WriteBuffer:            16 / 4 * opt.MiB, // Two of these are used internally
		Filter:                 filter.NewBloomFilter(10),
	})
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// 旧的 utxo
	// refTxid := "f2173ac96c47a8c11575cc44fc810cd60e59c7098e8538a57970cc156ea3da85"
	refTxidBytes, err := hex.DecodeString(refTxid)
	if err != nil {
		panic(err)
	}

	key := GenUtxoKeyWithPrefix([]byte(addr), []byte(refTxidBytes), 0)
	ui := UtxoItem{
		Amount:       big.NewInt(200),
		FrozenHeight: 0,
	}

	value, _ := json.Marshal(ui)

	err = db.Delete([]byte(key), nil)
	if err != nil {
		panic(err)
	}

	// refTxidNew := "f2173ac96c47a8c11575cc44fc810cd60e59c7098e8538a57970cc156ea3da85"
	refTxidNew := refTxid
	refTxidBytesNew, err := hex.DecodeString(refTxidNew)
	if err != nil {
		panic(err)
	}
	keyNew := GenUtxoKeyWithPrefix([]byte(addr), []byte(refTxidBytesNew), 0)

	err = db.Put([]byte(keyNew), value, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("updata utxo succ")
}

func get() {

	db, err := leveldb.OpenFile(p, &opt.Options{
		OpenFilesCacheCapacity: 16,
		BlockCacheCapacity:     16 / 2 * opt.MiB,
		WriteBuffer:            16 / 4 * opt.MiB, // Two of these are used internally
		Filter:                 filter.NewBloomFilter(10),
	})
	if err != nil {
		panic(err)
	}
	defer db.Close()

	refTxid := "f2173ac96c47a8c11575cc44fc810cd60e59c7098e8538a57970cc156ea3da85"

	refTxidBytes, err := hex.DecodeString(refTxid)
	if err != nil {
		panic(err)
	}
	key := GenUtxoKeyWithPrefix([]byte("XgW8QLK8kMHx3YQ1GpqvQTsjWfKthXhQ3"), refTxidBytes, int32(0))
	fmt.Println("key::", key)
	v, e := db.Get([]byte(key), nil)
	if e != nil {
		panic(e)
	}
	ui := new(UtxoItem)
	e = json.Unmarshal(v, ui)
	if e != nil {
		panic(e)
	}
	fmt.Println(ui)
}

func genUtxoKey(addr []byte, txid []byte, offset int32) string {
	return fmt.Sprintf("%s_%x_%d", addr, txid, offset)
}

func GenUtxoKeyWithPrefix(addr []byte, txid []byte, offset int32) string {
	baseUtxoKey := genUtxoKey(addr, txid, offset)
	return "U" + baseUtxoKey
}
