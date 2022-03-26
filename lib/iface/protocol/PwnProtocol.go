package protocol

/*
各种p2p自组网协议实现接口
*/
type IPwnProtocol interface {
	SetConfig(i interface{})
	GetConfig() interface{}
	Start()
	Stop()

	LogNoti(i interface{})
	DataNoti(i interface{})

	IMessage
	// PublishMessage(i interface{}) *IPwnProtocol
	// Subscribe(i interface{}) *IPwnProtocol
}

type IMessage interface {
	ReadMsg(i interface{})
	WriteMsg(i interface{})
}

/*//////////////////////////////////////

type xxxx struct {
}

func (n *xxx)SetConfig(i interface{}){}
func (n *xxx)GetConfig() interface{}{}
func (n *xxx)Start() *IPwnProtocol{}
func (n *xxx)Stop(){}

func (n *xxx)LogNoti(i interface{}){}
func (n *xxx)DataNoti(i interface{}){}
//////////////////////////////////////*/
