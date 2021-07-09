package BestGoTlv

import (
  "encoding/hex"
	"fmt"
  "log"
  "reflect"
)

type BerTlv struct {
	TheTag   BerTag
	TheValue []byte
	TheList  []*BerTlv
}

func NewBerTlvList(aTag BerTag, aList []*BerTlv) *BerTlv {
	return &BerTlv{
		TheTag:   aTag,
		TheValue: nil,
		TheList:  aList,
	}
}

func NewBerTlvValue(aTag BerTag, aValue []byte) *BerTlv {
	return &BerTlv{
		TheTag:   aTag,
		TheValue: aValue,
		TheList:  nil,
	}
}

func (b *BerTlv) GetTag() BerTag {
	return b.TheTag
}

func (b *BerTlv) IsPrimitive() bool {
	return !b.TheTag.IsConstructed()
}

func (b *BerTlv) IsConstructed() bool {
	return b.TheTag.IsConstructed()
}

/**
find tag
*/
func (b *BerTlv) Find(aTag BerTag) *BerTlv {
	if reflect.DeepEqual(aTag, b.GetTag()) {
		return b
	}
	if b.IsConstructed() {
		for _, tlv := range b.TheList {
			ret := tlv.Find(aTag)
			if ret != nil {
				return ret
			}
		}
	}
	return nil
}

func (b *BerTlv) FindAll(aTag BerTag) []*BerTlv {
	var list []*BerTlv
	if aTag.Equals(b.GetTag()) {
		list = append(list, b)
		return list
	} else if b.IsConstructed() {
		for _, tlv := range b.TheList {
			list = append(list, tlv.Find(aTag))
		}
	}
	return list
}

func (b *BerTlv) GetHexValue() string {
	if b.IsConstructed() {
		log.Fatalf("Tag is CONSTRUCTED =%x", b.TheTag.Bytes)
	}
	return hex.EncodeToString(b.TheValue)
}

func (b *BerTlv) GetTextValue() string {
	if b.IsConstructed() {
		log.Fatalf("TLV is constructed =%x", b.TheTag.Bytes)
		return ""
	}
	return string(b.TheValue)
}

func (b *BerTlv) GetBytesValue() []byte {
	if b.IsConstructed() {
		log.Fatalf("TLV [%s] is constructed =%x", b.TheTag, b.TheTag.Bytes)
		return nil
	}
	return b.TheValue
}

func (b *BerTlv) GetValues() []*BerTlv {
	if b.IsPrimitive() {
		log.Fatalf("  Tag [%s]is PRIMITIVE", b.TheTag)
		return nil
	}
	return b.TheList
}

func (b *BerTlv) Equals(o interface{}) bool {
	if o == nil {
		return false
	}
	if reflect.DeepEqual(b, o) {
		return true
	}
	var newBerTlv *BerTlv
	var ok bool
	if newBerTlv, ok = o.(*BerTlv); ok {
		if reflect.DeepEqual(newBerTlv.TheTag, b.TheTag) {
			return false
		}
		if reflect.DeepEqual(newBerTlv.TheValue, b.TheValue) {
			return false
		}
		if reflect.DeepEqual(newBerTlv.TheList, b.TheList) {
			return false
		}
	}
	return false
}

func (b *BerTlv) HashCode() int {
	var result int
	result = b.TheTag.HashCode()
	result = 31*result + String(string(b.TheValue))
	return result
}
func (b *BerTlv) ToString() string {
	return fmt.Sprintf("BerTlv{theTag=%s, theValue=%s ,theList=%v }", b.TheTag, string(b.TheValue), b.TheList)
}
