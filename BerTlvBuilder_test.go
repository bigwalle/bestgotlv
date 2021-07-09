package BestGoTlv

import "testing"

func TestNewBerTlvBuilder(t *testing.T) {
     bytes:= NewBerTlvBuilder().
        AddHex(NewBerTagFirstByte(0X50),"56495341").
        AddHex(NewBerTagFirstByte(0X57),"1000023100000033D44122011003400000481F").BuildArray()
  t.Logf("data.length=%x",len(bytes))
    t.Logf("data=%x",bytes)



  //mock tlv track response
  NewBerTlvBuilder().
       AddText(NewBerTagFirstByte(0x4E),"").
      AddHex(NewBerTagFirstByte(0x45),"").
      AddHex(NewBerTagFirstByte(0x4A),"").
      AddHex(NewBerTagFirstByte(0x97),"").BuildArray()
}
