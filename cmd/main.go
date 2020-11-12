package main

import (
	"fmt"
	"k8s.io/klog"
	"os"
	"s3uploader/cmd/app"
)

func main() {
	klog.InitFlags(nil)
	defer klog.Flush()

	cmd := app.NewUploaderCommand()
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v/n", err)
		os.Exit(1)
	}
}
