package BestGoTlv


type DefaultBerTagFactory struct {

}

func NewDefaultBerTagFactory() BerTagFactory {
    return   &DefaultBerTagFactory{}
}

func (*DefaultBerTagFactory)CreateTag (aBuf []byte ,aOffset , aLength int ) BerTag {
    return  NewBerTagOffSet(aBuf,aOffset,aLength)
}
