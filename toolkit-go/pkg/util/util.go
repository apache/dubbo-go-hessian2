package util

import (
	"os"
	"strings"
)

// FileExists 判断文件是否存在
func FileExists(p string) bool {
	if _, err := os.Stat(p); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func SnakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	for i := 0; i < len(s); i++ {
		d := s[i]
		if (i > 0 && d >= 'A' && d <= 'Z') &&
			(!(i+1 < len(s) && s[i+1] >= 'A' && s[i+1] <= 'Z') ||
				!(i-1 > 0 && s[i-1] >= 'A' && s[i-1] <= 'Z')) {
			data = append(data, '_')
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}

// CreateDir 创建文件夹，不存在就创建
func CreateDir(p string) error {
	f, err := os.Stat(p)
	if err != nil {
		return os.MkdirAll(p, 0755)
	}
	if !f.IsDir() {
		return os.MkdirAll(p, 0755)
	}
	return nil
}

var specialMap = map[string]string{
	"ACL":   "ACL",
	"API":   "API",
	"ASCII": "ASCII",
	"CPU":   "CPU",
	"CSS":   "CSS",
	"DNS":   "DNS",
	"EOF":   "EOF",
	"GUID":  "GUID",
	"HTML":  "HTML",
	"HTTP":  "HTTP",
	"HTTPS": "HTTPS",
	"ID":    "ID",
	"IP":    "IP",
	"JSON":  "JSON",
	"JSONP": "JSONP",
	"LHS":   "LHS",
	"QPS":   "QPS",
	"RAM":   "RAM",
	"RHS":   "RHS",
	"RPC":   "RPC",
	"SLA":   "SLA",
	"SMTP":  "SMTP",
	"SQL":   "SQL",
	"SSH":   "SSH",
	"TCP":   "TCP",
	"TLS":   "TLS",
	"TTL":   "TTL",
	"UDP":   "UDP",
	"UI":    "UI",
	"UID":   "UID",
	"UUID":  "UUID",
	"URI":   "URI",
	"URL":   "URL",
	"UTF8":  "UTF8",
	"VM":    "VM",
	"XML":   "XML",
	"XMPP":  "XMPP",
	"XSRF":  "XSRF",
	"XSS":   "XSS",

	"ECS":  "ECS",
	"VPC":  "VPC",
	"OSS":  "OSS",
	"SLB":  "SLB",
	"RDS":  "RDS",
	"CDN":  "CDN",
	"HPC":  "HPC",
	"NAT":  "NAT",
	"ICMP": "ICMP",
	"GW":   "GW",
	"EDAS": "EDAS",
	"DRDS": "DRDS",
	"ARMS": "ARMS",
	"MQ":   "MQ",
	"CSB":  "CSB",
	"TS":   "TS",
}

func UpFirst(s string) string {
	if len(s) == 0 {
		return s
	}

	if v, ok := specialMap[strings.ToUpper(s)]; ok {
		return v
	}

	return strings.ToUpper(string(s[0])) + s[1:]
}

// IsDir 是否是文件
func IsDir(p string) bool {
	f, err := os.Stat(p)
	if err != nil {
		return false
	}
	return f.IsDir()
}
