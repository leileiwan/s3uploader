# What's S3Uploader
S3Uploader is a data upload tool to upload dumped(mysql/mydumper) data to s3 object storage.

It still can upload data when local disk is not enough.

# How to build it
Build:
```
make
```

Clean:
```
make clean
```

# How to use it
Dependencies:
 
install mysql and minio, then dumper data to many some files with -s -F
```
apt-get install mysql
docker run -p 9000:9000   -e "MINIO_ACCESS_KEY=AKIAIOSFODNN7EXAMPLE"   -e "MINIO_SECRET_KEY=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"   minio/minio server /data
sudo mydumper -u root -p 123 -B mytest -s 100 -F 1 -t 1  -o /home/wanlei/PingCap/data
```

s3upload useage:
```
Usage:
  mydumper uploader [flags]

Flags:
      --access-key-id string       the access-key-id for minio (default "Q3AM3UQ867SPQQA43P2F")
      --bucket-name string         the bucket name (default "my-bucket")
      --content-type string        the s3uploader file type (default "application/zip")
      --endpoint string            The minio address (default "play.min.io")
  -h, --help                       help for mydumper
      --process-name string        dumper process name (default "mydumper")
      --region string              the minio region (default "us-east-1")
      --secret-access-key string   the access-key-id for minio (default "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG")
      --target-path string         dump mysql data target path (default "/data")
      --use-ssl                    use ssl (default true)
```


# how it works
[deign](./doc/design.md)
