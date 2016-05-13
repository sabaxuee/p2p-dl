# p2p-dl

## 测试方法

1. go build p2p-dl.go
2. ./create_testfile.sh 
    创建测试文件。
3. ./p2p-dl & 
4. 访问 http://127.0.0.1:9090/pull?f=test.dat&src=127.0.0.1:9090
     
    或使用curl: 
    curl --request GET
    'http://127.0.0.1:9090/pull?f=test.dat&src=127.0.0.1:9090'


    如果在多台机器测试，把上述url相应IP换成目标IP即可。

    最后会在程序所在目录看到下载的文件。
    用md5sum比较两个文件。

#issue
    单线程下载。
    下载过程中重复提交url会导致下载的文件出错，程序目前没有对重复提交作处理。


#Docker
    docker build -t="ubuntu:localbuild" .
    docker run -t -i ubuntu:localbuild bash
    #cd /home/p2p-dl/
    #./create_testfile.sh 
    #./p2p-dl &
