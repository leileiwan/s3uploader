package options

import "github.com/spf13/pflag"

type Config struct {
	Endpoint string
	AccessKeyID string
	SecretAccessKey string
	Bucket string
	ContentType string
	Region string
	UseSSL bool
	ProcessName string
	TargetPath string
}

func (cfg *Config)AddFlags(flags *pflag.FlagSet){
	flags.StringVar(&cfg.Endpoint,"endpoint","play.min.io","The minio address")
	flags.StringVar(&cfg.AccessKeyID,"access-key-id","Q3AM3UQ867SPQQA43P2F","the access-key-id for minio")
	flags.StringVar(&cfg.SecretAccessKey,"secret-access-key","zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG","the access-key-id for minio")
	flags.StringVar(&cfg.Bucket,"bucket-name","my-bucket","the bucket name")
	flags.StringVar(&cfg.ContentType,"content-type","application/zip","the s3uploader file type")
	flags.StringVar(&cfg.Region,"region","us-east-1","the minio region")
	flags.BoolVar(&cfg.UseSSL,"use-ssl",true,"use ssl")
	flags.StringVar(&cfg.ProcessName,"process-name","mydumper","dumper process name")
	flags.StringVar(&cfg.TargetPath,"target-path","/data","dump mysql data target path")
}
