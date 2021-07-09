package BestGoTlv

import (
  "fmt"
  "log"
)

type   BerTlvParser struct {
    TagFactory BerTagFactory
    Log        IBerTlvLogger
}
var  tagFactory BerTagFactory = NewDefaultBerTagFactory()
func NewBerTlvParser()*BerTlvParser{
    return &BerTlvParser{
      TagFactory: tagFactory,
      Log:        NewEmptyFactory(),
    }
}

func NewBerTlvParserLogFactory(aLogger IBerTlvLogger)*BerTlvParser{
  return &BerTlvParser{
    TagFactory: tagFactory,
    Log:       aLogger,
  }
}

func NewBerTlvParserFactory(aTagFactory BerTagFactory)*BerTlvParser{
    return &BerTlvParser{
      TagFactory: aTagFactory,
      Log:        NewEmptyFactory(),
    }
}

func NewBerTlvParserFactoryLog(aTagFactory BerTagFactory,aLogger IBerTlvLogger)*BerTlvParser{
  return &BerTlvParser{
    TagFactory: aTagFactory,
    Log:      aLogger,
  }
}

type  ParseResult struct {
  Tlv *BerTlv
  Offset int
}

func NewParseResult(aTlv *BerTlv,aOffset int )*ParseResult{
    return &ParseResult{
      Tlv:    aTlv,
      Offset: aOffset,
    }
}

func(B *BerTlvParser)ParseConstructed(aBuf []byte , aOffset ,aLen  int )*BerTlv {
  result:=parseWithResult(0,aBuf,aOffset,aLen)
  return  result.Tlv
}

func (B *BerTlvParser)Parse(aBuf []byte )*BerTlvs {
  return  B.ParseOffset(aBuf,0, len(aBuf))
}

func (B *BerTlvParser)ParseOffset(aBuf []byte , aOffset int , aLen int )*BerTlvs {
  var  tlvs   = make([]*BerTlv,0)
  if aLen == 0 {
    return  &BerTlvs{Tlvs: tlvs}
  }

  offset:= aOffset
  for i:=0; i<100; i++ {
    result:=  parseWithResult(0, aBuf, offset, aLen-offset)
    tlvs = append(tlvs,result.Tlv)

    if result.Offset>=aOffset+aLen  {
      break
    }
    offset = result.Offset

  }
  return  &BerTlvs{Tlvs: tlvs}
}

func ParseConstructed(aBuf []byte )*BerTlv {
  return  ParseConstructedOffset(aBuf,0,len(aBuf))
}

func ParseConstructedOffset(aBuf []byte,aOffset ,aLen int )*BerTlv {
   result:=  parseWithResult(0, aBuf, aOffset, aLen)
  return result.Tlv
}


func createLevelPadding(aLevel int )string {
  var s  string
  for i:=0; i<aLevel*4; i++ {
    s+= fmt.Sprintf("%s" , " " )
  }
  return  s
}

// parse result
func parseWithResult( aLevel int,aBuf  [] byte, aOffset int ,aLen  int ) *ParseResult{
  levelPadding:=createLevelPadding(aLevel)
  if aOffset+aLen > len(aBuf) {
      panic(fmt.Sprintf("Length is out of the range [offset=%d, len =%d, array.length=%d,level=%s",aOffset,aLen, len(aBuf),levelPadding))
  }
  // tag
  tagBytesCount:=getTagBytesCount(aBuf,aOffset)
  tag:=createTag(levelPadding,aBuf,aOffset,tagBytesCount)
  log.Printf("%s ,tag = %s, tagBytesCount=%d, tagBuf={}", levelPadding, tag, tagBytesCount)
  //length
  lengthBytesCount:=getLengthBytesCount(aBuf,aOffset+tagBytesCount)
  valueLength:=getDataLength(aBuf,aOffset + tagBytesCount)

  if tag.IsConstructed(){
    var list  = make([]*BerTlv,0)
    newList:=addChildren(aLevel, aBuf, aOffset + tagBytesCount + lengthBytesCount, levelPadding, lengthBytesCount, valueLength, list)
    resultOffset:= aOffset + tagBytesCount + lengthBytesCount + valueLength
    return &ParseResult{
      Tlv:    &BerTlv{TheTag: tag,TheList:newList},
      Offset: resultOffset,
    }
  }else{
    // value
    posStart:=aOffset+tagBytesCount+lengthBytesCount
    postEnd:=posStart+valueLength
    value:=aBuf[posStart:postEnd]
    resultOffset:= aOffset + tagBytesCount + lengthBytesCount + valueLength
    return &ParseResult{
      Tlv:    &BerTlv{TheTag: tag,TheValue:value},
      Offset: resultOffset,
    }
  }
}


func addChildren(aLevel int ,aBuf []byte ,aOffset  int ,levelPadding  string  , aDataBytesCount int ,valueLength  int ,  list []*BerTlv) []*BerTlv {
   startPosition:= aOffset
   len:= valueLength
  for startPosition < aOffset + valueLength{
        result:= parseWithResult(aLevel+1, aBuf, startPosition, len)
        list = append(list,result.Tlv)
        startPosition = result.Offset
       len           = (aOffset + valueLength) - startPosition
     }

     return  list
}

//get data length
func getDataLength(aBuf []byte ,aOffset int )int {
  var length  = int(aBuf[aOffset] & 0xff)
  if (length & 0x80) == 0x80  {
    numberOfBytes:= length & 0x7f
    if numberOfBytes>3{
      panic(fmt.Sprintf("At position %d the len is more then 3 [%d]", aOffset, numberOfBytes))
    }
    length  = 0
    for i:=aOffset+1; i<aOffset+1+int(numberOfBytes); i++ {
      length = (length * 0x100) + int(aBuf[i] & 0xff)
    }
  }
  return length
}

//get length
func getLengthBytesCount(aBuf  []byte,aOffset int )int{
  var  len = aBuf[aOffset] & 0xff
  if  (len & 0x80) == 0x80 {
    return int(1 + (len & 0x7f))
  } else {
    return 1
  }
}

func createTag(aLevelPadding string ,aBuf []byte,aOffset int ,aLength int ) BerTag {
  return tagFactory.CreateTag(aBuf,aOffset,aLength)
}


// get tag count
func getTagBytesCount(aBuf []byte ,aOffset int )int{
       if aBuf[aOffset] & 0X1F == 0X1F{
          len:=2
         for i:=aOffset+1; i<aOffset+10; i++  {
           if (aBuf[i] & 0x80) != 0x80 {
             break
           }
           len++
         }
         return len
       }else{
         return  1
       }
}
