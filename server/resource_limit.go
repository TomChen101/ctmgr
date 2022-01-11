package server

/*
@关键字 ResourceLimit
*/
import (
	"ctmgr/cs"

	"google.golang.org/protobuf/proto"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	log "github.com/sirupsen/logrus"
)

/*容器启动资源限制 设置容器启动配置参数
 */
type ResourceLimit struct {
	Config           *container.Config
	HostConfig       *container.HostConfig
	NetworkingConfig *network.NetworkingConfig
	ImageName        string
}

/*创建容器启动配置
*@param1  datas []byte pb数据
*@return1 obj *ResourceLimit
 */
func NewResourceLimit(datas []byte) (*ResourceLimit, error) {
	csMsg := &cs.C2S_CreateContainer{}
	err := proto.Unmarshal(datas, csMsg)
	if err != nil {
		log.Errorf("[ResourceLimit] proto.Unmarshal faild err %v", err)
		return nil, err
	}
	log.Infof("[ResourceLimit] proto.Marshal csMsg %v", csMsg)
	msg := &ResourceLimit{
		ImageName: csMsg.GetImageName(),
	}
	msg.Config = &container.Config{
		Tty:   true,
		Image: csMsg.GetImage(),
	}
	msg.HostConfig = &container.HostConfig{}
	msg.HostConfig.Memory = int64(csMsg.GetMemoryLimit())
	msg.HostConfig.CPUShares = int64(csMsg.GetCpuCoreLimit())
	msg.NetworkingConfig = &network.NetworkingConfig{}
	return msg, nil
}
