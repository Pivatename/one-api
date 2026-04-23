package common

import (
    "flag"
    "fmt"
    "github.com/songquanpeng/one-api/common/config"
    "github.com/songquanpeng/one-api/common/logger"
    "log"
    "os"
    "path/filepath"
    "strconv" // 必须有这个，用于解析端口
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

    // --- Vercel 适配：优先读取动态分配的端口 ---
    if os.Getenv("PORT") != "" {
        p, _ := strconv.Atoi(os.Getenv("PORT"))
        *Port = p
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

    // --- Vercel 适配：强制将日志路径指向 /tmp ---
    if os.Getenv("VERCEL") != "" {
        *LogDir = "/tmp"
    }

    if *LogDir != "" {
        var err error
        *LogDir, err = filepath.Abs(*LogDir)
        if err != nil {
            log.Fatal(err)
        }
        if _, err := os.Stat(*LogDir); os.IsNotExist(err) {
            err = os.MkdirAll(*LogDir, 0777)
            if err != nil {
                // 在 Vercel 上如果创建失败也不要让程序崩溃
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
