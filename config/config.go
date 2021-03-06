package conf

/*参数说明
app.port // 应用端口
app.upload_file_path // 图片上传的临时文件夹目录，绝对路径！
app.cookie_key // 生成加密session
app.serve_type // Serve:正常websocket  or GoServe: goroutine + channel websocket
mysql.dsn // mysql 连接地址dsn
*/

var AppJsonConfig = []byte(`
{
  "app": {
    "port": ":1234",
    "upload_file_path": "e:\\golang\\www\\go-gin-chat\\tmp_images\\",
    "cookie_key": "4238uihfieh49r3453kjdfg",
	"token_key": "gin-chat-server",
    "serve_type": "GoServe"
  },
  "mysql": {
    "dsn": "root:1998.cjc@tcp(localhost:3306)/gin-chat?charset=utf8mb4&parseTime=True&loc=Local"
  }
}
`)
