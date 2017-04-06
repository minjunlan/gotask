## 编译命令
### Windows
```bash
GOOS=windows GOARCH=386 go build -o Mytask.exe
```

> 1. 目前只能手动改数据库用户、密码等，改完后再编译。以后会写入文件
> 2. linux转windows时，main.go 中的模板文件路径要修改，以后会加自动识别。
> 3. 任务删除后，得重启服务器，不然删除任务不起作用 

