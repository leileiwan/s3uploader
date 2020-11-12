package worker

import (
	"io/ioutil"
	"k8s.io/apimachinery/pkg/util/sets"
	"sync"
	"s3uploader/pkg/util"
	"klog"
	"fmt"
)
type Dumper interface {
	 IsExist() (bool,error)
	 FindDumpingFiles() (sets.String,error)
	 FindAllFiles() ([]string,error)
	 FindDumpedFiles() ([]string,error)
}

type MyDumper struct {
	processName   string
	targetPath    string
	allFiles []string
	dumpingFiles sets.String
	dumpedFiles []string
}

var defaultMyDumper *MyDumper
var defaultMyDumperLock sync.Mutex

// use single-mod new MyDumper
func NewMyDumper(processName,targetPath string)*MyDumper{
	defaultMyDumperLock.Lock()
	defer defaultMyDumperLock.Unlock()
	if defaultMyDumper==nil{
		defaultMyDumper=&MyDumper{
			processName: processName,
			targetPath: targetPath,
			allFiles: []string{},
			dumpingFiles: sets.String{},
			dumpedFiles: []string{},
		}
	}
	return defaultMyDumper
}

// IsExist judge dumper process is exist
func (d *MyDumper)IsExist()(bool,error){
	_,found,err:=util.FindProcess(d.processName)
	return found,err
}

// FindAllFiles get all dump file
func (d *MyDumper)FindAllFiles()([]string,error){
	files,err:=ioutil.ReadDir(d.targetPath)
	if err!=nil{
		klog.Errorf("Get path %s files with error:%v",d.targetPath,err)
		return nil,err
	}
	fs:=[]string{}
	for i:=range files{
		if !files[i].IsDir(){
			eachPath:=fmt.Sprintf("%s/%s",d.targetPath,files[i].Name())
			fs=append(fs,eachPath)
		}
	}
	d.allFiles=fs
	return fs,nil
}

// FindDumpingFiles get all dumpling file
func (d *MyDumper)FindDumpingFiles() (sets.String, error) {
	process,found,err:=util.FindProcess(d.processName)
	if err!=nil{
		return nil, err
	}
	if !found{
		return nil, nil
	}

	res:=sets.String{}
	files,err:=util.FindProcessOpenfiles(process)
	if err!=nil{
		return nil, err
	}
	res.Insert(files...)
	return res,nil
}

// FindDumpedFiles get all dumped files
func (d *MyDumper)FindDumpedFiles()([]string,error){
	_,err:=d.FindAllFiles()
	if err!=nil{
		return nil,err
	}

	_,err=d.FindDumpingFiles()
	if err!=nil{
		return nil,err
	}

	res:=[]string{}
	for _,path:=range d.allFiles{
		if d.dumpingFiles.Has(path){
			continue
		}
		res=append(res,path)
	}
	d.dumpedFiles=res
	return res,err
}
