package BestGoTlv


type BerTagFactory interface {
   CreateTag (aBuf []byte ,aOffset , aLength int ) BerTag
}
