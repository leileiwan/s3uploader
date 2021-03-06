package worker

import (
	"context"
	"k8s.io/klog"
	"time"
)

type Worker struct {
	uploader        Uploader
	dumper          Dumper
	uploadableFiles []string
	signal1         chan bool
	signal2         chan bool
}

func NewWorker(endpoint, accessKeyID, secretAccessKey, bucket, contentType, region string, useSSL bool, processName, targetPath string) (*Worker, error) {
	uploader, err := NewMinioUploader(endpoint, accessKeyID, secretAccessKey, bucket, contentType, region, useSSL)
	if err != nil {
		klog.Errorf("New Uploader with error:%v", err)
		return nil, err
	}

	return &Worker{
		uploader:        uploader,
		dumper:          NewMyDumper(processName, targetPath),
		uploadableFiles: []string{},
		signal1:         make(chan bool, 1),
		signal2:         make(chan bool, 1),
	}, nil
}

func (w *Worker) Run(ctx context.Context, cancel context.CancelFunc) {
	w.signal1 <- true
	go func() {
		for {
			<-w.signal1
			_, err := w.findUploadableFiles()
			if err != nil {
				cancel()
				return
			}
			w.signal2 <- true
		}
	}()

	go func(ctx context.Context) {
		for {
			<-w.signal2
			err := w.upload(ctx)
			if err != nil {
				cancel()
				return
			}
			w.signal1 <- true
		}
	}(ctx)

	go func() {
		for {
			if w.isFinished() {
				klog.Infoln("upload worker finished")
				cancel()
			}
			time.Sleep(1 * time.Second)
		}
	}()
	<-ctx.Done()
}

func (w *Worker) findUploadableFiles() ([]string, error) {
	klog.Infoln("Get uploadable files...")

	dumpedFiles, err := w.dumper.FindDumpedFiles()
	if err != nil {
		klog.Errorf("Worker get dumpedFiles with error:%v", err)
		return nil, err
	}

	res := []string{}
	uploadedFiles := w.uploader.GetUploadedFiles()
	for _, path := range dumpedFiles {
		if !uploadedFiles.Has(path) {
			res = append(res, path)
		}
	}

	w.uploadableFiles = res
	return res, nil
}

func (w *Worker) upload(ctx context.Context) error {
	klog.Infoln("Upload file...")
	for _, path := range w.uploadableFiles {
		err := w.uploader.Upload(path, ctx)
		if err != nil {
			klog.Errorf("Upload file %s with error:%v", path, err)
			return err
		}
	}
	w.uploadableFiles = []string{}
	return nil
}

// isFinished judge whether uploading worker should finish
func (w *Worker) isFinished() bool {
	isExisted, err := w.dumper.IsExist()
	if err != nil {
		klog.Errorf("Run dumper IsExist with error:%v", err)
		return false
	}
	if isExisted {
		return false
	}

	allFiles, err := w.dumper.FindAllFiles()
	if len(w.uploadableFiles) == 0 && len(allFiles) == 0 {
		return true
	}
	return false
}
