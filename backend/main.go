package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
)

const DATA_PATH = "/lzcapp/var/"
const MNT_PATH = "/lzcapp/var/"

type DiskInfo struct {
	DiskTotal     string `json:"disk_total"`
	DiskUsed      string `json:"disk_used"`
	DiskAvailable string `json:"disk_available"`
	DataUsed      string `json:"data_used"`
}

func getDiskInfo() (DiskInfo, error) {
	diskInfo := DiskInfo{}

	// 获取 /data 的磁盘信息
	dfOutput, err := exec.Command("df", MNT_PATH).Output()
	if err != nil {
		return diskInfo, err
	}

	lines := strings.Split(string(dfOutput), "\n")
	if len(lines) < 2 {
		return diskInfo, fmt.Errorf("failed to parse df output")
	}

	// 解析 df 输出
	fields := strings.Fields(lines[1])
	if len(fields) < 4 {
		return diskInfo, fmt.Errorf("unexpected df output format")
	}
	diskInfo.DiskTotal = fields[1]
	diskInfo.DiskUsed = fields[2]
	diskInfo.DiskAvailable = fields[3]

	// 获取 /data/code 的已使用大小
	duOutput, err := exec.Command("du", "-s", DATA_PATH).Output()
	if err != nil {
		return diskInfo, err
	}

	duFields := strings.Fields(string(duOutput))
	if len(duFields) > 0 {
		diskInfo.DataUsed = duFields[0]
	} else {
		return diskInfo, fmt.Errorf("failed to parse du output")
	}

	return diskInfo, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	diskInfo, err := getDiskInfo()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(diskInfo)
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Listening on port 8081...")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		fmt.Println("Failed to start server:", err)
	}
}
