package torrent

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
	"crypto/sha1"
	"net/url"
)

//Contains all the info about the torrent file
type TorrentFile struct {
	Length			int				`json:"length"`
	Path			[]string		`json:"path"`
}

//Contains all the meta-info data from the original torrent file
type TorrentInfo struct {
	Length			int				`json:"length"`
	Name			string			`json:"name"`
	PieceLength		int				`json:"piece length"`
	Files			[]TorrentFile	`json:"files"`
	Pieces			[]string		`json:"pieces"`
}

//Contains all the meta-file data from the original torrent file
type TorrentMeta struct {
	AnnounceList	[]string		`json:"announce-list"`
	CreationDate 	time.Time		`json:"creation date"`
	Encoding		string			`json:"encoding"`
	Comment			string			`json:"comment"`
	CreatedBy		string			`json:"created by"`
	Info			TorrentInfo		`json:"info"`
	Tracker			TrackerData		`json:"-"`
}

//Non Torrent file parameter
//Used to help the tracker data resolution
//Example: c3 49 9c 27 29 73 0a 7f 80 7e fb 86 76 a9 2d cb 6f 8a 3f 8f
//IMPORTANT: Might need to be scaped into url before use
type TrackerData	struct {
	InfoHash		string			`json:"info_hash"`
}

const SHA1Size = 40

//Decode byte-array of torrent file into TorrentMeta struct
func TorrentDecode(b []byte) *TorrentMeta {
	obj := RawDecode(b)

	var root TorrentMeta
	/*Start Getting Parent Torrent Data*/
	if belem, ok := obj.dictionary["announce"]; ok {
		root.AnnounceList = append(root.AnnounceList, belem.(*BString).String())
	}
	if belem, ok := obj.dictionary["announce-list"]; ok {
		for _, elem := range *belem.(*BList) {
			value := (*elem.(*BList))[0].(*BString).String()
			if value != root.AnnounceList[0] { //Prevent duplicated with Announce_1
				root.AnnounceList = append(root.AnnounceList, value)
			}
		}
	}
	if belem, ok := obj.dictionary["creation date"]; ok {
		root.CreationDate = time.Unix(int64(belem.(*BInteger).value/1000), 0)
	}
	if belem, ok := obj.dictionary["encoding"]; ok {
		root.Encoding = belem.(*BString).String()
	}
	if belem, ok := obj.dictionary["comment"]; ok {
		root.Comment = belem.(*BString).String()
	}
	if belem, ok := obj.dictionary["created by"]; ok {
		root.CreatedBy = belem.(*BString).String()
	}
	/*End Getting Parent Torrent Data*/

	sobj := obj.dictionary["info"].(*BDictionary)

	var info TorrentInfo
	/*Start Getting Info Torrent Data*/
	if belem, ok := sobj.dictionary["length"]; ok {
		info.Length = belem.(*BInteger).value
	}
	if belem, ok := sobj.dictionary["name"]; ok {
		info.Name = belem.(*BString).String()
	}
	if belem, ok := sobj.dictionary["piece length"]; ok {
		info.PieceLength = belem.(*BInteger).value
	}
	if belem, ok := sobj.dictionary["pieces"]; ok {
		bf := fmt.Sprintf("%x", belem.(*BString).ByteArray())
		for i := 0; i < len(bf); i += SHA1Size {
			info.Pieces = append(info.Pieces, bf[i:i+SHA1Size])
		}
	}
	/*Start Getting files Torrent Data*/
	if belem, ok := sobj.dictionary["files"]; ok {
		var files []TorrentFile
		for _, bfile := range *belem.(*BList) {
			var file TorrentFile
			if belem, ok := bfile.(*BDictionary).dictionary["length"]; ok {
				file.Length = belem.(*BInteger).value
			}
			if belem, ok := bfile.(*BDictionary).dictionary["path"]; ok {
				var path []string
				for _, bpath := range *belem.(*BList) {
					path = append(path, bpath.(*BString).String())
				}
				file.Path = path
				files = append(files, file)
			}
		}
		info.Files = files
	}
	/*End Getting files Torrent Data*/
	/*End Getting Info Torrent Data*/

	h := sha1.New()
	fmt.Fprint(h, sobj.BEncodedDictionary())
	root.Tracker = TrackerData{
		InfoHash: url.PathEscape(fmt.Sprintf("%s", h.Sum(nil))),
	}

	root.Info = info
	return &root
}

//Decode byte-array into BEncode Dictionary Interface(BDictionary)
func RawDecode(b []byte) *BDictionary {
	obj, _ := AsBDictionary(b)
	return obj
}

//Prints BDictionary as JSON removing non-value objects
func (r *BDictionary)PrintAsJSON() {
	res2B, err := json.Marshal(r)
	if err != nil {
		log.Print(err)
	}
	fmt.Println(string(res2B))
}

//Prints TorrentMeta as JSON ordered.
//IMPORTANT: Generated tracker data will be obviated.
func (r *TorrentMeta)PrintAsJSON() {
	res2B, err := json.Marshal(r)
	if err != nil {
		log.Print(err)
	}
	fmt.Println(string(res2B))
}