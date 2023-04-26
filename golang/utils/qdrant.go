package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	QdrantBase = "127.0.0.1"
	QdrantPort = "6333"
)

//创建集合
func PutCollection(collectionName string) error {
	url := fmt.Sprintf("http://%s:%s/collections/%s", QdrantBase, QdrantPort, collectionName)
	requestBody, err := json.Marshal(map[string]interface{}{
		"name": collectionName,
		"vectors": map[string]interface{}{
			"size":     1536,
			"distance": "Cosine",
		},
	})
	if err != nil {
		return err
	}

	request, err := http.NewRequest("PUT", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return nil
}

//删除集合
func DeleteCollection(collectionName string) error {
	url := fmt.Sprintf("http://%s:%s/collections/%s", QdrantBase, QdrantPort, collectionName)

	request, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	return nil
}

//列出所有集合
func GetCollections() ([]byte, error) {
	url := fmt.Sprintf("http://%s:%s/collections", QdrantBase, QdrantPort)
	resp, err := http.Get(url)
	if err != nil {
		return []byte(""), err
	}
	defer resp.Body.Close()

	result, _ := ioutil.ReadAll(resp.Body)
	return result, nil
}

//查询集合信息
func GetCollection(collectionName string) ([]byte, error) {
	url := fmt.Sprintf("http://%s:%s/collections/%s", QdrantBase, QdrantPort, collectionName)
	resp, err := http.Get(url)
	if err != nil {
		return []byte(""), err
	}
	defer resp.Body.Close()

	result, _ := ioutil.ReadAll(resp.Body)
	return result, nil
}

//查询数据列表
func GetPoints(collectionName string, limit uint, offset interface{}) ([]byte, error) {
	// 构造请求体
	requestBody := map[string]interface{}{
		"limit":        limit,
		"offset":       offset,
		"with_payload": true,
		"with_vector":  false,
	}
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	// 构造请求
	url := fmt.Sprintf("http://%s:%s/collections/%s/points/scroll", QdrantBase, QdrantPort, collectionName)
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := http.DefaultClient
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// 处理响应
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return responseBody, nil
}

//增加向量数据
func PutPoints(collectionName string, points []map[string]interface{}) (string, error) {
	url := fmt.Sprintf("http://%s:%s/collections/%s/points", QdrantBase, QdrantPort, collectionName)

	// 构造请求体
	requestBody := map[string]interface{}{
		"points": points,
	}
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	// 发送请求
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(requestBodyBytes))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	res, _ := ioutil.ReadAll(resp.Body)
	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return string(res), fmt.Errorf("failed to PUT points to collection %s, status code: %d"+string(requestBodyBytes), collectionName, resp.StatusCode)
	}

	return string(res), nil
}

//增加向量数据
func DeletePoints(collectionName string, points interface{}) (string, error) {
	url := fmt.Sprintf("http://%s:%s/collections/%s/points/delete", QdrantBase, QdrantPort, collectionName)

	// 构造请求体
	requestBody := map[string]interface{}{
		"points": points,
	}
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	// 发送请求
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(requestBodyBytes))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	res, _ := ioutil.ReadAll(resp.Body)
	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return string(res), errors.New(string(res))
	}

	return string(res), nil
}

//搜索向量数据
func SearchPoints(collectionName string, params map[string]interface{}, vector []float64, limit int) ([]byte, error) {
	// 构造请求体
	requestBody := map[string]interface{}{
		"params":       params,
		"vector":       vector,
		"limit":        limit,
		"with_payload": true,
	}
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	// 构造请求
	url := fmt.Sprintf("http://%s:%s/collections/%s/points/search", QdrantBase, QdrantPort, collectionName)
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := http.DefaultClient
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// 处理响应
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return responseBody, nil
}
