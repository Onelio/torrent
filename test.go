package torrent

import (
	"fmt"
	"encoding/json"
	"io/ioutil"
	"net/url"
)

func TestString() {
	var estr = []byte("4:spam..")
	str, p := AsBString(&estr)
	fmt.Println(estr, str.string, p, estr[p:])
}

func TestInt() {
	var estr = []byte("i12...")
	bint, p := AsBInt(&estr)
	fmt.Println(estr, bint.value, p, estr[p:])
}

func TestLists() {
	var estr = []byte("l4:spam4:eggsl4:spam4:eggsee...")
	blist, p := AsBList(estr)
	res2B, _ := json.Marshal(blist)
	fmt.Println(estr, string(res2B), p, estr[p:])
}

func TestMaps() {
	var estr = []byte("d3:cow3:moo4:spam4:eggse...")
	blist, p := AsBDictionary(estr)
	res2B, _ := json.Marshal(blist)
	fmt.Println(estr, string(res2B), p, estr[p:])
}

func TestRecMaps() {
	var estr = []byte("d4:spaml1:a1:bee...")
	blist, p := AsBDictionary(estr)
	res2B, _ := json.Marshal(blist)
	fmt.Println(estr, string(res2B), p, estr[p:])
}

func TestTorrent() {
	raw, _ := ioutil.ReadFile("leaves.torrent")
	torrent := TorrentDecode(raw)
	fmt.Println(url.PathEscape(torrent.Tracker.InfoHash))
	torrent.PrintAsJSON()
}