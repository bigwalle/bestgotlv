package BestGoTlv
type IBerTlvLogger interface {
   IisDebugEnabled() bool
   Debug(aFormat string ,args ... string )
}
