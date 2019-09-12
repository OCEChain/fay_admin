package admin_handler

import (
	"fay_admin/admin_model"
	"fmt"
	"github.com/henrylee2cn/faygo"
	"strconv"
	"strings"
)

type RoleHandler struct {
}

func (r *RoleHandler) Serve(ctx *faygo.Context) error {
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
		render_data["delrole"] = ""
		render_data["add"] = ""
		render_data["edit"] = ""
		render_data["modify_power"] = ""
		if col, ok := admin.Actions["/admin/role_add"]; ok {

			render_data["add"] = col
		}
		if col, ok := admin.Actions["/admin/role_del"]; ok {
			render_data["delrole"] = col
		}
		if col, ok := admin.Actions["/admin/role_edit"]; ok {
			render_data["edit"] = col
		}
		if col, ok := admin.Actions["/admin/role_power_modify"]; ok {
			render_data["modify_power"] = col
		}

		return ctx.Render(200, "admin_view/role.html", render_data)
	}
	//获取所有的用户列表
	page := ctx.QueryParam("page")
	if page == "" {
		page = "1"
	}
	p, err := strconv.Atoi(page)
	if err != nil || p <= 0 {
		p = 1
	}
	roleModel := admin_model.DefaultRole
	list, err := roleModel.GetRole()
	if err != nil {
		return jsonReturn(ctx, -1, "渲染失败")
	}

	return jsonReturn(ctx, 0, list)
}

type RoleAddHandler struct {
	Rolename string `param:"<in:formData><required><name:rolename>"`
}

func (r *RoleAddHandler) Serve(ctx *faygo.Context) error {
	if err := ctx.BindForm(r); err != nil {
		return jsonReturn(ctx, 0, "参数解析出错")
	}
	roleModel := admin_model.DefaultRole
	err := roleModel.Add(r.Rolename)
	if err != nil {
		return jsonReturn(ctx, 0, err.Error())
	}
	return jsonReturn(ctx, 200, "添加成功")
}

type RoleEditHandler struct {
	Id       int    `param:"<in:formData><required><name:id>"`
	Rolename string `param:"<in:formData><name:rolename>"`
	Status   int    `param:"<in:formData><name:status>"`
	Typeid   int    `param:"<in:formData><required><name:typeid>"`
}

func (r *RoleEditHandler) Serve(ctx *faygo.Context) error {
	var err error
	if err := ctx.BindForm(r); err != nil {
		return jsonReturn(ctx, 0, "参数解析出错")
	}
	//超级管理员无法修改
	if r.Id <= 1 {
		return jsonReturn(ctx, 0, "非法参数")
	}
	roleModel := admin_model.DefaultRole
	switch r.Typeid {
	case 1:
		if r.Rolename == "" {
			return jsonReturn(ctx, 0, "角色名不能为空")
		}
		err = roleModel.EditRoleName(r.Id, r.Rolename)
	case 2:
		if r.Status < 0 || r.Status > 1 {
			return jsonReturn(ctx, 0, "非法参数")
		}
		err = roleModel.EditRoleStatus(r.Id, r.Status)
	default:
		return jsonReturn(ctx, 0, "非法参数")
	}

	if err != nil {
		return jsonReturn(ctx, 0, err.Error())
	}
	return jsonReturn(ctx, 200, "修改成功")
}

type RoleDelHandler struct {
	Id int `param:"<in:formData><required><name:id>"`
}

func (r *RoleDelHandler) Serve(ctx *faygo.Context) error {
	var err error
	if err := ctx.BindForm(r); err != nil {
		return jsonReturn(ctx, 0, "参数解析出错")
	}
	//超级管理员无法修改
	if r.Id <= 1 {
		return jsonReturn(ctx, 0, "非法参数")
	}
	err = admin_model.DefaultRole.DelRole(r.Id)
	if err != nil {
		return jsonReturn(ctx, 0, err.Error())
	}
	return jsonReturn(ctx, 200, "删除成功")
}

type RolePowerModifyHandler struct {
	Id   int    `param:"<in:query><required><name:id>"`
	Cids string `parmam:"<in:formData><name:cids>"`
}

func (r *RolePowerModifyHandler) Serve(ctx *faygo.Context) error {
	id_str := ctx.QueryParam("id")
	if id_str == "" {
		return jsonReturn(ctx, 0, "参数不能为空")
	}
	id, err := strconv.Atoi(id_str)
	if err != nil {
		return jsonReturn(ctx, 0, "参数解析出错")
	}
	if id <= 1 {
		return jsonReturn(ctx, 0, "非法参数")
	}
	if ctx.Method() == "GET" {
		//获取当前的角色权限
		role, err := admin_model.DefaultRole.GetPowerByRid(id)
		if err != nil {
			return jsonReturn(ctx, 0, err.Error())
		}
		//获取所有的权限数据
		powerModel := admin_model.DefaultPower
		columns, err := powerModel.GetAllColumns(role.Power)
		if err != nil {
			return jsonReturn(ctx, 0, err.Error())
		}
		columns_str := r.TreeColumns(columns, true)
		render_data := make(map[string]interface{})
		render_data["columns_str"] = columns_str
		render_data["rid"] = id
		return ctx.Render(200, "admin_view/role_power_modify.html", render_data)
	}
	err = ctx.BindForm(r)
	if err != nil {
		return jsonReturn(ctx, 0, "参数解析出错")
	}
	cids := strings.Split(r.Cids, ",")

	//获取所有的栏目id
	cids_map, _, err := admin_model.DefaultPower.GetAllCids()
	for _, v := range cids {
		cid, err := strconv.Atoi(v)
		if err != nil {
			return jsonReturn(ctx, 0, "参数解析出错")
		}
		//如果不存在该栏目id
		if _, ok := cids_map[cid]; !ok {
			return jsonReturn(ctx, 0, "参数有误")
		}
	}
	//修改角色的权限id
	err = admin_model.DefaultRole.ModifyPower(id, r.Cids)
	if err != nil {
		return jsonReturn(ctx, 0, err.Error())
	}
	return jsonReturn(ctx, 200, "修改成功")
}

//以树桩返回栏目列表的html
func (r *RolePowerModifyHandler) TreeColumns(columns []admin_model.Column, is_top bool) (str string) {

	if is_top {
		str = `<ul class="layui-box layui-tree" style="margin-left: 20px;">`
	} else {
		str = `<ul>`
	}
	for _, v := range columns {
		str += `<li>`
		if len(v.SonPower) > 0 {
			str += `<i class="layui-icon layui-tree-spread"></i>`
		}
		str += ` <div class="checkbox-group" >`
		if v.HasColumn {
			str += "<input type=\"checkbox\" id=\"" + fmt.Sprintf("column%v", v.Id) + "\" value=" + fmt.Sprintf("%v", v.Id) + " checked/>"
		} else {
			str += "<input type=\"checkbox\" id=\"" + fmt.Sprintf("column%v", v.Id) + "\" value=" + fmt.Sprintf("%v", v.Id) + " />"
		}

		str += "<label for=\"" + fmt.Sprintf("column%v", v.Id) + "\">" + v.Name + "</label>"
		str += `</div>`
		if len(v.SonPower) > 0 {
			str += r.TreeColumns(v.SonPower, false)
		}
		str += `</li>`

	}
	str += `</ul>`
	return
}
