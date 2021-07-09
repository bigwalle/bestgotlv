package BestGoTlv

type BerTlvs struct {
      Tlvs []*BerTlv
}


/**
获取根节点集合
 */
func (tlvs *BerTlvs)GetList() []*BerTlv {
      return  tlvs.Tlvs
}


func (tlvs *BerTlvs) FindAll(aTag BerTag)[]*BerTlv {
  list:=make([]*BerTlv,0)
  for _,tlv:=range  tlvs.Tlvs{

      tagList:=tlv.FindAll(aTag)

      for _,tagTlv:=range tagList{
        list =  append(list,tagTlv)
      }
  }
  return  list
}


func (tlvs *BerTlvs) Find(aTag BerTag)*BerTlv {
  for _,tlv:=range  tlvs.Tlvs{
    found:=tlv.Find(aTag)
    if found!=nil{
      return found
    }
  }
  return  nil
}
