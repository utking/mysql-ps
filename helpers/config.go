package helpers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"gopkg.in/ini.v1"
)

var (
	myCnfPath = "~/.my.cnf"
)

func ExpandMyCnfPath() string {
	if strings.HasPrefix(myCnfPath, "~/") {
		dirname, _ := os.UserHomeDir()
		return filepath.Join(dirname, myCnfPath[2:])
	}
	return myCnfPath
}

func LoadConfig() {
	path := ExpandMyCnfPath()

	cfg, err := ini.Load(path)

	if err == nil {
		host := cfg.Section("client").Key("host").String()
		socket := cfg.Section("client").Key("socket").String()
		user := cfg.Section("client").Key("user").String()
		password := cfg.Section("client").Key("password").String()

		if host != "" && os.Getenv("MYSQL_DSN") == "" {
			os.Setenv("MYSQL_DSN", fmt.Sprintf("tcp(%s)", host))
		}

		if socket != "" {
			os.Setenv("MYSQL_DSN", fmt.Sprintf("unix(%s)", socket))
		}

		if user != "" && os.Getenv("MYSQL_USER") == "" {
			os.Setenv("MYSQL_USER", user)
		}

		if user != "" && os.Getenv("MYSQL_PASSWORD") == "" {
			os.Setenv("MYSQL_PASSWORD", password)
		}
	}

	_ = godotenv.Load()
}
