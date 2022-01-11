package server

/*
@关键字 Server
*/
import (
	"context"
	"ctmgr/containermgr"
	"ctmgr/cs"
	"net"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func Init() {
	containermgr.GetInstance().Register(uint32(cs.Command_C2S_CREATE_CONTAINER_CMD), NewCreateContainerIns()) //创建容器
	containermgr.GetInstance().Register(uint32(cs.Command_C2S_STOP_CONTAINER_CMD), NewStopContainer())        //停止容器实例
	containermgr.GetInstance().Register(uint32(cs.Command_C2S_RESTART_CONTAINER_CMD), NewRestartContainer())  //从启容器实例
	containermgr.GetInstance().Register(uint32(cs.Command_C2S_REMOVE_CONTAINER_CMD), NewRemoveContainer())    //移除容器实例
}

type Server struct {
	listenAddr string
	listen     net.Listener
	grpcSvr    *grpc.Server //grpc 服务
}

func NewServer(lisAddr string) (*Server, error) {
	_ins := &Server{
		listenAddr: lisAddr,
		grpcSvr:    grpc.NewServer(),
	}
	_ins.registerGrpcServer() //注册grpc server
	return _ins, nil
}

func (this *Server) registerGrpcServer() {
	cs.RegisterContainerInstanceMgrServer(this.grpcSvr, NewContainerServer())
	reflection.Register(this.grpcSvr)
}
func (this *Server) Run() (err error) {
	if this.listen, err = net.Listen("tcp", this.listenAddr); err != nil {
		log.Errorf("[Server] Run net.Listen faild listenAddr %s err %v", this.listenAddr, err)
		return err
	}
	defer this.listen.Close()
	if err = this.grpcSvr.Serve(this.listen); err != nil {
		log.Errorf("[Serve] Run this.grpcSvr.Serve faild listenAddr %s err %v", this.listenAddr, err)
		return err
	}
	return nil
}

/*grpc server register
 */
type ContainerServer struct{}

func NewContainerServer() *ContainerServer {
	return &ContainerServer{}
}

func (this *ContainerServer) Execute(ctx context.Context, in *cs.BaseMsg) (*cs.BaseResponseMsg, error) {
	return containermgr.GetInstance().Execute(ctx, in)
}
