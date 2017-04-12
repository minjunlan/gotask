package common

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"time"

	"strings"

	"path"

	"github.com/minjunlan/gotask/log"
	"github.com/minjunlan/gotask/model"
	"github.com/pkg/sftp"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

// Inerfunc 内部函数
type Inerfunc struct {
}

// Callback 通过它调用函数
func Callback(object interface{}, methodName string, args ...interface{}) {
	inputs := make([]reflect.Value, len(args))
	for i := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}
	reflect.ValueOf(object).MethodByName(methodName).Call(inputs)
}

// Download 用于远程下载
func (Inerfunc) Download(servid int, savepath string, files []string) {
	var h model.Server
	model.DB.First(&h, uint(servid))
	config := &ssh.ClientConfig{
		User: h.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(h.Password),
		},
	}
	conn, err := ssh.Dial("tcp", h.Host+":"+h.Port, config)
	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"Host":  h.Host,
			"User":  h.User,
			"Pwd":   h.Password,
			"error": err,
		}).Fatal("运行远程命令失败！")
	}

	sftp, err := sftp.NewClient(conn)
	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"Host":  h.Host,
			"User":  h.User,
			"Pwd":   h.Password,
			"error": err,
		}).Fatal("服务器连接错误")
	}
	defer sftp.Close()

	ti := time.Now().Format("20060102")
	sp := strings.Replace(savepath, "{date}", ti, 1)
	os.MkdirAll(sp, os.ModePerm)

	fmt.Println(sp)

	ch := make(chan struct{})
	for _, f := range files {
		log.Log.WithFields(logrus.Fields{
			"file": f,
		}).Info("下载文件")
		go download(sftp, sp, f, ch)
	}

	//等待下载完成
	for range files {
		<-ch
	}
}

func download(sftp *sftp.Client, savepath string, filepath string, ch chan struct{}) {

	file, err := sftp.Open(filepath)
	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"file":  file,
			"error": err,
		}).Fatal("打开文件错误")
		return
	}

	info, err := file.Stat()
	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"info":  info,
			"error": err,
		}).Fatal("获取文件信息错误")
		return
	}

	buf := make([]byte, 1024*1024) // 读 1M 数据

	savefile, err := os.OpenFile(path.Join(savepath, info.Name()), os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"savefile": savefile,
			"error":    err,
		}).Info("创建文件错误！")
	}
	for {
		n, err1 := file.Read(buf)
		if n == 0 || err1 == io.EOF {
			savefile.Close()
			break
		}
		n1, err := savefile.Write(buf)
		//err = ioutil.WriteFile(path.Join(savepath, info.Name()), buf, os.ModePerm)
		if err != nil {
			log.Log.WithFields(logrus.Fields{
				"n1":     n1,
				"error":  err,
				"error1": err1,
				"n":      n,
			}).Info("写文件错误！")
			return
		}
	}
	buf = nil
	ch <- struct{}{}
}
