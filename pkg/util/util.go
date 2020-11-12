package util

import (
	"bytes"
	"fmt"
	"github.com/shirou/gopsutil/process"
	"k8s.io/klog"
	"os"
	"os/exec"
	"strings"
)

// findProcess find process with name
func FindProcess(name string) (*process.Process, bool, error) {
	ps, err := process.Processes()
	if err != nil {
		klog.Errorf("Get processes error :%v", err)
		return nil, false, err
	}
	for i := range ps {
		fn, err := ps[i].Name()
		if err != nil {
			klog.Warningf("Get process %v name with error:%v", ps[i].String(), err)
			continue
		}
		if fn == name {
			return ps[i], true, nil
		}
	}
	klog.Warningf("Not found the process %s", name)
	return nil, false, nil
}

func FindProcessOpenfiles(processName, path string) ([]string, error) {
	// not implement
	//fs,err:=p.OpenFiles()
	//if err!=nil{
	//	klog.Errorf("process %v open files error: %v",*p,err)
	//}
	strCmd := fmt.Sprintf("lsof -c  %s | grep %s | awk '{print $9}'", processName, path)
	cmd := exec.Command("/bin/sh", "-c", strCmd)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Errorf("run command %s with error:%v", strCmd, err)
		return nil, err
	}
	files := strings.Split(out.String(), "\n")

	res := []string{}
	for _, eachPath := range files {
		if isDir(eachPath) || len(eachPath) < len(path) {
			continue
		}
		// /Users/wanlei/go/src/s3uploader/data/.b.swp
		strs := strings.Split(eachPath, "/")
		swpName := strs[len(strs)-1]
		name := swpName[1 : len(swpName)-4]
		dirPath := eachPath[:len(eachPath)-len(swpName)-1]
		file := fmt.Sprintf("%s/%s", dirPath, name)
		res = append(res, file)
	}
	return res, nil
}

func isDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}
