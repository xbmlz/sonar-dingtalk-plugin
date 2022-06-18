package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var httpClient = &http.Client{}

// indexHandler
func indexHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sonarRsp := make(map[string]interface{})
	dingTalkToken := r.Form.Get("dingtalk_token")
	if dingTalkToken == "" {
		fmt.Fprintf(w, "access_token 不能为空")
		return
	}
	sonarToken := r.Form.Get("sonar_token")
	if err := json.NewDecoder(r.Body).Decode(&sonarRsp); err != nil {
		r.Body.Close()
		log.Fatal(err)
		fmt.Fprintf(w, "解析Sonar参数错误")
		return
	}

	serverUrl := sonarRsp["serverUrl"]
	projectName := sonarRsp["project"].(map[string]interface{})["name"]
	projectKey := sonarRsp["project"].(map[string]interface{})["key"]
	branch := sonarRsp["branch"].(map[string]interface{})["name"]

	// get measures info
	url := fmt.Sprintf("%s/api/measures/search?projectKeys=%s&metricKeys=alert_status,bugs,reliability_rating,vulnerabilities,security_rating,code_smells,sqale_rating,duplicated_lines_density,coverage,ncloc,ncloc_language_distribution",
		serverUrl, projectKey)
	req, _ := http.NewRequest("GET", url, nil)
	if sonarToken != "" {
		req.SetBasicAuth(sonarToken, "")
	}
	measuresRsp, err := httpClient.Do(req)
	if err != nil {
		fmt.Fprintf(w, "获取measures失败")
		return
	}
	measuresObj := make(map[string]interface{})
	if err := json.NewDecoder(measuresRsp.Body).Decode(&measuresObj); err != nil {
		measuresRsp.Body.Close()
		log.Fatal(err)
		fmt.Fprintf(w, "解析Measures失败")
	}

	measures := measuresObj["measures"].([]interface{})
	alertStatus := (measures[0].(map[string]interface{}))["value"].(string)
	bugs := (measures[1].(map[string]interface{}))["value"].(string)
	codeSmells := (measures[2].(map[string]interface{}))["value"].(string)
	coverage := (measures[3].(map[string]interface{}))["value"].(string)
	duplicatedLinesDensity := (measures[4].(map[string]interface{}))["value"].(string)
	// ncloc := (measures[5].(map[string]interface{}))["value"].(string)
	// nclocLanguageDistribution := (measures[6].(map[string]interface{}))["value"].(string)
	// neliabilityRating := (measures[7].(map[string]interface{}))["value"].(string)
	// securityRating := (measures[8].(map[string]interface{}))["value"].(string)
	// sqaleRating := (measures[9].(map[string]interface{}))["value"].(string)
	vulnerabilities := (measures[10].(map[string]interface{}))["value"].(string)

	// 成功失败标志
	var picUrl string
	if alertStatus == "OK" {
		picUrl = "http://s1.ax1x.com/2020/10/29/BGMeTe.png"
	} else {
		picUrl = "http://s1.ax1x.com/2020/10/29/BGMZwD.png"
	}
	// 发送钉钉消息
	msgUrl := fmt.Sprintf("https://oapi.dingtalk.com/robot/send?access_token=%s", dingTalkToken)

	messageUrl := fmt.Sprintf("%s/dashboard?id=%s", serverUrl, projectKey)

	link := make(map[string]string)
	link["title"] = fmt.Sprintf("%s[%s]代码扫描报告", projectName, branch)
	link["text"] = fmt.Sprintf("Bugs: %s | 漏洞: %s | 异味: %s\r覆盖率: %s%%\r重复率: %s%%",
		bugs, vulnerabilities, codeSmells, coverage, duplicatedLinesDensity)
	link["messageUrl"] = messageUrl
	link["picUrl"] = picUrl

	param := make(map[string]interface{})
	param["msgtype"] = "link"
	param["link"] = link

	// send dingtalk message
	paramBytes, _ := json.Marshal(param)
	dingTalkRsp, _ := http.Post(msgUrl, "application/json", bytes.NewBuffer(paramBytes))
	dingTalkObj := make(map[string]interface{})
	json.NewDecoder(dingTalkRsp.Body).Decode(&dingTalkObj)
	if dingTalkObj["errcode"].(float64) != 0 {
		fmt.Fprint(w, "消息推送失败，请检查钉钉机器人配置")
		return
	}
	fmt.Fprint(w, "消息推送成功")
}

func main() {
	http.HandleFunc("/", indexHandler)
	log.Println("Server started on port(s): 0.0.0.0:9010 (http)")
	log.Fatal(http.ListenAndServe("0.0.0.0:9010", nil))
}
