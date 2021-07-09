package BestGoTlv

import (
  "encoding/hex"
  "log"
)

type BerTlvBuilder struct {
	TheTemplate     *BerTag
	TheBuffer       []byte
	ThePos          int
	TheBufferOffset int
}

func NewBerTlvBuilder() *BerTlvBuilder {
	return NewBerTlvBuilderTemplate(nil)
}

func NewBerTlvBuilderTemplate(aTemplate *BerTag) *BerTlvBuilder {
	newByte := make([]byte, DEFAULT_SIZE)
	return NewBerTlvBuilderOffset(aTemplate, newByte, 0, DEFAULT_SIZE)
}

func NewBerTlvBuilderOffset(aTemplate *BerTag, aBuffer []byte, aOffset, aLength int) *BerTlvBuilder {
	return &BerTlvBuilder{
		TheTemplate:     aTemplate,
		TheBuffer:       aBuffer,
		ThePos:          aOffset,
		TheBufferOffset: aOffset,
	}
}

func from(aTlv BerTlv) {
	return
}

func fromBufferSize(aTlv BerTlv, bufferSize int) {
	if aTlv.IsConstructed() {

	}
}

func (b *BerTlvBuilder) AddHex(aTag BerTag, aHex string) *BerTlvBuilder {
	buffer, _ := hex.DecodeString(aHex)
	return b.addBytes(aTag, buffer, 0, len(buffer))
}

func (b *BerTlvBuilder) AddText(aTag BerTag, aText string) *BerTlvBuilder {
	buffer := []byte(aText)
	return b.addBytes(aTag, buffer, 0, len(buffer))
}

func (b *BerTlvBuilder) addBytes(aTag BerTag, aBytes []byte, aFrom, aLength int) *BerTlvBuilder {
	tagLength := len(aTag.Bytes)
	lengthBytesCount := calculateBytesCountForLength(tagLength)
	b.TheBuffer = copySlice(b.ThePos,aTag.Bytes, b.TheBuffer)
  b.ThePos+=tagLength
	fillLength(b.TheBuffer, b.ThePos, aLength)
	b.ThePos += lengthBytesCount

	// VALUE
	newABytes := aBytes[aFrom : aFrom+aLength]
	b.TheBuffer = copySlice(b.ThePos,newABytes,b.TheBuffer)
	b.ThePos += aLength
	return b
}

func copySlice(index int, dest []byte , src []byte) (ns []byte) {
  ns = append(ns, src[:index]...) // 切片后加..., 相当于拆包成单个元素
  ns = append(ns, dest...)
  ns = append(ns, src[index+len(dest):]...)
  return
}


const DEFAULT_SIZE = 5 * 1024

func template(aTemplate *BerTag) *BerTlvBuilder {
	newByte := make([]byte, DEFAULT_SIZE)
	return NewBerTlvBuilderOffset(aTemplate, newByte, 0, DEFAULT_SIZE)
}

func templateBufferSize(aTemplate *BerTag, bufferSize int) *BerTlvBuilder {
	newByte := make([]byte, bufferSize)
	return NewBerTlvBuilderOffset(aTemplate, newByte, 0, bufferSize)
}

func (b *BerTlvBuilder) build() int {
	if b.TheTemplate != nil {
		tagLen := len(b.TheTemplate.Bytes)
		lengthBytesCount := calculateBytesCountForLength(b.ThePos)
		b.TheBuffer = b.TheBuffer[b.TheBufferOffset : b.TheBufferOffset+b.ThePos]
		b.TheBuffer = b.TheBuffer[tagLen+lengthBytesCount : tagLen+lengthBytesCount+b.ThePos]
		b.TheTemplate.Bytes = b.TheTemplate.Bytes[b.TheBufferOffset : b.TheBufferOffset+len(b.TheTemplate.Bytes)]

		fillLength(b.TheBuffer, tagLen, b.ThePos)
		b.ThePos += tagLen + lengthBytesCount
	}
	return b.ThePos
}

func (b *BerTlvBuilder) BuildArray() []byte {
	count := b.build()
	buf := b.TheBuffer[:count]
	return buf
}

func (b *BerTlvBuilder) buildTlv() *BerTlv {
	count := b.build()
	return NewBerTlvParser().ParseConstructed(b.TheBuffer, b.TheBufferOffset, count)
}

func (b *BerTlvBuilder) buildTlvs() *BerTlvs {
	count := b.build()
	return NewBerTlvParser().ParseOffset(b.TheBuffer, b.TheBufferOffset, count)
}

func fillLength(aBuffer []byte, aOffset int, aLength int) {
	if aLength < 0x80 {
		aBuffer[aOffset] = byte(aLength)
	} else if aLength < 0x100 {
		aBuffer[aOffset] = byte(0x81)
		aBuffer[aOffset+1] = byte(aLength)
	} else if aLength < 0x10000 {
		aBuffer[aOffset] = byte(0x82)
		aBuffer[aOffset+1] = (byte)(aLength / 0x100)
		aBuffer[aOffset+2] = (byte)(aLength % 0x100)
	} else if aLength < 0x1000000 {
		aBuffer[aOffset] = byte(0x83)
		aBuffer[aOffset+1] = (byte)(aLength / 0x10000)
		aBuffer[aOffset+2] = (byte)(aLength / 0x100)
		aBuffer[aOffset+3] = (byte)(aLength % 0x100)
	} else {
		log.Fatalf("length [%d] out of range (0x1000000)", aLength)
	}
}

func calculateBytesCountForLength(aLength int) int {
	var ret int
	if aLength < 0x80 {
		ret = 1
	} else if aLength < 0x100 {
		ret = 2
	} else if aLength < 0x10000 {
		ret = 3
	} else if aLength < 0x1000000 {
		ret = 4
	} else {
    log.Fatalf("length [%d] out of range (0x1000000)", aLength)
	}
	return ret
}
