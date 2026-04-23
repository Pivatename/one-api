package common

import (
	"flag"
	"fmt"
	"github.com/songquanpeng/one-api/common/config"
	"github.com/songquanpeng/one-api/common/logger"
	"log"
	"os"
	"path/filepath"
)

var (
	Port         = flag.Int("port", 3000, "the listening port")
	PrintVersion = flag.Bool("version", false, "print version and exit")
	PrintHelp    = flag.Bool("help", false, "print help and exit")
	LogDir       = flag.String("log-dir", "./logs", "specify the log directory")
)

func printHelp() {
	fmt.Println("One API " + Version + " - All in one API service for OpenAI API.")
	fmt.Println("Copyright (C) 2023 JustSong. All rights reserved.")
	fmt.Println("GitHub: https://github.com/songquanpeng/one-api")
	fmt.Println("Usage: one-api [--port <port>] [--log-dir <log directory>] [--version] [--help]")
}

func Init() {
	flag.Parse()
	// 增加这一段：让程序支持 Vercel 动态端口
      if os.Getenv("PORT") != "" {
            port, _ := strconv.Atoi(os.Getenv("PORT"))
            *Port = port
        }

	if *PrintVersion {
		fmt.Println(Version)
		os.Exit(0)
	}

	if *PrintHelp {
		printHelp()
		os.Exit(0)
	}

	if os.Getenv("SESSION_SECRET") != "" {
		if os.Getenv("SESSION_SECRET") == "random_string" {
			logger.SysError("SESSION_SECRET is set to an example value, please change it to a random string.")
		} else {
			config.SessionSecret = os.Getenv("SESSION_SECRET")
		}
	}
	if os.Getenv("SQLITE_PATH") != "" {
		SQLitePath = os.Getenv("SQLITE_PATH")
	}
	if *LogDir != "" {
             // Vercel 环境兼容：强制将日志目录指向 /tmp
             if os.Getenv("VERCEL") != "" {
                 *LogDir = "/tmp"
             }
             
             var err error
             *LogDir, err = filepath.Abs(*LogDir)
             if err != nil {
                log.Fatal(err)
            }
            if _, err := os.Stat(*LogDir); os.IsNotExist(err) {
                err = os.Mkdir(*LogDir, 0777)
                if err != nil {
                    // 在 Vercel 上即使创建失败也不要让程序崩溃
                    if os.Getenv("VERCEL") != "" {
                        fmt.Printf("Warning: failed to create log dir: %v\n", err)
                    } else {
                        log.Fatal(err)
                    }
                }
            }
            logger.LogDir = *LogDir
       }

}
