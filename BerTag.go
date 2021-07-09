package BestGoTlv

import (
  "encoding/hex"
  "reflect"
)

type BerTag struct {
    Bytes []byte
}

/**
 * Creates a new tag from given byte array. Similar {@link BerTag#BerTag(byte[], int, int)} but using
 * the full array.
 *
 * @param aBuf to create the tag
 */
func NewBerTag(aBuf [] byte)BerTag{
    return NewBerTagOffSet(aBuf,0,len(aBuf))
}

func NewBerTagOffSet(aBuf []byte , aOffset int ,aLength int )BerTag{
  bytes:=aBuf[aOffset:aOffset+aLength]
  return  BerTag{Bytes:bytes}
}

func NewBerTagFirstSecondByte(aFirstByte , aSecondByte int)*BerTag{
      return &BerTag{Bytes: []byte{byte(aFirstByte),byte(aSecondByte)}}
}

func NewBerTagFirstSecondThirdByte(aFirstByte , aSecondByte ,aFirth int)*BerTag{
  return &BerTag{Bytes: []byte{byte(aFirstByte),byte(aSecondByte),byte(aFirth)}}
}

func NewBerTagFirstByte(aFirstByte int)BerTag{
  return BerTag{Bytes: []byte{byte(aFirstByte)}}
}

func (b *BerTag) IsConstructed() bool{
    return (b.Bytes[0]& 0x20) !=0
}

func (b *BerTag) Equals(o interface{}) bool{
    if  o==nil {return  false}
    if newBer,ok:=o.(BerTag);ok {
            return reflect.DeepEqual(b ,newBer)
    }
    return false
}

func (b *BerTag) HashCode() int {
    return String(string(b.Bytes))
}


func (b *BerTag)  ToString()string{
   if  b.IsConstructed(){
       return  "+ " +hex.EncodeToString(b.Bytes)
   }else{
     return  "- " +hex.EncodeToString(b.Bytes)
   }
}
