# s3uploader work process
1. mydumper dump mysql data to local disk, set appropriate file size
2. when mydumper dumping the data, s3uploader upload the data to minio,then rm file in local disk
3. when all files  uploaded, s3uploader process exit.

# we should solve some problems

##  dump data speed is faster than upload data speed
we found mydumper will io hang when local disk is not enough from the mydumper issue. so when disk is not enough, we have enough time to upload data to minio.

## not upload dumping file to minio
In fact, the dumping file is the file opened by mydumper process. so the file can be uploaded is file in target folder but not dumping and uploaded file.

we use go channel achieve run get file can uploaded and upload file one by one

## judge upload worker finish
when mydumper process exited and all files uploaded, s3uploader process should exit.  

