package admin_handler

import (
	. "fay_admin/admin_common"
	"fay_admin/admin_model"
	"github.com/henrylee2cn/faygo"
	"time"
)

type LoginHandler struct {
	User string `param:"<in:formData><desc:用户名><name:user>"`
	Pass string `param:"<in:formData><desc:密码><name:pass>"`
}

type UserInfo struct {
	admin_model.Admin
	Rolename string
	Powers   []admin_model.Column
	Actions  map[string]admin_model.Column
}

func (l *LoginHandler) Serve(ctx *faygo.Context) error {
	if ctx.Method() == "GET" {
		return ctx.Render(200, "admin_view/login.html", nil)
	}
	err := ctx.BindForm(l)
	if err != nil {
		return jsonReturn(ctx, 0, "参数解析出错")
	}

	if !IsEmail(l.User) {
		return jsonReturn(ctx, 0, "请输入正确的用户名")
	}

	if !CheckPwd(l.Pass) {
		return jsonReturn(ctx, 0, "请输入6~16位数字和大小写字符组成的密码")
	}
	adminModel := admin_model.DefaultAdmin
	//查询数据库
	admin, err := adminModel.Check(l.User, l.Pass)
	if err != nil {
		return jsonReturn(ctx, 0, err.Error())
	}

	//修改最后一次登陆的时间
	now_time := time.Now()
	err = adminModel.UpdateLastLogin(l.User, ctx.IP(), now_time.Unix())
	if err == nil {
		admin.Login_ip = ctx.IP()
		admin.Login_time = now_time.Unix()
	}
	//获取用户的角色
	role, err := admin_model.DefaultRole.GetPowerByRid(admin.Role_id)
	if err != nil {
		return jsonReturn(ctx, 0, err.Error())
	}
	//获取角色所有的权限栏目地址
	powers, actions, err := admin_model.DefaultPower.GetColumnsByPowerIds(role.Power)
	if err != nil {
		return jsonReturn(ctx, 0, err.Error())
	}
	userInfo := UserInfo{
		Admin:    admin,
		Rolename: role.Name,
		Powers:   powers,
		Actions:  actions,
	}
	//将用户信息存在session中
	ctx.SetSession("admin", userInfo)
	return jsonReturn(ctx, 200, "登陆成功")
}

func (l *LoginHandler) Doc() faygo.Doc {
	return faygo.Doc{}
}

type QuitHandler struct {
}

func (q *QuitHandler) Serve(ctx *faygo.Context) error {
	//清空session
	ctx.DelSession("admin")
	return jsonReturn(ctx, 200, "退出成功")
}
