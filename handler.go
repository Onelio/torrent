package torrent

import (
	"strconv"
	"fmt"
)

//Bencode string implementation,
//Works as a bridge between plain string/[]byte and bencoded string/[]byte text
//IMPORTANT: Can be encoded or decoded but never edited, this tool is not meant for that
type BString struct {
	size		int 		`json:"-"`
	string		[]byte
}

//Bencode integer implementation,
//Doc specify float64 but we implement int since we only work with torrent
//Works as a bridge between int and bencoded string/[]byte integer
//IMPORTANT: Can be encoded or decoded but never edited, this tool is not meant for that
type BInteger struct {
	value		int
}

//Bencode dictionary implementation,
//Works as a bridge between map[string]interface{} and bencoded string/[]byte dictionary
//IMPORTANT: Can be encoded or decoded but never edited, this tool is not meant for that
type BDictionary struct {
	keys       []string 	`json:"-"`
	dictionary Map
}

//Bencode list implementation,
//Works as a bridge between []interface{} and bencoded string/[]byte list
//IMPORTANT: Can be encoded or decoded but never edited, this tool is not meant for that
type BList []interface{}

type Map map[string]interface{}

//Returns plain string value of BString
func (obj *BString)String() string {
	return string(obj.string)
}

//Returns plain []byte value of BString
func (obj *BString)ByteArray() []byte {
	return obj.string
}

//Returns BString value size
func (obj *BString)Size() int {
	return obj.size
}

//Returns bencoded string as string
func (obj *BString)BEncodedString() string {
	return fmt.Sprint(strconv.Itoa(obj.size), ":", string(obj.string))
}

//Returns int value of BInteger
func (obj *BInteger)Integer() int {
	return obj.Integer()
}

//Returns bencoded integer as string
func (obj *BInteger)BEncodedInteger() string {
	return fmt.Sprint("i", strconv.Itoa(obj.value), "e")
}

//Returns bencoded list as string
func (obj *BList)BEncodedList() string {
	return encodeList(obj)
}

//Returns BDictionary map[string]interface{}
//IMPORTANT: Unlike BList, here we need to call this function
//since we implement a workaround to fix the golang map random order
//to be used for encode.
func (obj *BDictionary)GetMap() Map {
	return obj.dictionary
}

//Returns bencoded dictionary as string
func (obj *BDictionary)BEncodedDictionary() string {
	return encodeDic(obj)
}

//Create BString from torrent []byte string
func AsBString(b *[]byte) (*BString, int) {
	var obj = &BString{}
	p := decodeString(b, obj)
	return obj, p
}

//Create BString from plain []byte string
func NewBString(b []byte) *BString {
	return &BString{
		string: b,
		size: len(b),
	}
}

//Create BInteger from torrent []byte integer
func AsBInt(b *[]byte) (*BInteger, int) {
	var obj = &BInteger{}
	p := decodeInt(b, obj)
	return obj, p
}

//Create BInteger from int
func NewBInt(b int) *BInteger {
	return &BInteger{
		value: b,
	}
}

//Create BList from torrent []byte list
func AsBList(b []byte) (*BList, int) {
	return decodeList(b)
}

//Create BList from []interface{}
func NewBList(b []interface{}) *BList {
	return createList(b)
}

//Create BDictionary from torrent []byte dictionary
func AsBDictionary(b []byte) (*BDictionary, int) {
	return decodeDic(b)
}

//Create BDictionary from map[string]interface{}
func NewBDictionary(b map[string]interface{}) *BDictionary {
	return createDic(b)
}