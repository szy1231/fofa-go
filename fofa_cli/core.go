package fofa_cli

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"github.com/xiaoyu-0814/fofa-go/fofa"

	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"
)

var clt *fofa.Fofa

const version = "1.0.0"

//setClt Init fofa client
func setClt() {
	var err error
	*email, *key, err = getKey()
	fofaErr(err)
	clt = fofa.NewFofaClient([]byte(*email), []byte(*key))
}

//getInfo Get user info
func getInfo() {
	userInfo, err := clt.UserInfo()
	fofaErr(err)
	if userInfo.Err != "" {
		fofaErr(errors.New(userInfo.Err))
	}
	fmt.Printf("Email：%s\nUserName：%s\nFcoin：%d\nVip：%t\nVipLevel：%d\n", userInfo.Email,
		userInfo.UserName, userInfo.Fcoin, userInfo.Vip, userInfo.VipLevel)
	os.Exit(0)
}

//search Fofa search
func search() {
	if *page <= 0 {
		*page = 1
	}

	if *query == "" {
		usage()
		os.Exit(0)
	}

	*query = strings.ReplaceAll(*query, "+", "&&")
	*query = strings.ReplaceAll(*query, "-", "||")
	*query = queryDomainParse(*query)

	fofaData, err := clt.QueryAsObject(uint(*page), []byte(*query), []byte(*fields))
	fofaErr(err)
	if *count {
		fmt.Printf("\ntotal: %d", fofaData.Size)
		os.Exit(0)
	}
	if *out != "" {
		*out = *out + ".txt"
		file, err := os.Create(*out)
		fofaErr(err)

		writer := bufio.NewWriter(file)
		for _, res := range fofaData.Results {
			str := ""
			resLen := len(res)
			for k, v := range res {
				if v == "" {
					v = "nil"
				}
				if k == resLen-1 {
					str += v
					continue
				}
				str += v + *format
			}
			writer.Write([]byte(fmt.Sprintf("%s\n", str)))
		}
		writer.Write([]byte(fmt.Sprintf("\ntotal: %d", fofaData.Size)))
		writer.Flush()
		fmt.Println("[+] Successfully")
	} else {
		for _, res := range fofaData.Results {
			str := ""
			resLen := len(res)
			for k, v := range res {
				if v == "" {
					v = "nil"
				}
				if k == resLen-1 {
					str += v
					continue
				}
				str += v + *format
			}
			fmt.Printf("%s\n", str)
		}
		fmt.Printf("\ntotal: %d", fofaData.Size)
	}
}

// Parse parses the command-line flags from os.Args[2:]. Must be called
// after all flags are defined and before flags are accessed by the program.
func parse() {
	// Ignore errors; CommandLine is set for ExitOnError.
	flag.CommandLine.Parse(os.Args[2:])
}

//userInit Init user email and apikey
func userInit() {
	if *email == "" || *key == "" {
		fofaErr(errors.New("please input your email and key\nfofa_cli init -email example@fofa.so -key 32charsMD5String"))
	}

	clt = fofa.NewFofaClient([]byte(*email), []byte(*key))
	if clt == nil {
		fofaErr(errors.New("Allocate Failed! Out Of Memery!"))
	}
	userInfo, err := clt.UserInfo()
	fofaErr(err)
	fofaPath, keyPath, err := getPath()
	fofaErr(err)
	if !fileExist(fofaPath) {
		err := os.MkdirAll(fofaPath, 0666)
		fofaErr(err)
	}

	err = ioutil.WriteFile(keyPath, []byte(*email+"------"+*key), 0666)
	fofaErr(err)

	fmt.Println("[+] Successfully initialized")
	fmt.Printf("\nEmail：%s\nUserName：%s\nFcoin：%d\nVip：%t\nVipLevel：%d\n", userInfo.Email,
		userInfo.UserName, userInfo.Fcoin, userInfo.Vip, userInfo.VipLevel)
}

//getKey Get user email and apikey
func getKey() (email, key string, err error) {
	var (
		keyPath string
		data    []byte
	)
	_, keyPath, err = getPath()
	if err != nil {
		return
	}

	data, err = ioutil.ReadFile(keyPath)
	if err != nil {
		return
	}

	dataArr := strings.Split(string(data), "------")
	if len(dataArr) == 2 {
		return dataArr[0], dataArr[1], nil
	}
	err = errors.New("email or apikey get fail")
	return
}

//getPath Get home path
func getPath() (string, string, error) {
	user, err := user.Current()
	if nil == err {
		return setFofaPath(user.HomeDir, nil)
	}

	// cross compile support

	if "windows" == runtime.GOOS {
		windowsHome, err := homeWindows()
		return setFofaPath(windowsHome, err)
	}

	// Unix-like system, so just assume Unix
	home, err := homeUnix()
	if err == nil {
		return setFofaPath(home, nil)
	}

	return setFofaPath("", err)
}

//setFofaPath Set fofa dir
func setFofaPath(home string, err error) (string, string, error) {
	if err != nil {
		return "", "", err
	}
	fofaPath := home + "/.config/fofa/setting"
	keyPath := home + "/.config/fofa/setting/apikey"

	return fofaPath, keyPath, nil
}

//homeUnix Unix-like system, so just assume Unix
func homeUnix() (string, error) {
	// First prefer the HOME environmental variable
	if home := os.Getenv("HOME"); home != "" {
		return home, nil
	}

	// If that fails, try the shell
	var stdout bytes.Buffer
	cmd := exec.Command("sh", "-c", "eval echo ~$USER")
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return "", err
	}

	result := strings.TrimSpace(stdout.String())
	if result == "" {
		return "", errors.New("blank output when reading home directory")
	}

	return result, nil
}

func homeWindows() (string, error) {
	drive := os.Getenv("HOMEDRIVE")
	path := os.Getenv("HOMEPATH")
	home := drive + path
	if drive == "" || path == "" {
		home = os.Getenv("USERPROFILE")
	}
	if home == "" {
		return "", errors.New("HOMEDRIVE, HOMEPATH, and USERPROFILE are blank")
	}

	return home, nil
}

//fileExist Check to see if the file exists
func fileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

//cliVersion fmt cli version
func cliVersion() {
	fmt.Printf("\nVersion：%s", version)
}

//queryDomainParse Domain Exceptions
func queryDomainParse(query string) string {
	queryArr := strings.Split(query, "&&")
	for _, v := range queryArr {
		queryArr2 := strings.Split(v, "||")
		for _, v1 := range queryArr2 {
			queryLen := len(v1)
			if queryLen > 7 && v1[:7] == "domain=" {
				query = strings.ReplaceAll(query, v1[7:], `"`+v1[7:]+`"`)
			} else if queryLen > 8 && v1[:8] == "domain!=" {
				query = strings.ReplaceAll(query, v1[8:], `"`+v1[8:]+`"`)
			}
		}
	}
	return query
}
