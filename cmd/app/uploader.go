package app

import (
	"context"
	"fmt"
	"klog"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"os"
	"s3uploader/cmd/app/options"
	"s3uploader/pkg/worker"
)

func NewUploaderCommand()*cobra.Command{
	cfg:=&options.Config{}
	cmd:=&cobra.Command{
		Use:"mydumper uploader",
		Long: "The s3uploader tool s3uploader mydumper data to minio",
		Run: func(cmd *cobra.Command, args []string) {
			if err:=run(cmd,args,cfg);err!=nil{
				fmt.Fprintf(os.Stderr,"%v\v",err)
				os.Exit(1)
			}
		},
	}
	cfg.AddFlags(cmd.Flags())
	return cmd
}

func run(cmd *cobra.Command, args []string,config *options.Config)error{
	// logs the flags in the flagset
	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		klog.V(1).Infof("FLAG: --%s=%q", flag.Name, flag.Value)
	})

	worker,err:=worker.NewWorker(config.Endpoint,config.AccessKeyID,config.SecretAccessKey,config.Bucket,config.ContentType,config.Region,config.UseSSL,config.ProcessName,config.TargetPath)
	if err!=nil{
		klog.Errorf("new worker with error %v",err)
		return err
	}
	ctx,cancel:=context.WithCancel(context.Background())
	worker.Run(ctx,cancel)
	return nil
}
