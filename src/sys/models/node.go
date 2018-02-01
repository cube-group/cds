package models

import (
    "golang.org/x/crypto/ssh"
    "net"
    "bytes"
    "fmt"
    "strings"
    "errors"
    "alex/utils"
)

//创建SSH连接实例
func NewSsh(address, username, password string) *Ssh {
    return &Ssh{Address:address, Username:username, Password:password}
}

//创建Node连接实例
func NewNode(address, username, password string) *Node {
    return &Node{Ssh:*NewSsh(address, username, password)}
}

//创建DockerNode连接实例
//宿主机
func NewDockerNode(sets *FSets, address string) *DockerNode {
    node := &DockerNode{Node:*NewNode(address, sets.SshUsername, sets.SshPassword)}
    node.AppMode = sets.AppMode
    node.RegistryUsername = sets.RegistryUsername
    node.RegistryPassword = sets.RegistryPassword
    node.RegistryServiceDomain = sets.RegistryServiceDomain
    node.RegistryServiceHost = sets.RegistryServiceHost
    return node
}

//创建DockerNodeBuild连接实例
//编译机
func NewDockerNodeBuild(sets *FSets) *DockerNodeBuild {
    node := &DockerNodeBuild{DockerNode:*NewDockerNode(sets, sets.BuilderAddress)}
    node.RegistryBaseDomain = sets.RegistryBaseDomain
    node.RegistryBaseHost = sets.RegistryBaseHost
    node.SvnUsername = sets.SvnUsername
    node.SvnPassword = sets.SvnPassword
    node.SvnDomain = sets.SvnDomain
    node.SvnHost = sets.SvnHost
    node.GitUsername = sets.GitUsername
    node.GitPassword = sets.GitPassword
    node.GitDomain = sets.GitDomain
    node.GitHost = sets.GitHost
    return node
}

//创建DockerNodeRegistry连接实例
//基础镜像仓库机
func NewDockerNodeBaseRegistry(sets *FSets) *DockerNodeRegistry {
    node := &DockerNodeRegistry{
        DockerNode:*NewDockerNode(sets, sets.RegistryBaseAddress),
        Registry:Registry{Username:sets.RegistryUsername, Password:sets.RegistryPassword, Domain:sets.RegistryBaseDomain},
        RegistryDockerRunName:sets.RegistryBaseRunName,
    }
    return node
}

//创建DockerNodeRegistry连接实例
//微服务镜像仓库机
func NewDockerNodeServiceRegistry(sets *FSets) *DockerNodeRegistry {
    node := &DockerNodeRegistry{
        DockerNode:*NewDockerNode(sets, sets.RegistryServiceAddress),
        Registry:Registry{Username:sets.RegistryUsername, Password:sets.RegistryPassword, Domain:sets.RegistryServiceDomain},
        RegistryDockerRunName:sets.RegistryServiceRunName,
    }
    return node
}

//创建DockerNodeRegistry连接实例
//仓库机
func NewDockerNodeTask(sets *FSets, address string) *DockerNodeTask {
    return &DockerNodeTask{DockerNode:*NewDockerNode(sets, address)}
}

//操作基类
type Operate struct {
    //ssh、stdOut、stdErr操作日志记录
    logs []string
}

//ssh类
type Ssh struct {
    Operate

    //ssh client
    client   *ssh.Client

    //ssh连接地址如:127.0.0.3:22
    Address  string
    Username string
    Password string
}

//节点机
type Node struct {
    Ssh

    //外网ip
    Ip      string
    //内网ip
    IpInner string
}

//docker宿主机类
type DockerNode struct {
    Node

    //镜像仓库统一账号
    RegistryUsername      string
    //镜像仓库统一密码
    RegistryPassword      string

    //微服务镜像仓库服务域名
    RegistryServiceDomain string
    //微服务镜像仓库服务内网hosts绑定
    RegistryServiceHost   string

    //是否为测试环境
    AppMode               string
}

//docker build编译机
type DockerNodeBuild struct {
    DockerNode

    //基础镜像仓库服务域名
    RegistryBaseDomain string
    //基础镜像仓库服务内网hosts绑定
    RegistryBaseHost   string

    //svn账号
    SvnUsername        string
    //svn密码
    SvnPassword        string
    //svn服务域名
    SvnDomain          string
    //svn服务hosts绑定
    SvnHost            string

    //git账号
    GitUsername        string
    //git密码
    GitPassword        string
    //git服务域名
    GitDomain          string
    //git服务hosts绑定
    GitHost            string
}

//docker task任务机
type DockerNodeTask struct {
    DockerNode

    //默认初始镜像
    DefaultImages []string
}

//docker registry仓库机
type DockerNodeRegistry struct {
    DockerNode
    Registry

    //docker run时的名称
    RegistryDockerRunName string
}

//记录操作日志
func (this *Operate)Log(oType interface{}, values ...string) {
    logString := fmt.Sprintf("[%v][%v]%v", utils.GetFormatYmdHis(), oType, strings.Join(values, " "))
    this.logs = append(this.logs, logString)
}

//返回所有操作日志
func (this *Operate)Logs() string {
    return strings.Join(this.logs, "\n")
}

//关闭ssh连接
func (this *Ssh)Close() error {
    if this.client != nil {
        err := this.client.Close()
        this.client = nil
        return err
    }
    this.Log("SSH-Stop")
    return nil
}

//初始化SSH连接
func (this *Ssh)connect() error {
    if (this.client != nil) {
        return nil
    }

    config := &ssh.ClientConfig{
        User:this.Username,
        Auth:[]ssh.AuthMethod{ssh.Password(this.Password)},
        //需要验证服务端，不做验证返回nil就可以，点击HostKeyCallback看源码就知道了
        HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
            return nil
        },
    }
    this.Log("SSH-START", "address:", this.Address)
    c, err := ssh.Dial("tcp", this.Address, config)
    this.client = c
    if err != nil {
        return err
    }
    return nil
}

//通过ssh执行脚本
//返回依次为standardOut standardErr error
func (this *Ssh)Run(command string) (string, string, error) {
    err := this.connect()
    if err != nil {
        return "", err.Error(), err
    }

    this.Log("SSH-RUN", command)
    s, err := this.client.NewSession()
    defer s.Close()
    if err != nil {
        return "", err.Error(), err
    }

    var outBytes bytes.Buffer
    s.Stdout = &outBytes

    var errBytes bytes.Buffer
    s.Stderr = &errBytes

    err = s.Run(command)
    if err != nil {
        return outBytes.String(), errBytes.String(), err
    }

    return outBytes.String(), errBytes.String(), nil
}

//节点机器是否可以被正常操作
func (this *Node)IsValid() (bool, error) {
    stdOut, stdErr, err := this.Run("cat /etc/os-release")
    if err != nil {
        this.Log("Node-IsValid", err.Error(), stdErr)
        return false, err
    }
    if !strings.Contains(stdOut, `NAME="CentOS Linux"`) || !strings.Contains(stdOut, `VERSION="7 (Core)"`) {
        return false, errors.New("linux not centos 7*")
    }

    return true, nil
}

//获取系统版本号
func (this *Node)GetVersion() (string, error) {
    stdOut, stdErr, err := this.Run("lsb_release -a")
    if err != nil {
        this.Log("Node-GetVersion", err.Error(), stdErr)
        return stdErr, err
    }
    return stdOut, nil
}

//获取系统平均负载信息
func (this *Node)GetLoadAverage() (string, error) {
    stdOut, stdErr, err := this.Run("uptime")
    if err != nil {
        this.Log("Node-GetLoadAverage", err.Error(), stdErr)
        return stdErr, err
    }
    return stdOut, nil
}

//获取系统内存使用情况
func (this *Node)GetMemory() (string, error) {
    stdOut, stdErr, err := this.Run("free")
    if err != nil {
        this.Log("Node-GetLoadAverage", err.Error(), stdErr)
        return stdErr, err
    }
    return stdOut, nil
}

//获取系统硬件基础配置
func (this *Node)GetInfo() (string, error) {
    cmd := `echo [system] && uname -a &&` +
        `echo [disk] && fdisk -l | grep Disk &&` +
        `echo [cpu] && cat /proc/cpuinfo |grep "model name"`
    stdOut, stdErr, err := this.Run(cmd)
    if err != nil {
        this.Log("Node-GetInfo", err.Error(), stdErr)
        return stdErr, err
    }
    return stdOut, nil
}

//系统ping测试
func (this *Node)HostPing(domain, host string) (string, error) {
    if host != "" {
        hostArr := strings.Fields(host)
        if len(hostArr) != 2 || hostArr[1] != domain {
            this.Log("Node-HostPing", "host format error", domain, host)
            return "host error", errors.New("host error")
        }
        ip := hostArr[0]

        stdOut, stdErr, err := this.Run(fmt.Sprintf("ping -c 1 %v", domain))
        if err == nil && strings.Index(stdOut, ip) >= 0 {
            return stdOut, nil
        }

        stdOut, stdErr, err = this.Run(fmt.Sprintf(`echo "%v %v" >> /etc/hosts`, ip, domain))
        if err != nil {
            this.Log("Node-HostPing", err.Error(), stdErr)
            return stdErr, err
        }
        return stdOut, nil
    } else {
        stdOut, stdErr, err := this.Run(fmt.Sprintf("ping -c 1 %v", domain))
        if err != nil {
            return stdErr, err
        }
        return stdOut, nil
    }
}

//添加环境命令
func (this *Node)CmdAdd(name string, value string) (string, error) {
    this.CmdDel(name)

    //进行添加
    cmd := fmt.Sprintf(`echo -e "#!/bin/sh\n%v" >> /usr/local/sbin/%v && chmod 777 /usr/local/sbin/%v`,
        value, name, name)
    stdOut, stdErr, err := this.Run(cmd)
    if err != nil {
        this.Log("Node-CmdAdd", err.Error(), stdErr)
        return stdErr, err
    }
    return stdOut, nil
}

//系统删除/usr/local/sbin/命令
func (this *Node)CmdDel(name string) (string, error) {
    cmd := fmt.Sprintf("rm -rf /usr/local/sbin/%v", name)
    stdOut, stdErr, err := this.Run(cmd)
    if err != nil {
        this.Log("Node-CmdDel", err.Error(), stdErr)
        return stdErr, nil
    }
    return stdOut, nil
}

//系统是否包含某个全局命令
func (this *Node)CmdHas(name string) (bool, error) {
    cmd := fmt.Sprintf("which %v", name)
    _, stdErr, err := this.Run(cmd)
    if err != nil {
        this.Log("Node-CmdHas", err.Error(), stdErr)
        return false, nil
    }
    return true, nil
}

//系统创建目录
func (this *Node)DirMake(path string) (string, error) {
    cmd := fmt.Sprintf(`mkdir -p %v`, path)
    stdOut, stdErr, err := this.Run(cmd)
    if err != nil {
        this.Log("Node-DirMake", err.Error(), stdErr)
        return stdErr, err
    }
    return stdOut, nil
}

//系统删除目录
func (this *Node)DirRm(path string) (string, error) {
    return this.FileRm(path)
}

//系统移动或重命名目录
func (this *Node)DirMv(path string, rename string) (string, error) {
    return this.FileMv(path, rename)
}

//系统目录是否存在
func (this *Node)DirHas(path string) (bool, error) {
    return this.FileHas(path)
}

//系统创建文件
func (this *Node)FileMake(path string) (string, error) {
    cmd := fmt.Sprintf(`touch %v`, path)
    stdOut, stdErr, err := this.Run(cmd)
    if err != nil {
        this.Log("Node-FileMake", err.Error(), stdErr)
        return stdErr, err
    }
    return stdOut, nil
}

//系统移除文件
func (this *Node)FileRm(path string) (string, error) {
    cmd := fmt.Sprintf(`rm -rf %v`, path)
    stdOut, stdErr, err := this.Run(cmd)
    if err != nil {
        this.Log("Node-FileRm", err.Error(), stdErr)
        return stdErr, err
    }
    return stdOut, nil
}

//系统移动或重命名文件
func (this *Node)FileMv(path string, rename string) (string, error) {
    cmd := fmt.Sprintf(`mv %v %v`, path, rename)
    stdOut, stdErr, err := this.Run(cmd)
    if err != nil {
        this.Log("Node-FileMv", err.Error(), stdErr)
        return stdErr, err
    }
    return stdOut, nil
}

//文件是否存在
func (this *Node)FileHas(path string) (bool, error) {
    cmd := fmt.Sprintf(`ls %v`, path)
    _, stdErr, err := this.Run(cmd)
    if err != nil {
        this.Log("Node-FiFileHasleMv", err.Error(), stdErr)
        return false, err
    }
    return true, nil
}

//当前Docker宿主机是否可以正常操作
func (this *DockerNode)IsValid() (bool, error) {
    //该Node节点机是否可用
    valid, err := this.Node.IsValid()
    if !valid {
        return false, err
    }
    //微服务和其备份域必须存在
    if this.RegistryServiceDomain == "" {
        this.Log(
            "DockerNode-IsValid",
            "ServiceRegistryDomain is null",
            this.RegistryServiceDomain,
        )
        return false, errors.New("ServiceRegistryDomain is null")
    }
    //docker必须安装完毕
    _, err = this.DockerInstall()
    if err != nil {
        return false, err
    }
    //微服务镜像仓库连通性检测
    _, err = this.HostPing(this.RegistryServiceDomain, this.RegistryServiceHost)
    if err != nil {
        return false, err
    }
    return true, nil
}

//系统是否安装了docker
func (this *DockerNode)DockerHas() (bool, error) {
    return this.CmdHas("docker")
}

//系统安装docker
func (this *DockerNode)DockerInstall() (string, error) {
    has, err := this.CmdHas("docker")
    if err != nil || !has {
        _, stdErr, err := this.Run("yum -y install docker")
        if err != nil {
            this.Log("DockerNode-DockerInstall", err.Error(), stdErr)
            return stdErr, err
        }
    }

    stdErr, err := this.DockerStart()
    if err != nil {
        return stdErr, err
    }
    return "", nil
}

//启动docker server
func (this *DockerNode)DockerStart() (string, error) {
    _, stdErr1, err1 := this.Run("service docker start")
    _, stdErr2, err2 := this.Run("systemctl start docker")
    _, stdErr3, err3 := this.Run("chkconfig docker on")
    if err1 != nil && err2 != nil {
        this.Log("DockerNode-DockerStart", stdErr1, stdErr2)
        return fmt.Sprintf("service-start %v systemctl-start %v", stdErr1, stdErr2), err2
    }
    if err3 != nil {
        this.Log("DockerNode-DockerStart", stdErr3)
        return fmt.Sprintf("chkconfig %v %v", err3.Error(), stdErr3), nil
    }
    return "docker start success", nil
}

//获取系统docker上具备的所有镜像
func (this *DockerNode)DockerImages() ([]map[string]string, string, error) {
    stdOut, stdErr, err := this.Run("docker images")
    if err != nil {
        this.Log("DockerNode-DockerImages", err.Error(), stdErr)
        return nil, stdErr, err
    }
    var images []map[string]string
    imagesArr := strings.Split(stdOut, "\n")
    imagesArr = imagesArr[1:len(imagesArr) - 1]
    for _, item := range imagesArr {
        itemArr := strings.Fields(item)
        images = append(images, map[string]string{
            "name":itemArr[0],
            "tag":itemArr[1],
            "id":itemArr[2],
            "time":itemArr[3],
            "size":itemArr[4],
        })
    }
    return images, "", nil
}

//获取系统docker上正在运行中的容器
func (this *DockerNode)DockerContainers() ([]map[string]string, string, error) {
    stdOut, stdErr, err := this.Run("docker ps")
    if err != nil {
        this.Log("DockerNode-DockerContainers", err.Error(), stdErr)
        return nil, stdErr, err
    }
    var containers []map[string]string
    containersArr := strings.Split(stdOut, "\n")
    containersArr = containersArr[1:len(containersArr) - 1]
    for _, item := range containersArr {
        itemArr := strings.Fields(item)
        containers = append(containers, map[string]string{
            "id":itemArr[0],
            "image":itemArr[1],
            "status":itemArr[4],
            "name":itemArr[len(itemArr) - 1],
        })
    }
    return containers, "", nil
}

//系统docker是否已经下载了某个镜像
func (this *DockerNode)DockerImageHas(tag string) (bool, string, error) {
    images, stdErr, err := this.DockerImages()
    if err != nil {
        return false, stdErr, err
    }

    for _, image := range images {
        if strings.EqualFold(tag, fmt.Sprintf("%v:%v", image["name"], image["tag"])) {
            return true, "", nil
        }
    }
    return false, "no image", nil
}

//系统docker是否已经运行了某个容器
func (this *DockerNode)DockerContainerHas(name string) (bool, string, error) {
    containers, stdErr, err := this.DockerContainers()
    if err != nil {
        return false, stdErr, err
    }

    for _, container := range containers {
        if container["name"] == name {
            return true, "", nil
        }
    }
    return false, "no container", nil
}

//系统docker拉取镜像
func (this *DockerNode)DockerImagePull(name, version string) (string, error) {
    cmd := fmt.Sprintf(
        `docker login %v --username=%v --password=%v && docker pull %v:%v`,
        this.RegistryServiceDomain,
        this.Username,
        this.Password,
        name, version,
    )
    stdOut, stdErr, err := this.Run(cmd)
    if err != nil {
        this.Log("DockerNode-DockerImagePull", err.Error(), name, version, stdErr)
        return stdErr, err
    }
    return stdOut, nil
}

//系统docker移除镜像
func (this *DockerNode)DockerImageRemove(tag string) (string, error) {
    cmd := fmt.Sprintf("docker rmi -f %v", tag);
    stdOut, stdErr, err := this.Run(cmd)
    if err != nil {
        this.Log("DockerNode-DockerImageRemove", err.Error(), tag, stdErr)
        return stdErr, err
    }
    return stdOut, nil
}

//系统docker创建一个容器
func (this *DockerNode)DockerContainerRun(name, version, image string, ports []int, hosts map[string]string) (string, error) {
    //移除container
    this.DockerContainerRemove(fmt.Sprintf("%v%v", name, version))
    //组装env
    envArr := []string{
        fmt.Sprintf("-e %v=%v", "APP_NAME", name),
        fmt.Sprintf("-e %v=%v", "APP_VERSION", version),
        fmt.Sprintf("-e %v=%v", "APP_MODE", this.AppMode),
        fmt.Sprintf("-e %v=/data/%v%v", "APP_PATH", name, version),
        fmt.Sprintf("-e %v=%v", "NODE_IP", this.Ip),
        fmt.Sprintf("-e %v=%v", "NODE_IP_INNER", this.IpInner),
        fmt.Sprintf("-e %v=%v", "NODE_PORT", ports[0]),
    }
    envString := strings.Join(envArr, " ")
    //组装host
    var hostsArr []string
    if hosts != nil {
        for i, e := range hosts {
            envArr = append(hostsArr, fmt.Sprintf("--add-host %v:%v", i, e))
        }
    }
    hostsString := strings.Join(envArr, " ")
    //组装docker run
    cmd := `docker run -d --restart=always ` +
        `%v %v ` +
        `-p %v:%v ` +
        `-v /data/%v%v:/data ` +
        `--name %v%v %v`
    cmd = fmt.Sprintf(
        cmd,
        envString,
        hostsString,
        ports[0], ports[1],
        name, version,
        name, version, image,
    )
    stdOut, stdErr, err := this.Run(cmd)
    if err != nil {
        this.Log("DockerNode-DockerContainerRun", err.Error(), stdErr)
        return stdErr, err
    }
    return stdOut, nil
}

//系统docker移除容器
func (this *DockerNode)DockerContainerRemove(name string) (string, error) {
    cmd := fmt.Sprintf(`docker rm -f %v`, name);
    stdOut, stdErr, err := this.Run(cmd)
    if err != nil {
        this.Log("DockerNode-DockerContainerRemove", err.Error(), stdErr)
        return stdErr, err
    }
    return stdOut, nil
}

//系统docker显示容器信息
func (this *DockerNode)DockerContainerInfo(name string) (string, error) {
    cmd := fmt.Sprintf(`docker inspect %v`, name);
    stdOut, stdErr, err := this.Run(cmd)
    if err != nil {
        this.Log("DockerNode-DockerContainerInfo", err.Error(), stdErr)
        return stdErr, err
    }
    return stdOut, nil
}

//镜像编译机是否可以进行正常操作
func (this *DockerNodeBuild)IsValid() (bool, error) {
    //是否符合可操作的Docker节点机
    valid, err := this.DockerNode.IsValid()
    if !valid {
        return false, err
    }
    //基础镜像仓库连通性检测
    _, err = this.HostPing(this.RegistryBaseDomain, this.RegistryBaseHost)
    if err != nil {
        return false, err
    }
    //svn & git domain check
    if this.SvnDomain == "" && this.GitDomain == "" {
        this.Log("DockerNodeBuild-IsValid", "svn or git domain is null")
        return false, errors.New("svn or git domain is null")
    }
    //install and linked check
    if this.SvnDomain != "" {
        has, err := this.SvnInstall()
        if !has {
            return false, err
        }
        linked, err := this.SvnLinked()
        if !linked {
            return false, err
        }
    } else if this.GitDomain != "" {
        fmt.Println("=============", this.GitDomain)
        has, err := this.GitInstall()
        if !has {
            return false, err
        }
        linked, err := this.GitLinked()
        if !linked {
            return false, err
        }
    }
    return true, nil
}

//系统是否安装了svn客户端
func (this *DockerNodeBuild)SvnHas() (bool, error) {
    cmd := "which svn"
    _, stdErr, err := this.Run(cmd)
    if err != nil {
        this.Log("DockerNodeBuild-SvnHas", err.Error(), stdErr)
        return false, nil
    }
    return true, nil
}

//svn服务连通性检测
func (this *DockerNodeBuild)SvnLinked() (bool, error) {
    stdErr, err := this.HostPing(this.SvnDomain, this.SvnHost)
    if err != nil {
        this.Log("DockerNodeBuild-SvnLinked", err.Error(), stdErr)
        return false, err
    }
    return true, nil
}

//系统安装svn客户端
func (this *DockerNodeBuild)SvnInstall() (bool, error) {
    has, err := this.SvnHas()
    if has {
        return true, nil
    }

    cmd := "yum -y install svn"
    _, stdErr, err := this.Run(cmd)
    if err != nil {
        this.Log("DockerNodeBuild-SvnInstall", err.Error(), stdErr)
        return false, err
    }
    return true, nil
}

//系统svn下载项目
//成功后返回src仓库的uuid
func (this *DockerNodeBuild)SvnCheckout(src string) (string, error) {
    uuid := utils.GetUUID()
    cmd := fmt.Sprintf(`mkdir -p /opt/building/%v && cd /opt/building/%v && svn export %v %v ` +
        `--username %v --password %v --no-auth-cache --non-interactive`,
        uuid, uuid, src, uuid, this.SvnUsername, this.SvnPassword)
    _, stdErr, err := this.Run(cmd)
    if err != nil {
        this.Log("DockerNodeBuild-SvnCheckout", err.Error(), stdErr)
        return stdErr, err
    }
    return uuid, nil
}


//系统是否安装了git客户端
func (this *DockerNodeBuild)GitHas() (bool, error) {
    cmd := "which git"
    _, stdErr, err := this.Run(cmd)
    if err != nil {
        this.Log("DockerNodeBuild-GitHas", err.Error(), stdErr)
        return false, nil
    }
    return true, nil
}


//系统安装git客户端
func (this *DockerNodeBuild)GitInstall() (bool, error) {
    has, err := this.GitHas()
    if has {
        return true, nil
    }

    cmd := "yum -y install git"
    _, _, err = this.Run(cmd)
    if err != nil {
        return false, nil
    }
    return true, nil
}

//git服务连通性测试
func (this *DockerNodeBuild)GitLinked() (bool, error) {
    stdErr, err := this.HostPing(this.GitDomain, this.GitHost)
    if err != nil {
        this.Log("DockerNodeBuild-GitLinked", err.Error(), stdErr)
        return false, err
    }
    return true, nil
}

//系统git下载项目
//成功后返回src仓库的uuid
func (this *DockerNodeBuild)GitClone(src string) (string, error) {
    has, err := this.GitHas()
    if err != nil || !has {
        return "no git", errors.New("no git")
    }

    linked, err := this.GitLinked()
    if err != nil || !linked {
        return "git no linked", errors.New("git no linked")
    }

    uuid := utils.GetUUID()
    cmd := fmt.Sprintf(`mkdir -p /opt/sources && cd /opt/sources
	&& git clone http://%v:%v@%v %v`, this.GitUsername, this.GitPassword, src, uuid)
    _, stdErr, err := this.Run(cmd)
    if err != nil {
        this.Log("DockerNodeBuild-GitClone", err.Error(), stdErr)
        return stdErr, err
    }
    return uuid, nil
}

//系统docker开始构建微服务镜像
//name 微服务名称
//version 微服务版本号
//image 基础镜像名称(image:tag)
//src 微服务项目svn或git地址
//镜像构建version一定要是累加的不可重复
func (this *DockerNodeBuild)DockerBuildAndPush(name, version, image, src string) (string, error) {
    var err error
    var srcUuid string
    if strings.Index(src, "svn") >= 0 {
        srcUuid, err = this.SvnCheckout(src)
    } else {
        srcUuid, err = this.GitClone(src)
    }
    if err != nil {
        return srcUuid, err
    }

    //微服务全称
    msName := fmt.Sprintf("%v%v", name, version);
    //Dockerfile
    var dockerFileBuffer bytes.Buffer
    dockerFileBuffer.WriteString("FROM %v/%v\n")
    dockerFileBuffer.WriteString("MAINTAINER dev@foryou56.com\n")
    dockerFileBuffer.WriteString("RUN mkdir -p /opt/%v\n")
    dockerFileBuffer.WriteString("RUN mkdir -p /opt/run/ && touch /opt/run/run.log\n")
    dockerFileBuffer.WriteString("COPY . /opt/%v/\n")
    dockerFileBuffer.WriteString("CMD sh /opt/%v/run.sh && tail -f /opt/run/run.log")
    dockerFile := fmt.Sprintf(dockerFileBuffer.String(), this.RegistryBaseDomain, image, msName, msName, msName)
    //image build
    dockerFileCmd := `docker login %v --username=%v --password=%v ` +
        `&& echo "%v" > /opt/building/%v/Dockerfile ` +
        `&& cd /opt/building/%v ` +
        `&& docker build -t %v/%v:%v . ` +
        `&& rm -rf /opt/building/%v`
    dockerFileCmd = fmt.Sprintf(
        dockerFileCmd,
        this.RegistryBaseDomain, this.RegistryUsername, this.RegistryPassword,
        dockerFile,
        srcUuid,
        srcUuid,
        this.RegistryServiceDomain, name, version,
        srcUuid,
    )
    _, stdErr, err := this.Run(dockerFileCmd)
    if err != nil {
        this.Log("DockerNodeBuild-DockerImageBuild-Build", err.Error(), stdErr)
        return stdErr, err
    }
    //image push
    var pushCmd []string
    registryList := []string{this.RegistryServiceDomain}
    for _, registry := range registryList {
        if registry == "" {
            continue
        }
        pushCmd = append(
            pushCmd,
            fmt.Sprintf(
                `docker login %v --username=%v --password=%v && docker push %v/%v:%v && docker rmi -f %v/%v:%v`,
                registry,
                this.RegistryUsername,
                this.RegistryPassword,
                registry, name, version,
                registry, name, version,
            ),
        )
    }
    _, stdErr, err = this.Run(strings.Join(pushCmd, " && "))
    if err != nil {
        this.Log("DockerNodeBuild-DockerImageBuild-Push", err.Error(), stdErr)
        return stdErr, err
    }

    return srcUuid, nil
}

//移除仓库机的某个版本镜像
//逻辑删除+物理删除镜像
//仓库请使用latest版本
func (this *DockerNodeRegistry)DeleteRegistryImage(name, version string) (bool, error) {
    ok, err := this.Registry.HasImage(name, version)
    if !ok {
        this.Log("DockerNodeRegistry-DeleteRegistryImage", err.Error(), "Registry no", name, version)
        return false, err
    }

    //image logic delete
    ok, err = this.Registry.DeleteImage(name, version)
    if !ok || err != nil {
        this.Log("DockerNodeRegistry-DeleteRegistryImage", err.Error(), "Registry delete", name, version)
        return false, err
    }

    //registry garbage-collect
    cmd := "docker exec -i %v bin/registry garbage-collect /etc/docker/registry/config.yml"
    _, stdErr, err := this.Run(fmt.Sprintf(cmd, this.RegistryDockerRunName))
    if err != nil {
        this.Log("DockerNodeRegistry-DeleteRegistryImage", err.Error(), stdErr, "DockerRegistry Collect Error", name, version)
        return false, err
    }

    return true, nil
}