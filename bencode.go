package torrent

import (
	"strconv"
	"fmt"
	"log"
)

func decodeString(raw *[]byte, obj *BString) int {
	var (
		size string
	)

	for _, b := range *raw {
		//Add numbers to the string until ":",
		//then begin the atoi conversion.
		if b != ':' {
			size += string(b)
		} else {
			obj.size, _ = strconv.Atoi(size)
			break
		}
	}
	p := len(size) + 1
	obj.string = (*raw)[p:p+obj.size]
	return p+obj.size
}

func decodeInt(raw *[]byte, obj *BInteger) int {
	var (
		buffer string
		count int
	)

	for i, b := range *raw {
		count++
		if i == 0 {
			//Prevent false positive,
			//Must be implemented as sub if to let the first
			//conditional always handle the first letter.
			if b != 'i' {
				break
			}
		} else {
			if b == 'e' {
				obj.value, _ = strconv.Atoi(buffer)
				break
			}
			buffer += string(b)
		}
	}
	p := count
	return p
}

func decodeList(b []byte) (*BList, int) {
	var (
		a     BList
		count int
	)
	if b[0] != 'l' {
		return nil, 0
	}
	b = b[1:]
	count++
	for len(b) > 0 {
		var (
			obj interface{}
			p int
		)
		switch GetType(b[0]) {
		case BSTR:
			obj, p = AsBString(&b)
			a = append(a, obj)
		case BINT:
			obj, p = AsBInt(&b)
			a = append(a, obj)
		case BLIST:
			obj, p = AsBList(b)
			a = append(a, obj)
		case BDIR:
			obj, p = AsBDictionary(b)
			a = append(a, obj)
		case BEND:
			b = []byte{} //End, exit for
		default:
			log.Println(fmt.Sprint("unknown output ", b))
			b = []byte{}
		}
		b = b[p:]
		count += p
	}
	return &a, count + 1
}

func decodeDic(b []byte) (*BDictionary, int) {
	var (
		m = BDictionary{dictionary:make(Map)}
		count int
		key string
	)
	if b[0] != 'd' {
		return nil, 0
	}
	b = b[1:]
	count++
	for len(b) > 0 {
		var (
			obj interface{}
			p int
		)
		switch GetType(b[0]) {
		case BSTR:
			if key == "" { //No key defined, Set as key
				obj, p = AsBString(&b)
				key = obj.(*BString).String()
			} else { //Already defined, Set as value
				obj, p = AsBString(&b)
				m.keys = append(m.keys, key)
				m.dictionary[key] = obj
				key = ""
			}
		case BINT:
			obj, p = AsBInt(&b)
			m.keys = append(m.keys, key)
			m.dictionary[key] = obj
			key = ""
		case BLIST:
			obj, p = AsBList(b)
			m.keys = append(m.keys, key)
			m.dictionary[key] = obj
			key = ""
		case BDIR:
			obj, p = AsBDictionary(b)
			m.keys = append(m.keys, key)
			m.dictionary[key] = obj
			key = ""
		case BEND:
			b = []byte{} //End, exit for
		default:
			log.Println(fmt.Sprint("unknown output ", b))
			b = []byte{}
		}
		b = b[p:]
		count += p
	}
	return &m, count + 1
}

func encodeList(obj *BList) string {
	out := "l"
	for _, elem := range *obj {
		if selem, ok := elem.(*BString); ok {
			out += selem.BEncodedString()
		}
		if selem, ok := elem.(*BInteger); ok {
			out += selem.BEncodedInteger()
		}
		if selem, ok := elem.(*BList); ok {
			out += selem.BEncodedList()
		}
		if selem, ok := elem.(*BDictionary); ok {
			out += selem.BEncodedDictionary()
		}
	}
	return out + "e"
}

func encodeDic(obj *BDictionary) string {
	out := "d"
	for _, key := range obj.keys {
		elem := obj.dictionary[key]
		//SetKey
		def := NewBString([]byte(key))
		out += def.BEncodedString()
		//SetValue
		if selem, ok := elem.(*BString); ok {
			out += selem.BEncodedString()
		}
		if selem, ok := elem.(*BInteger); ok {
			out += selem.BEncodedInteger()
		}
		if selem, ok := elem.(*BList); ok {
			out += selem.BEncodedList()
		}
		if selem, ok := elem.(*BDictionary); ok {
			out += selem.BEncodedDictionary()
		}
	}
	return out + "e"
}

func createList(b []interface{}) *BList {
	var list BList
	for _, elem := range b {
		if sub, ok := elem.([]byte); ok {
			list = append(list, NewBString(sub))
		} else if sub, ok := elem.(int); ok {
			list = append(list, NewBInt(sub))
		} else if sub, ok := elem.([]interface{}); ok {
			list = append(list, NewBList(sub))
		} else if sub, ok := elem.(map[string]interface{}); ok {
			list = append(list, NewBDictionary(sub))
		} else {
			log.Print("error, invalid interface type in slice")
			return nil
		}
	}
	return &list
}

func createDic(b map[string]interface{}) *BDictionary {
	var list = BDictionary{dictionary:make(Map)}
	for key, elem := range b {
		if sub, ok := elem.([]byte); ok {
			list.keys = append(list.keys, key)
			list.dictionary[key] = NewBString(sub)
		} else if sub, ok := elem.(int); ok {
			list.keys = append(list.keys, key)
			list.dictionary[key] = NewBInt(sub)
		} else if sub, ok := elem.([]interface{}); ok {
			list.keys = append(list.keys, key)
			list.dictionary[key] = NewBList(sub)
		} else if sub, ok := elem.(map[string]interface{}); ok {
			list.keys = append(list.keys, key)
			list.dictionary[key] = NewBDictionary(sub)
		} else {
			log.Print("error, invalid interface type in map")
			return nil
		}
	}
	return &list
}