package api

import (
	"bufio"
	"clashconfig/util"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"

	"gopkg.in/yaml.v2"
)

type Vmess struct {
	Add  string      `json:"add"`
	Aid  interface{} `json:"aid"`
	Host string      `json:"host"`
	ID   string      `json:"id"`
	Net  string      `json:"net"`
	Path string      `json:"path"`
	Port interface{} `json:"port"`
	PS   string      `json:"ps"`
	TLS  string      `json:"tls"`
	Type string      `json:"type"`
	V    string      `json:"v"`
}

type ClashVmess struct {
	Name           string            `json:"name,omitempty"`
	Type           string            `json:"type,omitempty"`
	Server         string            `json:"server,omitempty"`
	Port           interface{}       `json:"port,omitempty"`
	UUID           string            `json:"uuid,omitempty"`
	AlterID        interface{}       `json:"alterId,omitempty"`
	Cipher         string            `json:"cipher,omitempty"`
	TLS            bool              `json:"tls,omitempty"`
	Network        string            `json:"network,omitempty"`
	WSPATH         string            `json:"ws-path,omitempty"`
	WSHeaders      map[string]string `json:"ws-headers,omitempty"`
	SkipCertVerify bool              `json:"skip-cert-verify,omitempty"`
}

type ClashRSSR struct {
	Name          string      `json:"name"`
	Type          string      `json:"type"`
	Server        string      `json:"server"`
	Port          interface{} `json:"port"`
	Password      string      `json:"password"`
	Cipher        string      `json:"cipher"`
	Protocol      string      `json:"protocol"`
	ProtocolParam string      `json:"protocolparam"`
	OBFS          string      `json:"obfs"`
	OBFSParam     string      `json:"obfsparam"`
}

type Clash struct {
	Port      int `yaml:"port"`
	SocksPort int `yaml:"socks-port"`
	// RedirPort          int                      `yaml:"redir-port"`
	// Authentication     []string                 `yaml:"authentication"`
	AllowLan           bool   `yaml:"allow-lan"`
	Mode               string `yaml:"mode"`
	LogLevel           string `yaml:"log-level"`
	ExternalController string `yaml:"external-controller"`
	// ExternalUI         string                   `yaml:"external-ui"`
	// Secret             string                   `yaml:"secret"`
	// Experimental       map[string]interface{} 	`yaml:"experimental"`
	Proxy             []map[string]interface{} `yaml:"Proxy"`
	ProxyGroup        []map[string]interface{} `yaml:"Proxy Group"`
	Rule              []string                 `yaml:"Rule"`
	CFWByPass         []string                 `yaml:"cfw-bypass"`
	CFWLatencyTimeout int                      `yaml:"cfw-latency-timeout"`
}

func (this *Clash) LoadTemplate(path string, protos []interface{}) []byte {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		log.Printf("[%s] template doesn't exist.", path)
		return nil
	}
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("[%s] template open the failure.", path)
		return nil
	}
	err = yaml.Unmarshal(buf, &this)
	if err != nil {
		log.Printf("[%s] Template format error.", path)
	}

	this.Proxy = nil

	var proxys []map[string]interface{}
	var proxies []string

	for _, proto := range protos {
		o := reflect.ValueOf(proto)
		nameField := o.FieldByName("Name")
		proxy := make(map[string]interface{})
		j, _ := json.Marshal(proto)
		json.Unmarshal(j, &proxy)
		proxys = append(proxys, proxy)
		this.Proxy = append(this.Proxy, proxy)
		proxies = append(proxies, nameField.String())
	}

	this.Proxy = proxys

	for _, group := range this.ProxyGroup {
		groupProxies := group["proxies"].([]interface{})
		for i, proxie := range groupProxies {
			if "1" == proxie {
				groupProxies = groupProxies[:i]
				var tmpGroupProxies []string
				for _, s := range groupProxies {
					tmpGroupProxies = append(tmpGroupProxies, s.(string))
				}
				tmpGroupProxies = append(tmpGroupProxies, proxies...)
				group["proxies"] = tmpGroupProxies
				break
			}
		}

	}

	d, err := yaml.Marshal(this)
	if err != nil {
		return nil
	}

	return d
}

func V2ray2Clash(c *gin.Context) {

	decodeBody := c.GetString("decodebody")

	scanner := bufio.NewScanner(strings.NewReader(decodeBody))
	var vmesss []interface{}
	for scanner.Scan() {
		if !strings.HasPrefix(scanner.Text(), "vmess://") {
			continue
		}
		s := scanner.Text()[8:]
		s = strings.TrimSpace(s)
		vmconfig, err := util.Base64DecodeStripped(s)
		if err != nil {
			continue
		}
		vmess := Vmess{}
		err = json.Unmarshal(vmconfig, &vmess)
		if err != nil {
			log.Println(err)
			continue
		}
		clashVmess := ClashVmess{}
		clashVmess.Name = vmess.PS
		clashVmess.Type = "vmess"
		clashVmess.Server = vmess.Add
		switch vmess.Port.(type) {
		case string:
			clashVmess.Port, _ = vmess.Port.(string)
		case int:
			clashVmess.Port, _ = vmess.Port.(int)
		case float64:
			clashVmess.Port, _ = vmess.Port.(float64)
		default:
			continue
		}
		clashVmess.UUID = vmess.ID
		clashVmess.AlterID = vmess.Aid
		clashVmess.Cipher = vmess.Type
		if "" != vmess.TLS {
			clashVmess.TLS = true
		} else {
			clashVmess.TLS = false
		}
		if "ws" == vmess.Net {
			clashVmess.Network = vmess.Net
			clashVmess.WSPATH = vmess.Path
		}

		vmesss = append(vmesss, clashVmess)
	}
	clash := Clash{}
	r := clash.LoadTemplate("ConnersHua.yaml", vmesss)
	if r == nil {
		c.String(http.StatusBadRequest, "sublink 返回数据格式不对")
		return
	}

	requestURL := c.GetString("request_url")
	if requestURL != "" {
		clashx := fmt.Sprintf("clash://install-config?url=%s", url.PathEscape(requestURL))
		c.Redirect(http.StatusMovedPermanently, clashx)
	} else {
		c.String(http.StatusOK, util.UnicodeEmojiDecode(string(r)))
	}
}

const (
	SSRServer = iota
	SSRPort
	SSRProtocol
	SSRCipher
	SSROBFS
	SSRSuffix
)

func SSR2ClashR(c *gin.Context) {
	decodeBody := c.GetString("decodebody")

	scanner := bufio.NewScanner(strings.NewReader(decodeBody))
	var ssrs []interface{}
	for scanner.Scan() {
		if !strings.HasPrefix(scanner.Text(), "ssr://") {
			continue
		}
		s := scanner.Text()[6:]
		s = strings.TrimSpace(s)
		rawSSRConfig, err := util.Base64DecodeStripped(s)
		if err != nil {
			continue
		}
		params := strings.Split(string(rawSSRConfig), `:`)
		if 6 != len(params) {
			continue
		}
		ssr := ClashRSSR{}
		ssr.Type = "ssr"
		ssr.Server = params[SSRServer]
		ssr.Port = params[SSRPort]
		ssr.Protocol = params[SSRProtocol]
		ssr.Cipher = params[SSRCipher]
		ssr.OBFS = params[SSROBFS]

		// 如果兼容ss协议，就转换为clash的ss配置
		// https://github.com/Dreamacro/clash
		if "origin" == ssr.Protocol && "plain" == ssr.OBFS {
			switch ssr.Cipher {
			case "aes-128-gcm", "aes-192-gcm", "aes-256-gcm",
				"aes-128-cfb", "aes-192-cfb", "aes-256-cfb",
				"aes-128-ctr", "aes-192-ctr", "aes-256-ctr",
				"rc4-md5", "chacha20", "chacha20-ietf", "xchacha20",
				"chacha20-ietf-poly1305", "xchacha20-ietf-poly1305":
				ssr.Type = "ss"
			}
		}

		suffix := strings.Split(params[SSRSuffix], "/?")
		if 2 != len(suffix) {
			continue
		}
		passwordBase64 := suffix[0]
		password, err := util.Base64DecodeStripped(passwordBase64)
		if err != nil {
			continue
		}
		ssr.Password = string(password)

		m, err := url.ParseQuery(suffix[1])
		if err != nil {
			continue
		}

		for k, v := range m {
			de, err := util.Base64DecodeStripped(v[0])
			if err != nil {
				continue
			}
			switch k {
			case "obfsparam":
				ssr.OBFSParam = string(de)
				continue
			case "protoparam":
				ssr.ProtocolParam = string(de)
				continue
			case "remarks":
				ssr.Name = string(de)
				continue
			case "group":
				continue
			}
		}

		ssrs = append(ssrs, ssr)
	}
	clash := Clash{}
	r := clash.LoadTemplate("ConnersHua.yaml", ssrs)
	if r == nil {
		c.String(http.StatusBadRequest, "sublink 返回数据格式不对")
		return
	}
	c.String(http.StatusOK, util.UnicodeEmojiDecode(string(r)))
}
