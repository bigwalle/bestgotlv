package BestGoTlv

type  EmptyLogFactory struct {

}


func NewEmptyFactory ()*EmptyLogFactory{
    return   &EmptyLogFactory{}
}
func (EmptyLogFactory)IisDebugEnabled() bool{
  return   true
}
func (EmptyLogFactory)Debug(aFormat string ,args ... string ){
  return
}
