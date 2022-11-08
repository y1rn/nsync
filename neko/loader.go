package neko

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
)

func LoadRule(url, token string) (*RuleResp, error) {
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("token", token)
	req.Header.Set("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("fail to get rule with status code [%v]", resp)
	}
	if err != nil {
		return nil, err
	}

	var rule RuleResp
	respDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(resp.Body).Decode(&rule)

	defer resp.Body.Close()

	if err != nil {
		return nil, fmt.Errorf("fail to get rule with body [%s]", string(respDump))
	}

	return &rule, nil
}
