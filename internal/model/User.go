package model

type User struct {
	uuid	string				//用户标识
	depot	map[string]string	//用户仓库
}
