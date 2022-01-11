package server

/*
@关键字 ExecuteCommand
*/

import (
	"context"
	"ctmgr/cs"
	"errors"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
)

/*=========================================创建容器实例=========================================*/
type CreateContainerIns struct {
	resourceCfg *ResourceLimit
}

func NewCreateContainerIns() *CreateContainerIns {
	return &CreateContainerIns{}
}

/*启动容器实例
@param1 ctx context.Context 上下文
@param2 responseMsg *cs.BaseResponseMsg
@param3 datas []bytes pb 序列化消息
@return1 err error 状态
*/
func (this *CreateContainerIns) LaunchCommand(ctx context.Context, responseMsg *cs.BaseResponseMsg, datas []byte) error {
	responseMsg.Cmd = uint32(cs.Command_S2C_CREATE_CONTAINER_CMD)
	var err error
	this.resourceCfg, err = NewResourceLimit(datas)
	if err != nil {
		log.Errorf("[ExecuteCommand] CreateContainerIns::LaunchCommand faild err %v", err)
		responseMsg.Ret = uint32(cs.RET_RET_ERROR)
		return err
	}
	if err, exist := CheckContainerExist(this.resourceCfg.ImageName); err != nil || exist {
		log.Errorf("[ExecuteCommand] CreateContainerIns::LaunchCommand faild err %v exist %v src container name %s", err, exist, this.resourceCfg.ImageName)
		responseMsg.Ret = uint32(cs.RET_RET_CONTAINER_IS_EXIST)
		return errors.New("check contianer faild")
	}
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Errorf("[ExecuteCommand] CreateContainerIns::LaunchCommand client.NewClientWithOpts faild err %v", err)
		responseMsg.Ret = uint32(cs.RET_RET_ERROR)
		return err
	}
	resp, err := cli.ContainerCreate(ctx, this.resourceCfg.Config, this.resourceCfg.HostConfig, nil, nil, this.resourceCfg.ImageName)
	if err != nil {
		log.Errorf("[ExecuteCommand] CreateContainerIns:LaunchCommand cli.ContainerCreate faild err %v", err)
		responseMsg.Ret = uint32(cs.RET_RET_ERROR)
		return err
	}
	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		log.Errorf("[ExecuteCommand] CreateContainerIns::LaunchCommand  cli.ContainerStart faild err %v", err)
		responseMsg.Ret = uint32(cs.RET_RET_START_CONTAINER_FAILD)
		return err
	}
	msg := &cs.S2C_CreateContainer{
		ContainerId: resp.ID,
	}
	body, err := proto.Marshal(msg)
	if err != nil {
		log.Errorf("[ExecuteCommand] CreateContainerIns::LaunchCommand proto.Unmarshal faild err %v ", err)
		responseMsg.Ret = uint32(cs.RET_RET_ERROR)
		return err
	}
	responseMsg.Body = body
	log.Infof("[ExecuteCommand] CreateContainerIns::LaunchCommand stat container success id %v", resp.ID)
	return nil
}

/*==============================================关闭容器实例=======================================*/
type StopContainer struct {
}

func NewStopContainer() *StopContainer {
	return &StopContainer{}
}

func (this *StopContainer) LaunchCommand(ctx context.Context, responseMsg *cs.BaseResponseMsg, datas []byte) error {
	responseMsg.Cmd = uint32(cs.Command_S2C_STOP_CONTAINER_CMD)
	csMsg := &cs.C2S_StopContainer{}
	if err := proto.Unmarshal(datas, csMsg); err != nil {
		log.Errorf("[ExecuteCommand] StopContainer::LaunchCommand proto.Unmarshal faild err %v", err)
		responseMsg.Ret = uint32(cs.RET_RET_ERROR)
		return err
	}
	if err, exist := CheckContainerExist(csMsg.ImageName); err != nil || !exist {
		log.Errorf("[ExecuteCommand] StopContainer::LaunchCommand CheckContainerExist faild err %v exist %v src container name %s", err, exist, csMsg.ImageName)
		responseMsg.Ret = uint32(cs.RET_RET_ERROR)
		return errors.New("check contianer faild")
	}
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Errorf("[ExecuteCommand] StopContainer::LaunchCommand client.NewClientWithOpts faild err %v", err)
		responseMsg.Ret = uint32(cs.RET_RET_ERROR)
		return err
	}
	timeout := time.Duration(10000)
	if err := cli.ContainerStop(ctx, csMsg.GetImageName(), &timeout); err != nil {
		log.Errorf("[ExecuteCommand] StopContainer::LaunchCommand cli.ContainerStop faild err %v coontainer image name %s", err, csMsg.GetImageName())
		responseMsg.Ret = uint32(cs.RET_RET_STOP_CONTAINER_FAILD)
		return err
	}
	msg := &cs.S2C_StopContainer{
		ImageName: csMsg.GetImageName(),
	}
	body, err := proto.Marshal(msg)
	if err != nil {
		log.Errorf("[ExecuteCommand] StopContainer::LaunchCommand proto.Unmarshal faild err %v ", err)
		responseMsg.Ret = uint32(cs.RET_RET_STOP_CONTAINER_FAILD)
		return err
	}
	responseMsg.Body = body
	log.Infof("[ExecuteCommand] StopContainer::LaunchCommand stop container %s success", csMsg.GetImageName())
	return nil
}

/*==========================================从启容器实例=================================*/
type RestartContainer struct {
}

func NewRestartContainer() *RestartContainer {
	return &RestartContainer{}
}

/*重启容器实例
@param1 ctx context.Context 上下文
@param2 responseMsg *cs.BaseResponseMsg
@param3 datas []byte
@return err error 状态
*/
func (this *RestartContainer) LaunchCommand(ctx context.Context, responseMsg *cs.BaseResponseMsg, datas []byte) error {
	responseMsg.Cmd = uint32(cs.Command_S2C_RESTART_CONTAINER_CMD)
	csMsg := &cs.C2S_RestartContainer{}
	if err := proto.Unmarshal(datas, csMsg); err != nil {
		log.Errorf("[ExecuteCommand] RestartContainer::LaunchCommand proto.Unmarshal faild err %v", err)
		return err
	}
	if err, exist := CheckContainerExist(csMsg.ImageName); err != nil || exist {
		log.Errorf("[ExecuteCommand] RestartContainer::LaunchCommand CheckContainerExist faild err %v exist %v src container name %s", err, exist, csMsg.ImageName)
		responseMsg.Ret = uint32(cs.RET_RET_ERROR)
		return errors.New("check contianer faild")
	}
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Errorf("[ExecuteCommand] RestartContainer::LaunchCommand client.NewClientWithOpts faild err %v", err)
		responseMsg.Ret = uint32(cs.RET_RET_ERROR)
		return err
	}
	timeout := time.Duration(10000)
	if err := cli.ContainerRestart(ctx, csMsg.ImageName, &timeout); err != nil {
		log.Errorf("[ExecuteCommand] RestartContainer::LaunchCommand cli.ContainerRestart faild err %v", err)
		responseMsg.Ret = uint32(cs.RET_RET_RESTART_CONTAINER_FAILD)
		return err
	}
	msg := &cs.S2C_RestartContainer{
		ImageName: csMsg.ImageName,
	}
	body, err := proto.Marshal(msg)
	if err != nil {
		log.Errorf("[ExecuteCommand] RestartContainer::LaunchCommand proto.Marshal daild err %v", err)
		responseMsg.Ret = uint32(cs.RET_RET_RESTART_CONTAINER_FAILD)
		return err
	}
	responseMsg.Body = body
	log.Infof("[ExecuteCommand] RestartContainer::LaunchCommand success")
	return nil
}

/*============================================移除容器实例===================================*/
type RemoveContainer struct {
}

func NewRemoveContainer() *RemoveContainer {
	return &RemoveContainer{}
}

/*移除容器实例
@param1 ctx context.Context 上下文
@param2 responseMsg *cs.BaseResponseMsg
@param3 datas []byte
@return err error 状态
*/
func (this *RemoveContainer) LaunchCommand(ctx context.Context, responseMsg *cs.BaseResponseMsg, datas []byte) error {
	responseMsg.Cmd = uint32(cs.Command_S2C_REMOVE_CONTAINER_CMD)
	csMsg := &cs.C2S_RestartContainer{}
	if err := proto.Unmarshal(datas, csMsg); err != nil {
		log.Errorf("[ExecuteCommand] RemoveContainer::LaunchCommand proto.Unmarshal faild err %v", err)
		responseMsg.Ret = uint32(cs.RET_RET_ERROR)
		return err
	}
	if err, exist := CheckContainerExist(csMsg.ImageName); err != nil || !exist {
		log.Errorf("[ExecuteCommand] RemoveContainer::LaunchCommand CheckContainerExist faild err %v exist %v src container name %s", err, exist, csMsg.ImageName)
		responseMsg.Ret = uint32(cs.RET_RET_ERROR)
		return errors.New("check contianer faild")
	}
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Errorf("[ExecuteCommand] RemoveContainer::LaunchCommand client.NewClientWithOpts faild err %v", err)
		responseMsg.Ret = uint32(cs.RET_RET_ERROR)
		return err
	}
	if err := cli.ContainerRemove(ctx, csMsg.ImageName, types.ContainerRemoveOptions{Force: true}); err != nil {
		log.Errorf("[ExecuteCommand] RemoveContainer::LaunchCommand cli.ContainerRemove faild err %v", err)
		responseMsg.Ret = uint32(cs.RET_RET_REMOVE_CONTAINER_FAILD)
		return err
	}
	msg := &cs.S2C_RemoveContainer{
		ImageName: csMsg.ImageName,
	}
	body, err := proto.Marshal(msg)
	if err != nil {
		log.Errorf("[ExecuteCommand] RemoveContainer::LaunchCommand proto.Marshal faild err %v", err)
		responseMsg.Ret = uint32(cs.RET_RET_REMOVE_CONTAINER_FAILD)
		return err
	}
	responseMsg.Body = body
	log.Infof("[ExecuteCommand] RemoveContainer::LaunchCommand success")
	return nil
}
