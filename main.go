package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
)

var addr = flag.String("addr", "0.0.0.0:9001", "输入监听地址")
var token = flag.String("token", "", "输入sonarqube token")
var httpClient = &http.Client{}

func getMeasures(sonarUrl, projectKey interface{}) (*http.Response, error) {
	url := fmt.Sprintf("%s/api/measures/search?projectKeys=%s&metricKeys=alert_status,bugs,reliability_rating,vulnerabilities,security_rating,code_smells,sqale_rating,duplicated_lines_density,coverage,ncloc,ncloc_language_distribution", sonarUrl, projectKey)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(*token, "")
	resp, err := httpClient.Do(req)
	return resp, err
}

// dingTalkHandler
func dingTalkHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sonarRsp := make(map[string]interface{})
	accessToken := r.Form.Get("access_token")
	if err := json.NewDecoder(r.Body).Decode(&sonarRsp); err != nil {
		r.Body.Close()
		log.Fatal(err)
	}
	// sonar地址
	sonarUrl := sonarRsp["serverUrl"]
	// 项目名称
	projectName := sonarRsp["project"].(map[string]interface{})["name"]
	// 项目key
	projectKey := sonarRsp["project"].(map[string]interface{})["key"]
	// 分支名称
	branchName := sonarRsp["branch"].(map[string]interface{})["name"]
	// sonar prop
	var totalBugs, vulnerabilities, codeSmells, coverage, duplicatedLinesDensity, alertStatus string
	// dingtalk prop
	var sendUrl, text, picUrl, messageUrl string

	// get sonar info
	resp, err := getMeasures(sonarUrl, projectKey)
	if err != nil {
		fmt.Printf("request measures error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "request measures error: %v", err)
		return
	}
	measuresObj := make(map[string]interface{})
	if err := json.NewDecoder(resp.Body).Decode(&measuresObj); err == nil {
		measures := measuresObj["measures"].([]interface{})
		log.Println(len(measures))
		for i := 0; i < len(measures); i++ {
			metric := measures[i].(map[string]interface{})
			if metric["metric"] == "bugs" {
				totalBugs = metric["value"].(string)
			} else if metric["metric"] == "vulnerabilities" {
				vulnerabilities = metric["value"].(string)
			} else if metric["metric"] == "code_smells" {
				codeSmells = metric["value"].(string)
			} else if metric["metric"] == "coverage" {
				coverage = metric["value"].(string)
			} else if metric["metric"] == "duplicated_lines_density" {
				duplicatedLinesDensity = metric["value"].(string)
			} else if metric["metric"] == "alert_status" {
				alertStatus = metric["value"].(string)
			}
		}
		// 成功失败标志
		if "ERROR" == alertStatus {
			picUrl = "http://s1.ax1x.com/2020/10/29/BGMZwD.png"
		} else {
			picUrl = "http://s1.ax1x.com/2020/10/29/BGMeTe.png"
		}
		// 钉钉消息
		sendUrl = fmt.Sprintf("https://oapi.dingtalk.com/robot/send?access_token=%s", accessToken)

		messageUrl = fmt.Sprintf("%s/dashboard?id=%s", sonarUrl, projectKey)

		text = fmt.Sprintf("%s[%s]代码扫描结果：BUG数：%s个，漏洞数：%s个，异味数：%s个，覆盖率：%s%%，重复率：%s%%",
			projectName, branchName, totalBugs, vulnerabilities, codeSmells, coverage, duplicatedLinesDensity)

		link := make(map[string]string)
		link["title"] = "代码质量报告"
		link["text"] = text
		link["picUrl"] = picUrl
		link["messageUrl"] = messageUrl

		param := make(map[string]interface{})
		param["msgtype"] = "link"
		param["link"] = link

		// send message
		paramBytes, _ := json.Marshal(param)
		response, _ := http.Post(sendUrl, "application/json", bytes.NewBuffer(paramBytes))
		fmt.Fprint(w, response)
	}
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func main() {
	flag.Parse()
	if !isFlagPassed("token") {
		fmt.Println("token参数是必须的")
		flag.Usage()
		return
	}
	http.HandleFunc("/dingtalk", dingTalkHandler)
	log.Printf("Server started on port(s): %s (http)\n", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
