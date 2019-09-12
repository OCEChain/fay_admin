package admin_handler

import (
	. "fay_admin/admin_common"
	"fay_admin/admin_model"
	"github.com/henrylee2cn/faygo"
)

type IndexHandler struct {
}

func (i *IndexHandler) Serve(ctx *faygo.Context) error {
	//获取当前用的所有的栏目
	sess_admin := ctx.GetSession("admin")
	if sess_admin == nil {
		return ctx.Redirect(301, "/admin/login")
	}
	userInfo, ok := sess_admin.(UserInfo)
	if !ok {
		return ctx.HTML(0, "系统出错")
	}

	render_data := make(map[string]interface{})
	render_data["columns"] = userInfo.Powers
	render_data["admin"] = userInfo.Admin
	return ctx.Render(200, "admin_view/index.html", render_data)
}

//修改用户资料
type EditAdminHandler struct {
	Nickname string `param:"<in:formData><name:nickname>"`
	Face     string `param:"<in:formData><name:face>"`
}

func (e *EditAdminHandler) Serve(ctx *faygo.Context) error {
	//从session中取出当前用户信息
	session := ctx.GetSession("admin")
	if session == nil {
		return ctx.Redirect(301, "/admin/login")
	}
	//从ssession中读取出用户的权限
	admin, ok := session.(UserInfo)
	if !ok {
		return ctx.HTML(0, "服务器出错")
	}
	if ctx.Method() == "GET" {
		render_data := make(map[string]interface{})
		render_data["user"] = admin.Admin
		faygo.Debug(admin.Admin)
		return ctx.Render(200, "admin_view/edit_admin.html", render_data)
	}
	err := ctx.BindForm(e)
	if err != nil {
		return jsonReturn(ctx, 0, "解析参数出错")
	}
	if e.Nickname == "" {
		return jsonReturn(ctx, 0, "昵称不能为空")
	}
	if e.Face == "" {
		return jsonReturn(ctx, 0, "请上传头像")
	}

	if e.Nickname == admin.Nickname && e.Face == admin.Face {
		return jsonReturn(ctx, 200, "修改成功")
	}

	//改变用户信息
	err = admin_model.DefaultAdmin.Edit(admin.Id, e.Nickname, e.Face)
	if err != nil {
		return jsonReturn(ctx, 0, err.Error())
	}
	//更新缓存中的信息
	admin.Face = e.Face
	admin.Nickname = e.Nickname
	ctx.SetSession("admin", admin)
	return jsonReturn(ctx, 200, "修改成功")
}

//修改密码
type EditPassHandler struct {
	Oldpass string `param:"<in:formData><name:oldpass>"`
	Newpass string `param:"<in:formData><name:newpass>"`
}

func (e *EditPassHandler) Serve(ctx *faygo.Context) error {
	session := ctx.GetSession("admin")
	if session == nil {
		return ctx.Redirect(301, "/admin/login")
	}

	if ctx.Method() == "GET" {
		return ctx.Render(200, "admin_view/edit_pass.html", nil)
	}
	//从ssession中读取出用户的权限
	admin, ok := session.(UserInfo)
	if !ok {
		return ctx.HTML(0, "服务器出错")
	}
	if !CheckPwd(e.Oldpass) {
		return jsonReturn(ctx, 0, "请输入6~16位数字和大小写字符组成的原密码")
	}
	if !CheckPwd(e.Newpass) {
		return jsonReturn(ctx, 0, "请输入6~16位数字和大小写字符组成的新密码")
	}
	if admin.Pass != MakeMd5([]byte(e.Oldpass)) {
		return jsonReturn(ctx, 0, "原密码错误")
	}
	//修改密码
	err := admin_model.DefaultAdmin.EditPass(admin.Id, e.Newpass)
	if err != nil {
		return jsonReturn(ctx, 0, err.Error())
	}
	//清空session
	ctx.DelSession("admin")
	return jsonReturn(ctx, 200, "修改成功")
}

type UploaderHandler struct {
}

//上传图片
func (u *UploaderHandler) Serve(ctx *faygo.Context) error {

	saveInfo, err := ctx.SaveFile("file", true)
	if err != nil {
		return jsonReturn(ctx, 0, "上传失败")
	}

	return jsonReturn(ctx, 200, saveInfo.Url)
}
