package util

import (
	"klog"
	"github.com/shirou/gopsutil/process"
)




// findProcess find process with name
func FindProcess(name string)(*process.Process,bool,error){
	ps,err:=process.Processes()
	if err!=nil{
		klog.Errorf("Get processes error :%v",err)
		return nil,false,err
	}
	for i:=range ps{
		fn,err:= ps[i].Name()
		if err!=nil{
			klog.Warningf("Get process %v name with error:%v",ps[i].String(),err)
			continue
		}
		if fn==name{
			return ps[i],true,nil
		}
	}
	klog.Warningf("Not found the process %s",name)
	return nil,false,nil
}

func FindProcessOpenfiles(p *process.Process)([]string,error){
	fs,err:=p.OpenFiles()
	if err!=nil{
		klog.Errorf("process %v open files error: %v",*p,err)
	}

	res:=[]string{}
	for i:=range fs{
		path:=fs[i].Path
		res=append(res,path)
	}
	return res,nil
}

