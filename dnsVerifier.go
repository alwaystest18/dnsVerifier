package main

import (
	"flag"
	"io/ioutil"
	"math/rand"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/Workiva/go-datastructures/queue"
	"github.com/miekg/dns"
	"github.com/remeh/sizedwaitgroup"
)

var vaildResolversList []string

// 将文件内容转为字符列表
func FileContentToList(filePath string) []string {
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return []string{""}
	}
	contentList := strings.Split(string(fileContent), "\n")
	var newList []string
	for _, element := range contentList {
		if element != "" {
			newList = append(newList, element)
		}
	}
	return newList
}

// 去除换行
func RemoveCRLF(resolver string) string {
	r := strings.NewReplacer("\r", "", "\n", "")
	return r.Replace(resolver)
}

// 检查是否为正确的ip格式
func CheckIp(str string) bool {
	regCheckIp := regexp.MustCompile(`\d+\.\d+\.\d+\.\d+(:\d+)?`)
	if regCheckIp.MatchString(str) {
		return true
	} else {
		return false
	}
}

// 生成指定位数的随机值
func RandString(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// 判断字符串在数组元素中是否包含
func In(target string, str_array []string) bool {
	sort.Strings(str_array)
	index := sort.SearchStrings(str_array, target)
	if index < len(str_array) && str_array[index] == target {
		return true
	}
	return false
}

// 把列表内容写入文件
func CreateFileWithArr(stringArr []string, fileName string) bool {
	if len(stringArr) > 0 {
		currFile, err := os.OpenFile(fileName, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			panic("File create fail")
		}
		for _, content := range stringArr {
			currFile.WriteString(content + "\n")
		}
		defer currFile.Close()
		return true
	} else {
		return false
	}
}

// 获取指定dns解析域名的A记录结果
func getDnsTypeA(resolver string, domain string, timeout int64) ([]string, error) {
	var resultTypeA []string
	if !strings.Contains(resolver, ":") {
		resolver = resolver + ":53"
	}
	c := new(dns.Client)
	m := new(dns.Msg)
	c.Timeout = time.Duration(timeout) * time.Millisecond
	m.SetQuestion(domain, dns.TypeA)
	r, _, err := c.Exchange(m, resolver)
	if err != nil {
		return resultTypeA, err
	}

	for _, ans := range r.Answer {
		record, _ := ans.(*dns.A)
		if record != nil && len(record.String()) < 100 { //避免record为nil导致的panic及忽略部分返回内容过大的dns（增加网络波动）
			resultTypeA = append(resultTypeA, record.A.String())
		}
	}
	return resultTypeA, err
}

// 筛选可用dns服务器
func VerifyResolver(resolver string, timeout int64) {
	domain := RandString(20) + ".com."
	result, err := getDnsTypeA(resolver, "public1.114dns.com.", timeout)
	if err != nil {
		return
	}

	//1.正确解析指定域名 2.解析不存在域名无结果 3.解析不存在域名不超时 满足此三个条件判断dns服务器可用
	if In("114.114.114.114", result) {
		noExistResult, err := getDnsTypeA(resolver, domain, timeout)
		if err == nil && len(noExistResult) == 0 {
			vaildResolversList = append(vaildResolversList, resolver)
		}
	}

}

func main() {

	r := flag.String("r", "", "dns resolvers file")
	o := flag.String("o", "vaildResolvers.txt", "output vaild resolvers to file")
	t := flag.Int64("t", 100, "Number of threads")
	rate := flag.Int("rate", 500, "rate limit")
	timeout := flag.Int64("timeout", 300, "DNS maximum connection time(ms)")
	flag.Parse()

	resolversFile := *r
	vaildResolversFile := *o
	threads := *t
	rateLimit := *rate
	dnsTimeout := *timeout

	var wg sizedwaitgroup.SizedWaitGroup = sizedwaitgroup.New(rateLimit)

	que := queue.New(threads)

	resolverList := FileContentToList(resolversFile)
	for _, resolver := range resolverList {
		resolver = RemoveCRLF(resolver)
		if CheckIp(resolver) {
			que.Put(resolver)
		}
	}

	for que.Len() > 0 {
		wg.Add()
		queResolverList, _ := que.Get(1)
		queResolver := queResolverList[0].(string)
		go func() {
			defer wg.Done()
			VerifyResolver(queResolver, dnsTimeout)
		}()
	}
	wg.Wait()
	if len(vaildResolversList) > 0 {
		CreateFileWithArr(vaildResolversList, vaildResolversFile)
	}
}
