package main

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	jarName         = "your-service.jar"           // 要监控的jar包名称
	jarStartCommand = "java -jar your-service.jar" // 启动jar包的命令
	checkInterval   = 10 * time.Second             // 检查的时间间隔
	logFilePath     = "monitor.log"                // 日志文件路径
)

func main() {
	// 打开日志文件
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("无法打开日志文件: %v", err)
	}
	defer logFile.Close()

	// 将日志输出到文件
	logger := log.New(logFile, "", log.LstdFlags)

	logger.Println("启动Jar包监控服务...")

	for {
		// 检查jar包是否正在运行
		isRunning := checkJarRunning(jarName)

		if isRunning {
			logger.Printf("服务 %s 正在运行。\n", jarName)
		} else {
			logger.Printf("服务 %s 未运行，正在尝试重启...\n", jarName)
			err := restartJar(jarStartCommand, logger)
			if err != nil {
				logger.Printf("无法启动服务 %s: %v\n", jarName, err)
			} else {
				logger.Printf("服务 %s 已成功启动。\n", jarName)
			}
		}

		// 等待下一个检查周期
		time.Sleep(checkInterval)
	}
}

// 检查指定的jar包是否正在运行
func checkJarRunning(jarName string) bool {
	cmd := exec.Command("ps", "aux")
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		log.Printf("检查进程时出错: %v", err)
		return false
	}

	scanner := bufio.NewScanner(&out)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, jarName) {
			return true
		}
	}

	return false
}

// 重启jar包服务
func restartJar(command string, logger *log.Logger) error {
	parts := strings.Fields(command)
	cmd := exec.Command(parts[0], parts[1:]...)

	// 设置日志输出到文件
	cmd.Stdout = logger.Writer()
	cmd.Stderr = logger.Writer()

	err := cmd.Start()
	if err != nil {
		return err
	}

	return nil
}
