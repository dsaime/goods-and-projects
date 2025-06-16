package goodsEvent

//const (
//	GoodsEventCreated = "goods_created"
//	GoodsEventUpdated = "goods_updated"
//)

type Logger interface {
	Log(event Event)
}
