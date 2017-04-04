package model

// Task 任务表
type Task struct {
	ID          uint   `gorm:"primary_key;AUTO_INCREMENT"`
	Name        string `grom:"size:32"`
	Ntype       string `grom:"size:10"`
	Opportunity string `grom:"size:255"`
	Time        int64  //添加任务时间
	Stype       string `grom:"size:20"` //执行类型
	Sbody       string `grom:"type:text"`
	Cron        string `grom:"size:255"`  //执行周期 0 * * * * *
	Flag        byte   `grom:"default:0"` //状态 1为执行，0为未执行
	Rtime       int64  //执行时间
	ServerId    uint   `grom:"default:0"`
}
