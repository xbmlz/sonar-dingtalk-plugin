package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

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
	// 项目标识
	projectKey := sonarRsp["project"].(map[string]interface{})["key"]
	// 分支名称
	branchName := sonarRsp["branch"].(map[string]interface{})["name"]
	// sonar prop
	var totalBugs, vulnerabilities, codeSmells, coverage, duplicatedLinesDensity, alertStatus string
	// dingtalk prop
	var sendUrl, text, picUrl, messageUrl string

	// get sonar info
	resp, _ := http.Get(fmt.Sprintf("%s/api/measures/search?projectKeys=%s&metricKeys=alert_status,bugs,reliability_rating,vulnerabilities,security_rating,code_smells,sqale_rating,duplicated_lines_density,coverage,ncloc,ncloc_language_distribution",
		sonarUrl, projectKey))
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

func main() {
	http.HandleFunc("/dingtalk", dingTalkHandler)
	log.Println("Server started on port(s): 9001 (http)")
	log.Fatal(http.ListenAndServe(":9001", nil))
}
