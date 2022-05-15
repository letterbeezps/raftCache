package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

func JoinRaft(joinAddr, raftAddr, nodeID string) error {
	body := map[string]string{"addr": raftAddr, "id": nodeID}
	byteInfo, err := json.Marshal(body)
	fmt.Println(body)
	if err != nil {
		zap.S().Error(err)
		return err
	}

	joinURL := fmt.Sprintf("http://%s/v1/cache/join", joinAddr)

	zap.S().Infof("send Join request to %s with URL %s", joinAddr, joinURL)

	resp, err := http.Post(
		joinURL,
		"application-type/json",
		bytes.NewReader(byteInfo),
	)
	if err != nil {
		zap.S().Error(err)
		return err
	}
	defer resp.Body.Close()

	return nil
}
