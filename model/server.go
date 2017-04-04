package model

// Server 表
type Server struct {
	ID       uint   `gorm:"primary_key;AUTO_INCREMENT"`
	Host     string `grom:"size:100"`
	User     string `grom:"size:100"`
	Password string `grom:"size:255"`
	Port     string `grom:"default:'22';size:5"`
}
