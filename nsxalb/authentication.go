package nsxalb

// func (a *Agent) Login(cred map[string]string) error {
// 	target_url := a.BaseUrl + "/login"
// 	credJson, _ := json.Marshal(cred)
// 	req, _ := http.NewRequest("POST", target_url, bytes.NewBuffer(credJson))
// 	req.Header.Set("Content-Type", "application/json")
// 	res, err := a.httpClient.Do(req)
// 	if err != nil {
// 		return err
// 	}
// 	defer res.Body.Close()
// 	if res.StatusCode != 200 {
// 		if a.Debug {
// 			log.Printf("StatusCode=%d\n", res.StatusCode)
// 			log.Println(req)
// 			log.Println(res)
// 			data := readResponseBody(res)
// 			log.Println(data)
// 		}
// 		return fmt.Errorf("authentication failed")
// 	}
// 	url, _ := url.Parse(target_url)
// 	cookies := a.httpClient.Jar.Cookies(url)
// 	for i := 0; i < len(cookies); i++ {
// 		if cookies[i].Name == "csrftoken" {
// 			a.Token = cookies[i].Value
// 		}
// 	}
// 	data := readResponseBody(res)
// 	if a.Debug {
// 		log.SetOutput(os.Stderr)
// 		log.Println("login successful")
// 		log.Println(res.Header)
// 		log.Println(data)
// 	}
// 	versionData := data.(map[string]interface{})["version"]
// 	version := versionData.(map[string]interface{})["Version"]
// 	a.Version = version.(string)

// 	return nil
// }

// func (a *Agent) Logout() {
// 	target_url := a.BaseUrl + "/logout"
// 	req, _ := http.NewRequest("POST", target_url, nil)
// 	req.Header.Set("X-CSRFToken", a.Token)
// 	req.Header.Set("Referer", a.BaseUrl)
// 	res, err := a.httpClient.Do(req)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	defer res.Body.Close()
// 	if res.StatusCode != 200 {
// 		log.Printf("StatusCode=%d\n", res.StatusCode)
// 		return
// 	}
// 	if a.Debug {
// 		log.Printf("logout successful\n")
// 		log.Println(res.Header)
// 	}
// }

//func (a *Agent) Cluster() {
// target_url := a.BaseUrl + "/api/pool"
// req, _ := http.NewRequest("GET", target_url, nil)
// res, err := a.httpClient.Do(req)
// if err != nil {
// 	log.Println(err)
// }
// defer res.Body.Close()
// res_body, err := io.ReadAll(res.Body)
// if err != nil {
// 	log.Println(err)
// 	return
// }
// fmt.Printf("StatusCode: %d\n", res.StatusCode)
// var data interface{}
// if len(res_body) > 0 {
// 	err = json.Unmarshal(res_body, &data)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	j, _ := json.MarshalIndent(data, "", "  ")
// 	fmt.Println(string(j))
// } else {
// 	fmt.Println("no response body")
// }
//}
