package goincv

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"

	fastJson "github.com/goccy/go-json"
)

func JsonToFile(file string, v interface{}) bool {
	data, err := fastJson.Marshal(v)
	if err != nil {
		return false
	}
	return ioutil.WriteFile(file, data, 0644) == nil
}

func ReadFileContent(path string) []byte {
	r, _ := ioutil.ReadFile(path)
	return r
}

// 递归获取指定目录下的所有文件名.
func ReadFilesInDir(pathname string) []string {
	result := []string{}

	fis, err := ioutil.ReadDir(pathname)
	if err != nil {
		fmt.Printf("读取文件目录失败，pathname=%v, err=%v \n", pathname, err)
		return result
	}

	// 所有文件/文件夹
	for _, fi := range fis {
		fullname := filepath.Join(pathname, fi.Name())
		// 是文件夹则递归进入获取;是文件，则压入数组
		if fi.IsDir() {
			temp := ReadFilesInDir(fullname)
			result = append(result, temp...)
		} else {
			result = append(result, fullname)
		}
	}

	return result
}
func RunFuncName(i int) string {
	pc := make([]uintptr, 1)
	runtime.Callers(i, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}

var TimeConsumingOutput = ""

type TimeConsumingType string

const (
	TimeConsumingModeNone   TimeConsumingType = "None"
	TimeConsumingModeStdout TimeConsumingType = "Stdout"
	TimeConsumingModeFile   TimeConsumingType = "File"
)

var TimeConsumingModeWithName map[string]TimeConsumingType = map[string]TimeConsumingType{}

func TimeConsuming(names ...string) func() {
	name := RunFuncName(3)
	if len(names) > 0 {
		name = names[0]
	}
	unt := time.Now().UnixNano()
	return func() {
		end := time.Now().UnixNano()
		switch TimeConsumingModeWithName[name] {
		case TimeConsumingModeStdout:
			log.Println(name, "耗时:", (end-unt)/(1000*1000), "ms")
		case TimeConsumingModeFile:
			if TimeConsumingOutput != "" {
				old := ReadFileContent(TimeConsumingOutput)
				old = []byte(fmt.Sprintln(string(old), "\n", name, "耗时:", (end-unt)/(1000*1000), "ms"))
				ioutil.WriteFile(TimeConsumingOutput, old, 0644)
			}
		}
	}
}

func SftpUpload(user, password, host, port string, remoteFile string, data []byte) error {

	var connect = func(user, password, host, port string) (*sftp.Client, error) {
		var (
			auth         []ssh.AuthMethod
			addr         string
			clientConfig *ssh.ClientConfig
			sshClient    *ssh.Client
			sftpClient   *sftp.Client
			err          error
		)
		// get auth method
		auth = make([]ssh.AuthMethod, 0)
		auth = append(auth, ssh.Password(password))

		clientConfig = &ssh.ClientConfig{
			User:    user,
			Auth:    auth,
			Timeout: 30 * time.Second,
			HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				return nil
			},
		}

		// connet to ssh
		addr = fmt.Sprintf("%s:%s", host, port)

		if sshClient, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
			return nil, err
		}

		// create sftp client
		if sftpClient, err = sftp.NewClient(sshClient); err != nil {
			return nil, err
		}

		return sftpClient, nil
	}
	// 这里换成实际的 SSH 连接的 用户名，密码，主机名或IP，SSH端口
	sftpClient, err := connect(user, password, host, port)
	if err != nil {
		return err
	}
	defer sftpClient.Close()

	dstFile, err := sftpClient.Create(remoteFile)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = dstFile.Write(data)

	return err
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func WriteFileWithURI(uri string, buff []byte) error {
	return WriteFileWithURICallback(uri, func(tmpFile string) error {
		return ioutil.WriteFile(tmpFile, buff, 0644)
	})
}
func WriteFileWithURICallback(uri string, pre func(tmpFile string) error) error {
	fileExt := filepath.Ext(uri)
	url, _ := url.Parse(uri)
	if url == nil {
		return pre(uri)
	}
	switch url.Scheme {
	case "sftp":
		pwd, _ := url.User.Password()
		user := url.User.Username()
		port := url.Port()
		if port == "" || port == "0" {
			port = "22"
		}
		log.Println("url.Scheme:", url.Scheme, url.User.Username(), pwd, port)
		rand.Seed(time.Now().Unix())
		tmpFile := filepath.Join(os.TempDir(), fmt.Sprintf("%d_%d%s", time.Now().UnixNano(), rand.Int(), fileExt))
		err := pre(tmpFile)
		if err != nil {
			return err
		}
		tmpContent := ReadFileContent(tmpFile)
		if !PathExists(tmpFile) {
			return fmt.Errorf("文件写入失败")
		}
		// if len(tmpContent) == 0 {
		// 	return fmt.Errorf("文件为空")
		// }
		err = SftpUpload(user, pwd, url.Hostname(), port, url.Path, tmpContent)
		if err != nil {
			return err
		}
	default:
		return pre(uri)
	}
	return nil
}

func InterfaceReload(in interface{}, out interface{}) {
	str, _ := fastJson.Marshal(in)
	fastJson.Unmarshal(str, out)
}

func DeepCopy(in interface{}, out interface{}) {
	str, _ := fastJson.Marshal(in)
	fastJson.Unmarshal(str, out)
}
