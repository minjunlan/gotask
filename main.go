package main

import (
	"strconv"

	"time"

	"os/exec"

	"strings"

	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/hypersleep/easyssh"
	"github.com/jinzhu/gorm"
	"github.com/minjunlan/gotask/common"
	"github.com/minjunlan/gotask/log"
	"github.com/minjunlan/gotask/model"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
)

const DATE_FORMAT_TYPE = "2006-01-02 15:04:05"

type Wherel map[string]string

type CallFunc func(args ...string)

var whe = Wherel{
	"day":      "每天",
	"day-n":    "N天",
	"hour":     "每小时",
	"hour-n":   "N小时",
	"minute-n": "N分钟",
	"week":     "每星期",
	"month":    "每月",
}

var week = map[string]string{
	"0": "星期天",
	"1": "星期一",
	"2": "星期二",
	"3": "星期三",
	"4": "星期四",
	"5": "星期五",
	"6": "星期六",
}

var num = 0

func main() {
	r := gin.Default()

	//设置计划任务
	cron := cron.New()
	defer cron.Stop()
	cron.Start()

	//启动数据库
	db := model.Init()
	defer db.Close()

	// 加载原来的计划任务
	go loadCron(db, cron)

	r.Static("/public", "public")
	r.LoadHTMLGlob("tpl/*")

	r.GET("/", func(c *gin.Context) {
		var task, rs []model.Task
		db.Find(&task)
		for _, v := range task {
			v.Ntype = whe[v.Ntype]
			rs = append(rs, v)
		}

		c.HTML(200, "index.html", gin.H{"title": "计划任务", "data": rs})
	})

	r.GET("/servers", func(c *gin.Context) {
		var servs []model.Server
		db.Find(&servs)
		c.HTML(200, "serverlist.html", gin.H{"title": "服务器列表", "Servs": servs})
	})

	r.GET("/addServer", func(c *gin.Context) {
		var servs []model.Server
		db.Find(&servs)
		c.HTML(200, "addServer.html", gin.H{"title": "服务器列表"})
	})
	r.POST("/addServer", func(c *gin.Context) {
		s := model.Server{Host: c.PostForm("Host"), User: c.PostForm("User"), Password: c.PostForm("Password"), Port: c.PostForm("Port")}
		if db.NewRecord(s) {
			db.Create(&s)
			c.JSON(200, gin.H{"status": true})
		} else {
			c.JSON(200, gin.H{"status": false})
		}

	})

	r.GET("/editServer/:id", func(c *gin.Context) {
		id := c.Param("id")
		var servs model.Server
		db.First(&servs, id)
		c.HTML(200, "addServer.html", servs)
	})

	r.POST("/editServer/:id", func(c *gin.Context) {
		id := c.Param("id")
		var servs model.Server
		db.First(&servs, id)
		servs.Host = c.PostForm("Host")
		servs.User = c.PostForm("User")
		servs.Password = c.PostForm("Password")
		servs.Port = c.PostForm("Port")
		db.Save(servs)
		c.JSON(200, gin.H{"status": true})
	})

	r.POST("/delServer/:id", func(c *gin.Context) {
		id := c.Param("id")
		var servs model.Server
		db.First(&servs, id)
		db.Delete(&servs)
		go delCron(db, cron)
		c.JSON(200, gin.H{"status": true})
	})

	r.POST("/addCron", func(c *gin.Context) {
		fmt.Printf("%v \r\n", c)
		//得到form值
		f := make(map[string]string)
		f["name"] = c.PostForm("name")
		f["type"] = c.PostForm("type")
		f["where1"] = c.PostForm("where1")
		f["week"] = c.PostForm("week")
		f["hour"] = c.PostForm("hour")
		f["minute"] = c.PostForm("minute")
		f["sType"] = c.PostForm("sType")
		f["sBody"] = c.PostForm("sBody")
		f["serverid"] = c.PostForm("serverid")
		fmt.Printf("from:%#v \r\n", f)

		//添加 定时
		opportunity, corn := RunCron(cron, &f)
		fmt.Printf("%#v \r\n", opportunity)
		id, _ := strconv.Atoi(f["serverid"])
		fmt.Printf("%#v \r\n", id)
		//插入数据库
		task := model.Task{Name: f["name"], Ntype: f["type"], Opportunity: opportunity, Time: time.Now().Unix(), Cron: corn, Stype: f["sType"], Sbody: f["sBody"], ServerId: uint(id)}
		fmt.Printf("task:%#v \r\n", task)
		if db.NewRecord(task) {
			db.Create(&task)
			//返回值
			c.JSON(200, gin.H{"Name": f["name"], "Ntype": whe[f["type"]], "Opportunity": opportunity, "Time": time.Now().Format(DATE_FORMAT_TYPE)})
			//return
		} else {
			c.JSON(500, "插入数据出错！")
		}

	})
	r.POST("/delCron/:id", func(c *gin.Context) {
		id := c.Param("id")
		var task model.Task
		db.First(&task, id)
		db.Delete(&task)
		c.JSON(200, gin.H{"status": true})
	})
	r.GET("/getServers", func(c *gin.Context) {
		var servs []model.Server
		db.Find(&servs)
		c.JSON(200, gin.H{"status": true, "data": servs})
	})
	r.Run(":24115")
	//endless.ListenAndServe(":24115", r)

}

// RunCron ...
func RunCron(cron *cron.Cron, f *map[string]string) (string, string) {

	str := ""
	msg := ""
	switch form := *f; form["type"] {
	case "day":
		str = "0 " + form["minute"] + " " + form["hour"] + " * * *"
		msg = "每天的，" + form["hour"] + "点" + form["minute"] + "分钟运行"

	case "day-n":
		str = "0 " + form["minute"] + " " + form["hour"] + " */" + form["where1"] + " * *"
		msg = "每隔" + form["where1"] + "天，" + form["hour"] + "点" + form["minute"] + "分钟运行"

	case "hour":
		str = "0 " + form["minute"] + " * * * *"
		msg = "每小时的，" + form["minute"] + "分钟运行"

	case "hour-n":
		str = "0 " + form["minute"] + " */" + form["where1"] + " * * *"
		msg = "每隔" + form["where1"] + "小时，" + form["minute"] + "分钟运行"

	case "minute-n":
		str = "0 " + " */" + form["where1"] + " * * * *"
		msg = "每隔" + form["where1"] + "分钟，" + "运行"

	case "week":
		str = "0 " + form["minute"] + " " + form["hour"] + " * * " + form["week"]
		msg = "每星期的," + week[form["week"]] + form["hour"] + "点" + form["minute"] + "分钟运行"

	case "month":
		str = "0 " + form["minute"] + " " + form["hour"] + " " + form["where1"] + " * *"
		msg = "每月的，" + form["where1"] + "日" + form["hour"] + "点" + form["minute"] + "分钟运行"

	}

	cron.AddFunc(str, func() {
		form := *f
		switch form["sType"] {
		case "toLocal":
			go runLocalCmd(form["sBody"])
		case "toShell":
			var h model.Server
			id, _ := strconv.Atoi(form["serverid"])
			model.DB.First(&h, uint(id))
			go runRemoteSSH(form["sBody"], &h)
		case "toFunc":
			id, _ := strconv.Atoi(form["serverid"])
			ss := strings.Split(form["sBody"], "|")
			sss := strings.Split(ss[2], ",")
			go common.Callback(new(common.Inerfunc), ss[0], id, ss[1], sss)
			//go download()
		}

	})

	return msg, str
}

// runRemoteSSH 运行远程ssh
func runRemoteSSH(str string, h *model.Server) {
	ssh := &easyssh.MakeConfig{}
	ssh.User = h.User
	ssh.Server = h.Host
	ssh.Password = h.Password
	ssh.Port = h.Port
	log.Log.WithFields(logrus.Fields{
		"Host": ssh.Server,
	}).Info("运行远程命令")
	_, _, done, err := ssh.Run(str, 3600)
	// Handle errors
	if err != nil && !done {
		panic("Can't run remote command: " + err.Error())
		log.Log.WithFields(logrus.Fields{
			"Host":  ssh.Server,
			"error": err.Error(),
		}).Fatal("运行远程命令失败！")
	}
}

// runLocalCmd 运行本地命令
func runLocalCmd(cmd string) {
	exec.Command(cmd)
	log.Log.WithFields(logrus.Fields{
		"cmd": cmd,
	}).Info("运行本地命令")
}

func loadCron(db *gorm.DB, cron *cron.Cron) {
	var task []model.Task
	db.Find(&task)
	for _, v := range task {
		cron.AddFunc(v.Cron, func() {
			switch v.Stype {
			case "toLocal":
				go runLocalCmd(v.Sbody)
			case "toShell":
				var h model.Server
				model.DB.First(&h, v.ServerId)
				go runRemoteSSH(v.Sbody, &h)
			case "toFunc":
				ss := strings.Split(v.Sbody, "|")
				sss := strings.Split(ss[2], ",")
				go common.Callback(new(common.Inerfunc), ss[0], int(v.ServerId), ss[1], sss)
				//go download()
			}

		})
	}

}

func delCron(db *gorm.DB, cron *cron.Cron) {
	cron.Stop()

	defer func() {
		cron.Start()
		loadCron(db, cron)
	}()

}
