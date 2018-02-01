package models

import (
    "testing"
    "fmt"
)

const (
    SSH_USERNAME = "root"
    SSH_PASSWORD = "qwe2#$%75$Ghhhhh"
    SSH_ADDRESS = "47.94.42.154:22"

    NODE_IP = "47.94.42.154"
    NODE_IP_INNER = ""

    NODE_SVN_USERNAME = "linyang"
    NODE_SVN_PASSWORD = "xyq2525307"
    NODE_SVN_DOMAIN = "svn.fuyoukache.com"
    NODE_SVN_HOSTS = "119.2.6.189 svn.fuyoukache.com"

    NODE_GIT_USERNAME = ""
    NODE_GIT_PASSWORD = ""
    NODE_GIT_DOMAIN = ""
    NODE_GIT_HOSTS = ""

    NODE_APP_MODE = "develop"
)

var node *Node = &Node{
    Ssh:Ssh{Address:SSH_ADDRESS, Username:SSH_USERNAME, Password:SSH_PASSWORD},
    Ip:NODE_IP,
    IpInner:NODE_IP_INNER,
}

var dockerNode = DockerNode{
    Node:Node{
        Ssh:Ssh{Address:SSH_ADDRESS, Username:SSH_USERNAME, Password:SSH_PASSWORD},
        Ip:NODE_IP,
        IpInner:NODE_IP_INNER,
    },
    AppMode :NODE_APP_MODE,
    RegistryUsername:"linyang",
    RegistryPassword:"123456",
    RegistryServiceDomain:"registry-service.fuyoukache.com",
    RegistryServiceHost:"127.0.0.1 registry-service.fuyoukache.com",
}

//docker宿主机
var docker *DockerNode = &dockerNode
//docker编译机
var build *DockerNodeBuild = &DockerNodeBuild{
    DockerNode:dockerNode,
    RegistryBaseDomain:"registry-base.fuyoukache.com",
    RegistryBaseHost:"127.0.0.1 registry-base.fuyoukache.com",
    SvnUsername:NODE_SVN_USERNAME,
    SvnPassword:NODE_SVN_PASSWORD,
    SvnDomain:NODE_SVN_DOMAIN,
    SvnHost:NODE_SVN_HOSTS,
}
//docker任务机
var task *DockerNodeTask = &DockerNodeTask{
    DockerNode:dockerNode,
    DefaultImages:[]string{"yaf:1.0.0"},
}
//docker仓库机
var registry *DockerNodeRegistry = &DockerNodeRegistry{
    DockerNode:dockerNode,
    Registry:Registry{Username:docker.RegistryUsername, Password:docker.RegistryPassword, Domain:docker.RegistryServiceDomain},
    RegistryDockerRunName:"registry-service",
}

//ssh run
func TestConnect_Run(t *testing.T) {
    ssh := NewSsh(SSH_ADDRESS, SSH_USERNAME, SSH_PASSWORD)
    defer ssh.Close()
    stdOut, stdErr, err := ssh.Run(`echo "hello world"`)
    if err != nil {
        t.Error(stdErr)
    } else {
        t.Log(stdOut)
    }
    fmt.Println(ssh.Logs())
}

//ssh close
func TestConnect_Close(t *testing.T) {
    ssh := NewSsh(SSH_ADDRESS, SSH_USERNAME, SSH_PASSWORD)
    _, stdErr, err := ssh.Run(`echo "hello world"`)
    if err != nil {
        t.Error(stdErr, err.Error())
        return
    }
    err = ssh.Close()
    if err != nil {
        t.Error(err.Error())
    } else {
        t.Log("success")
    }
    fmt.Println(ssh.Logs())
}

func TestNode_GetInfo(t *testing.T) {
    valid, err := node.IsValid()
    if !valid {
        t.Error(valid, err.Error())
    } else {
        out, err := node.GetInfo()
        if err != nil {
            t.Error(out, err.Error())
        } else {
            fmt.Println(out)
        }
    }
    fmt.Println(node.Logs())
}

func TestNode_GetLoadAverage(t *testing.T) {
    valid, err := node.IsValid()
    if !valid {
        t.Error(valid, err.Error())
    } else {
        out, err := node.GetLoadAverage()
        if err != nil {
            t.Error(out, err.Error())
        } else {
            fmt.Println(out)
        }
    }
    fmt.Println(node.Logs())
}

func TestNode_GetMemory(t *testing.T) {
    valid, err := node.IsValid()
    if !valid {
        t.Error(valid, err.Error())
    } else {
        out, err := node.GetMemory()
        if err != nil {
            t.Error(out, err.Error())
        } else {
            fmt.Println(out)
        }
    }
    fmt.Println(node.Logs())
}

func TestNode_GetVersion(t *testing.T) {
    valid, err := node.IsValid()
    if !valid {
        t.Error(valid, err.Error())
    } else {
        out, err := node.GetVersion()
        if err != nil {
            t.Error(out, err.Error())
        } else {
            fmt.Println(out)
        }
    }
    fmt.Println(node.Logs())
}

func TestNode_HostPing(t *testing.T) {
    valid, err := node.IsValid()
    if !valid {
        t.Error(valid, err.Error())
        fmt.Println(node.Logs())
        return
    }
    out, err := node.HostPing("www.163.com", "")
    if err != nil {
        t.Error(out, err.Error())
    } else {
        fmt.Println(out)
    }
    fmt.Println(node.Logs())
}

func TestNode_CmdHas(t *testing.T) {
    valid, err := node.IsValid()
    if !valid {
        t.Error(valid, err.Error())
        fmt.Println(node.Logs())
        return
    }
    out, err := node.CmdHas("linyang")
    if err != nil {
        t.Error(out, err.Error())
    } else {
        fmt.Println(out)
    }
    fmt.Println(node.Logs())
}

func TestNode_CmdAdd(t *testing.T) {
    valid, err := node.IsValid()
    if !valid {
        t.Error(valid, err.Error())
        fmt.Println(node.Logs())
        return
    }
    out, err := node.CmdAdd("linyang", "ls /")
    if err != nil {
        t.Error(out, err.Error())
    } else {
        fmt.Println(out)
        stdOut, stdErr, err := node.Run("linyang")
        fmt.Println(stdOut, stdErr, err)
    }
    fmt.Println(node.Logs())
}

func TestNode_DirHas(t *testing.T) {
    valid, err := node.IsValid()
    if !valid {
        t.Error(valid, err.Error())
        fmt.Println(node.Logs())
        return
    }
    _, err = node.DirHas("/opt/sources")
    fmt.Println(err)
    fmt.Println(node.Logs())
}

func TestNode_DirRm(t *testing.T) {
    out, err := node.DirRm("/opt/sources")
    if err != nil {
        t.Error("dir rm error", out, err.Error())
    } else {
        fmt.Println("dir rm")
        t.Log(out)
    }
}

func TestNode_DirMake(t *testing.T) {
    out, err := node.DirMake("/opt/sources")
    if err != nil {
        t.Error("dir mk has", out, err.Error())
    } else {
        fmt.Println("dir mk")
        t.Log(out)
    }
}

func TestNode_SvnHas(t *testing.T) {
    has, err := build.SvnHas()
    fmt.Println("SvnHas", has, err)
    t.Log(has)
}

func TestNode_SvnInstall(t *testing.T) {
    has, err := build.SvnInstall()
    if err != nil || !has {
        t.Error("svn install error", err.Error())
    } else {
        fmt.Println("svn installed")
        t.Log(has)
    }
}

func TestNode_SvnCheckout(t *testing.T) {
    out, err := build.SvnCheckout("http://svn.fuyoukache.com/tp/trunk/yaf")
    if err != nil {
        t.Error("[TestNode_SvnCheckout]", out, err.Error())
    } else {
        fmt.Println("TestNode_SvnCheckout", out)
        t.Log(out)

        TestNode_DirRm(t)
    }
}

func TestNode_SvnLinked(t *testing.T) {
    out, err := build.SvnLinked()
    fmt.Println("TestNode_SvnLinked", out, err)
    t.Log(out)
}

func TestDockerNode_DockerHas(t *testing.T) {
    out, err := docker.DockerHas()
    fmt.Println("TestDockerNode_DockerHas", out, err)
    t.Log(out)
}

func TestDockerNode_DockerInstall(t *testing.T) {
    out, err := docker.DockerInstall()
    if err != nil {
        t.Error("TestDockerNode_DockerInstall", out, err)
    } else {
        fmt.Println("TestDockerNode_DockerInstall", out)
        t.Log(out)
    }
}

func TestDockerNode_DockerStart(t *testing.T) {
    out, err := docker.DockerStart()
    if err != nil {
        t.Error("TestDockerNode_DockerStart", out, err)
    } else {
        fmt.Println("TestDockerNode_DockerStart", out, err)
        t.Log(out)
    }
}

func TestDockerNode_DockerImages(t *testing.T) {
    images, out, err := docker.DockerImages()
    if err != nil {
        t.Error("TestDockerNode_DockerImages", images, out, err)
    } else {
        fmt.Println("TestDockerNode_DockerImages", images, out)
        t.Log(out)
    }
}

func TestDockerNode_DockerImageHas(t *testing.T) {
    has, out, err := docker.DockerImageHas("docker.io/nginx:latest")
    if err != nil {
        t.Error("TestDockerNode_DockerImageHas", out, err)
    } else {
        fmt.Println("TestDockerNode_DockerImageHas", has, out, err)
        t.Log(out)
    }
}

func TestDockerNode_DockerImagePull(t *testing.T) {
    out, err := docker.DockerImagePull("alpine", "latest")
    if err != nil {
        t.Error("TestDockerNode_DockerImagePull", out, err)
    } else {
        fmt.Println("TestDockerNode_DockerImagePull", out, err)
        t.Log(out)
    }
}

func TestDockerNodeBuild_DockerBuildAndPush(t *testing.T) {
    valid, err := build.IsValid()
    if !valid {
        t.Error("TestDockerNodeBuild_DockerBuildAndPush", err.Error())
        fmt.Println(build.Logs())
        return
    }
    out, err := build.DockerBuildAndPush("test", "latest", "alpine", "http://svn.fuyoukache.com/tp/trunk/yaf")
    if err != nil {
        t.Error("TestDockerNodeBuild_DockerBuildAndPush", out, err)
    } else {
        fmt.Println("TestDockerNodeBuild_DockerBuildAndPush", out, err)
    }
    fmt.Println(build.Logs())
}

func TestDockerNode_DockerImageRemove(t *testing.T) {
    out, err := docker.DockerImageRemove("registry-service.fuyoukache.com/yaf:1.0.0")
    if err != nil {
        t.Error("TestDockerNode_DockerImageRemove", out, err)
    } else {
        fmt.Println("TestDockerNode_DockerImageRemove", out, err)
        t.Log("TestDockerNode_DockerImageRemove", out)
    }
}

func TestDockerNode_DockerContainers(t *testing.T) {
    containers, out, err := docker.DockerContainers()
    if err != nil {
        t.Error("TestDockerNode_DockerContainers", containers, out, err)
    } else {
        fmt.Println("TestDockerNode_DockerContainers", containers, out)
        t.Log(out)
    }
}

func TestDockerNode_DockerContainerHas(t *testing.T) {
    has, out, err := docker.DockerContainerHas("nginx")
    if err != nil {
        t.Error("TestDockerNode_DockerContainerHas", out, err)
    } else {
        fmt.Println("TestDockerNode_DockerContainerHas", has, out, err)
        t.Log(out)
    }
}

func TestDockerNode_DockerContainerInfo(t *testing.T) {
    out, err := docker.DockerContainerInfo("nginx")
    if err != nil {
        t.Error("TestDockerNode_DockerContainerInfo", out, err)
    } else {
        fmt.Println("TestDockerNode_DockerContainerInfo", out, err)
        t.Log(out)
    }
}

func TestDockerNode_DockerContainerRun(t *testing.T) {
    out, err := docker.DockerContainerRun("y", "1.0.0", "registry-service.fuyoukache.com/yaf:v1.0.0", []int{1234, 10000}, nil)
    if err != nil {
        t.Error("TestDockerNode_DockerContainerRun", out, err)
    } else {
        fmt.Println("TestDockerNode_DockerContainerRun", out, err)
        t.Log(out)
    }
}

func TestDockerNode_DockerContainerRemove(t *testing.T) {
    out, err := docker.DockerContainerRemove("yaf1.0.0")
    if err != nil {
        t.Error("TestDockerNode_DockerContainerRemove", out, err)
    } else {
        fmt.Println("TestDockerNode_DockerContainerRemove", out, err)
        t.Log(out)
    }
}

func TestDockerNodeRegistry_DeleteRegistryImage(t *testing.T) {
    out, err := registry.DeleteRegistryImage("yaf", "0.0.3")
    if err != nil {
        t.Error("TestDockerNodeRegistry_DeleteRegistryImage", out, err)
    } else {
        fmt.Println("TestDockerNodeRegistry_DeleteRegistryImage", out, err)
        t.Log(out)
    }
}