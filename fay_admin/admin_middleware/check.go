package admin_middleware

import (
	"fay_admin/admin_handler"
	"github.com/henrylee2cn/faygo"
	"strings"
)

//检验用户是否登陆，且有访问当前url的权限是否有权限
var Check = faygo.HandlerFunc(func(ctx *faygo.Context) error {
	session := ctx.GetSession("admin")
	//如果没有，则跳转到登陆界面
	if session == nil {
		return ctx.Redirect(301, "/admin/login")
	}
	//从ssession中读取出用户的权限
	admin, ok := session.(admin_handler.UserInfo)
	if !ok {
		return nil
	}
	//检测是否有当前访问url的权限
	url := strings.TrimSpace(ctx.URL().Path)
	if _, ok = admin.Actions[url]; !ok {
		faygo.Debug(url)
		faygo.Debug(admin.Actions)
		return ctx.Redirect(301, "/admin/index")
	}
	return nil
})
