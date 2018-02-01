package models

//上下文共享信息
type ContextInfo struct {
	//用户个人信息
	User *FUsers
	//系统配置信息
	Sets *FSets
}