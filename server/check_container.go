package server

import "ctmgr/utils"

/*检查容器实例是否存在
*@param1 srcContainerName string
*@return1 err error 状态
*@return2 status bool 是否存在
 */
func CheckContainerExist(srcContainerName string) (error, bool) {
	err, containerNames := utils.GetAllContainerInsNames()
	if err != nil {
		return err, false
	}
	for _, name := range containerNames {
		if srcContainerName == name {
			return nil, true
		}
	}
	return nil, false
}
