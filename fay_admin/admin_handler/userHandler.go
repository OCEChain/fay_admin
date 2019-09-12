package admin_handler

import (
	. "fay_admin/admin_common"
	"fay_admin/admin_model"
	"github.com/henrylee2cn/faygo"
	"strings"
)

type UserHandler struct {
	Page  int `param:"<in:formData><name:page>"`
	Limit int `param:"<in:formData><name:limit>"`
}

func (u *UserHandler) Serve(ctx *faygo.Context) error {
	//session := ctx.GetSession("admin")
	if ctx.Method() == "GET" {
		session := ctx.GetSession("admin")
		if session == nil {
			return ctx.Redirect(301, "/admin/login")
		}
		//从ssession中读取出用户的权限
		admin, ok := session.(UserInfo)
		if !ok {
			return ctx.HTML(0, "服务器出错")
		}
		render_data := make(map[string]interface{})
		render_data["deluser"] = ""
		render_data["add"] = ""
		render_data["edit"] = ""
		if col, ok := admin.Actions["/admin/user_add"]; ok {

			render_data["add"] = col
		}
		if col, ok := admin.Actions["/admin/user_del"]; ok {
			render_data["deluser"] = col
		}
		if col, ok := admin.Actions["/admin/user_setstatus"]; ok {
			render_data["edit"] = col
		}

		return ctx.Render(200, "admin_view/user.html", render_data)
	}
	//获取所有的用户列表
	err := ctx.BindForm(u)
	if err != nil {
		return jsonReturn(ctx, -1, "参数解析出错")
	}
	if u.Page <= 0 {
		u.Page = 1
	}
	if u.Limit <= 0 {
		u.Limit = 10
	}
	adminModel := admin_model.DefaultAdmin

	count, err := adminModel.Count()
	if err != nil {
		return jsonReturn(ctx, -1, "渲染失败")
	}

	var list []*admin_model.User
	if count > 0 {
		offset := (u.Page - 1) * u.Limit
		list, err = adminModel.List(offset, u.Limit)
		if err != nil {
			return jsonReturn(ctx, -1, "渲染失败")
		}
	}
	return jsonReturn(ctx, 0, list, count)
}

type UserAddHandler struct {
	User   string `param:"<in:formData><name:user>"`
	Pass   string `param:"<in:formData><name:pass>"`
	Role   int    `param:"<in:formData><name:role>"`
	Status int    `param:"<in:formData><name:status>"`
}

func (u *UserAddHandler) Serve(ctx *faygo.Context) error {
	if ctx.Method() == "GET" {
		//获取所有的角色
		roleModel := admin_model.DefaultRole
		role_list, err := roleModel.GetRole()
		if err != nil {
			return ctx.HTML(0, "渲染失败")
		}
		render_data := make(map[string]interface{})
		render_data["role_list"] = role_list
		return ctx.Render(200, "admin_view/user_add.html", render_data)
	}
	err := ctx.BindForm(u)
	if err != nil {
		return jsonReturn(ctx, 0, "参数解析失败")
	}
	u.User = strings.TrimSpace(u.User)
	u.Pass = strings.TrimSpace(u.Pass)
	if !IsEmail(u.User) {
		return jsonReturn(ctx, 0, "请输入邮箱格式的账号")
	}
	if !CheckPwd(u.Pass) {
		return jsonReturn(ctx, 0, "请输入6～16位数字，字母组成的密码")
	}
	if u.Role == 0 {
		return jsonReturn(ctx, 0, "请选择角色")
	}
	//判断是否是合法的角色id
	role_list, err := admin_model.DefaultRole.GetMapRole()
	if err != nil {
		return jsonReturn(ctx, 0, err.Error())
	}
	if _, ok := role_list[u.Role]; !ok {
		return jsonReturn(ctx, 0, "不存在的角色id")
	}
	if u.Status != 0 && u.Status != 1 {
		return jsonReturn(ctx, 0, "非法参数")
	}
	//数据入库
	err = admin_model.DefaultAdmin.Add(u.User, u.Pass, u.Role, u.Status)
	if err != nil {
		return jsonReturn(ctx, 0, err.Error())
	}
	return jsonReturn(ctx, 200, "添加成功")
}

type UserDelHandler struct {
	Id int `param:"<in:formData><required><name:id>"`
}

func (u *UserDelHandler) Serve(ctx *faygo.Context) error {
	err := ctx.BindForm(u)
	if err != nil {
		return jsonReturn(ctx, 0, "参数解析失败")
	}
	if u.Id <= 0 {
		return jsonReturn(ctx, 0, "非法参数")
	}

	//判断是否是删除自己
	sess_admin := ctx.GetSession("admin")
	if sess_admin == nil {
		return ctx.Redirect(301, "/admin/login")
	}
	userInfo, ok := sess_admin.(UserInfo)
	if !ok {
		return ctx.HTML(0, "系统出错")
	}
	if userInfo.Id == u.Id {
		return jsonReturn(ctx, 0, "不能删除自己")
	}
	//获取用户信息
	admin, err := admin_model.DefaultAdmin.GetAdminById(u.Id)
	if err != nil {
		return jsonReturn(ctx, 0, err.Error())
	}
	if admin.Role_id == 1 {
		return jsonReturn(ctx, 0, "无法删除超级管理员")
	}
	//执行删除
	err = admin_model.DefaultAdmin.Del(u.Id)
	if err != nil {
		return jsonReturn(ctx, 0, err.Error())
	}
	return jsonReturn(ctx, 200, "删除成功")
}

type UserSetStatus struct {
	Id     int `param:"<in:formData><required><name:id>"`
	Status int `param:"<in:formData><required><name:status>"`
}

func (u *UserSetStatus) Serve(ctx *faygo.Context) error {
	err := ctx.BindForm(u)
	if err != nil {
		return jsonReturn(ctx, 0, "参数解析失败")
	}
	if u.Id <= 0 || u.Status < 0 || u.Status > 1 {
		return jsonReturn(ctx, 0, "非法参数")
	}
	err = admin_model.DefaultAdmin.SetStatus(u.Id, u.Status)
	if err != nil {
		return jsonReturn(ctx, 0, err.Error())
	}
	return jsonReturn(ctx, 200, "修改成功")
}
