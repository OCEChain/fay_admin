package router

import (
	"fay_admin/admin_handler"
	md "fay_admin/admin_middleware"
	"github.com/henrylee2cn/faygo"
)

// Route register router in a tree style.
func Route(frame *faygo.Framework) {
	frame.Route(
		frame.NewGroup("admin/",
			frame.NewNamedAPI("edit_admin", "GET POST", "/edit_admin", &admin_handler.EditAdminHandler{}),
			frame.NewNamedAPI("edit_pass", "GET POST", "/edit_pass", &admin_handler.EditPassHandler{}),
			frame.NewNamedAPI("quit", "POST", "/quit", &admin_handler.QuitHandler{}),
			frame.NewNamedAPI("uploader", "POST", "/uploader", &admin_handler.UploaderHandler{}),
			frame.NewNamedAPI("login", "GET POST", "/login", &admin_handler.LoginHandler{}),
			frame.NewNamedAPI("index", "GET", "/index", &admin_handler.IndexHandler{}),
			frame.NewNamedAPI("user", "GET POST", "/user", &admin_handler.UserHandler{}).Use(md.Check),
			frame.NewNamedAPI("user_add", "GET POST", "/user_add", &admin_handler.UserAddHandler{}).Use(md.Check),
			frame.NewNamedAPI("user_del", "GET POST", "/user_del", &admin_handler.UserDelHandler{}).Use(md.Check),
			frame.NewNamedAPI("user_setstatus", "POST", "/user_setstatus", &admin_handler.UserSetStatus{}).Use(md.Check),
			frame.NewNamedAPI("role", "GET POST", "/role", &admin_handler.RoleHandler{}).Use(md.Check),
			frame.NewNamedAPI("role_add", "POST", "/role_add", &admin_handler.RoleAddHandler{}).Use(md.Check),
			frame.NewNamedAPI("role_edit", "POST", "/role_edit", &admin_handler.RoleEditHandler{}).Use(md.Check),
			frame.NewNamedAPI("role_del", "POST", "/role_del", &admin_handler.RoleDelHandler{}).Use(md.Check),
			frame.NewNamedAPI("role_power_modify", "GET POST", "/role_power_modify", &admin_handler.RolePowerModifyHandler{}).Use(md.Check),
		),
		frame.NewNamedStaticFS("assets", "/assets", faygo.MarkdownFS(
			"./admin_static/assets",
		)),
		frame.NewNamedStaticFS("assets", "/layui", faygo.MarkdownFS(
			"./admin_static/layui",
		)),
		frame.NewNamedStaticFS("myjs", "/myjs", faygo.MarkdownFS(
			"./admin_static/myjs",
		)),
		frame.NewNamedStaticFS("upload", "/upload", faygo.MarkdownFS(
			"./upload/",
		)),
	)

}
