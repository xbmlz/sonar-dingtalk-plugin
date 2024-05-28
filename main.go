package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

// dingtalkHandler
func dingtalkHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sonarRsp := make(map[string]interface{})
	accessToken := r.Form.Get("access_token")
	sonarToken := r.Form.Get("sonar_token")
	if accessToken == "" {
		fmt.Fprintf(w, "access_token不能为空")
	}
	if err := json.NewDecoder(r.Body).Decode(&sonarRsp); err != nil {
		r.Body.Close()
		fmt.Fprintf(w, "解析Sonar参数错误:"+err.Error())
		return
	}

	serverUrl := sonarRsp["serverUrl"]
	projectName := sonarRsp["project"].(map[string]interface{})["name"]
	projectKey := sonarRsp["project"].(map[string]interface{})["key"]
	branch := sonarRsp["branch"].(map[string]interface{})["name"]
	// create http client
	httpClient := &http.Client{
		Transport: &http.Transport{
			// 设置代理 HTTPS_PROXY
			Proxy: http.ProxyFromEnvironment,
		},
	}
	// get measures info
	sonarUrl := fmt.Sprintf("%s/api/measures/search?projectKeys=%s&metricKeys=alert_status,bugs,reliability_rating,vulnerabilities,security_rating,code_smells,sqale_rating,duplicated_lines_density,coverage,ncloc,ncloc_language_distribution",
		serverUrl, projectKey)
	req, _ := http.NewRequest("GET", sonarUrl, nil)
	if sonarToken != "" {
		req.SetBasicAuth(sonarToken, "")
	}
	measuresRsp, err := httpClient.Do(req)
	if err != nil {
		fmt.Fprintf(w, "获取measures失败: "+err.Error())
		return
	}
	measuresObj := make(map[string]interface{})
	if err := json.NewDecoder(measuresRsp.Body).Decode(&measuresObj); err != nil {
		measuresRsp.Body.Close()
		fmt.Fprintf(w, "解析Measures失败: "+err.Error())
		return
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
	var statusText string
	if alertStatus == "OK" {
		statusText = "<font color=\"#008000\">正常</font>"
	} else {
		statusText = "<font color=\"#f90202\">异常</font>"
	}
	// 发送钉钉消息
	msgUrl := fmt.Sprintf("https://oapi.dingtalk.com/robot/send?access_token=%s", accessToken)

	messageUrl := fmt.Sprintf("%s/dashboard?id=%s", serverUrl, projectKey)

	actionCard := make(map[string]string)
	title := fmt.Sprintf("### %s[%s]代码扫描报告", projectName, branch)
	actionCard["title"] = title
	actionCard["text"] = fmt.Sprintf("%s <br/>\n- **状态:** %s \n - **Bugs:** %s \n- **漏洞:** %s \n- **异味:** %s \n- **覆盖率:** %s%% \n- **重复率:** %s%%",
		title, statusText, bugs, vulnerabilities, codeSmells, coverage, duplicatedLinesDensity)
	actionCard["btnOrientation"] = "0"
	actionCard["singleURL"] = fmt.Sprintf("dingtalk://dingtalkclient/page/link?url=%s&&pc_slide=false", url.QueryEscape(messageUrl))
	actionCard["singleTitle"] = "查看详情"

	param := make(map[string]interface{})
	param["msgtype"] = "actionCard"
	param["actionCard"] = actionCard

	// send dingtalk message
	paramBytes, _ := json.Marshal(param)
	dingTalkRsp, _ := http.Post(msgUrl, "application/json", bytes.NewBuffer(paramBytes))
	dingTalkObj := make(map[string]interface{})
	json.NewDecoder(dingTalkRsp.Body).Decode(&dingTalkObj)
	if dingTalkObj["errcode"].(float64) != 0 {
		fmt.Fprint(w, "消息推送失败，请检查钉钉机器人配置,错误信息："+dingTalkObj["errmsg"].(string))
		return
	}
	fmt.Fprint(w, "消息推送成功")
}

func main() {
	http.HandleFunc("/dingtalk", dingtalkHandler)
	log.Println("Server started on port(s): 0.0.0.0:9010 (http)")
	log.Fatal(http.ListenAndServe("0.0.0.0:9010", nil))
}
