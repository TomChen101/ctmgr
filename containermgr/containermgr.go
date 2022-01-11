package containermgr

/*
@关键字 ContainerMgr
*/
import (
	"context"
	"ctmgr/cs"
	"errors"

	log "github.com/sirupsen/logrus"
)

var g_Mgr *ContainerMgr = nil

type CommonBase interface {
	LaunchCommand(context.Context, *cs.BaseResponseMsg, []byte) error
}

func GetInstance() *ContainerMgr {
	if g_Mgr == nil {
		panic("ContainerMgr's g_Mgr is not init")
	}
	return g_Mgr
}

type ContainerMgr struct {
	cmdMap map[uint32]CommonBase
}

func CreateInstance() error {
	if g_Mgr != nil {
		return errors.New("ContainerMgr's g_Mgr had init")
	}
	g_Mgr = &ContainerMgr{
		cmdMap: make(map[uint32]CommonBase),
	}
	return nil
}
func (this *ContainerMgr) Register(cmd uint32, obj CommonBase) {
	this.cmdMap[cmd] = obj
}
func (this *ContainerMgr) Execute(ctx context.Context, req *cs.BaseMsg) (*cs.BaseResponseMsg, error) {
	cmd := req.GetCmd()
	obj, exist := this.cmdMap[cmd]
	if !exist {
		log.Errorf("[ContainerMgr] Execute cmd %d not found", cmd)
		return nil, errors.New("cmd not found")
	}
	responseMsg := &cs.BaseResponseMsg{
		Ret: uint32(cs.RET_RET_OK),
	}
	err := obj.LaunchCommand(ctx, responseMsg, req.GetDatas())
	if err != nil {
		log.Errorf("[ContainerMgr] Execute cmd %d LaunchCommand faild err %v", cmd, err)
		return responseMsg, err
	}
	log.Infof("[ContainerMgr] Execute responseMsg %+v", responseMsg)
	return responseMsg, nil
}
