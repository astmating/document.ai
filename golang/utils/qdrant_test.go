package utils

import (
	"encoding/json"
	"log"
	"testing"
)

func TestPutCollection(t *testing.T) {
	collectionName := "data_collection"
	err := PutCollection(collectionName)
	if err != nil {
		t.Errorf("Error putting collection: %v", err)
	}
	log.Println(err)
}
func TestDeleteCollection(t *testing.T) {
	collectionName := "data_collection"
	err := DeleteCollection(collectionName)
	if err != nil {
		t.Errorf("Error putting collection: %v", err)
	}
	log.Println(err)
}
func TestPutPoints(t *testing.T) {

	collectionName := "data_collection"
	points := []map[string]interface{}{
		{
			"id":      1,
			"payload": map[string]interface{}{"title": "测试标题", "text": "测试内容"},
			"vector":  []float64{0, 9, 0.9, 0.9},
		},
	}
	res, err := PutPoints(collectionName, points)
	if err != nil {
		t.Errorf("Error putting points: %v", err)
	}
	log.Println(res, err)
}
func TestSearchPoints(t *testing.T) {

	collectionName := "data_collection"
	params := map[string]interface{}{"exact": false, "hnsw_ef": 128}
	vector := []float64{0, 9, 0.9, 0.9}
	limit := 10
	points, err := SearchPoints(collectionName, params, vector, limit)
	if err != nil {
		t.Errorf("Error searching points: %v", err)
	}
	log.Println(string(points))
}

func TestPutPoints2(t *testing.T) {
	gpt := NewChatGptTool("https://openai.api2d.net", "fk188528-JAPbwe87SKzXwGBroAIdcOLfSC1bAMVU")
	response, err := gpt.GetEmbedding("测试", "text-embedding-ada-002")
	var embeddingResponse EmbeddingResponse
	json.Unmarshal([]byte(response), &embeddingResponse)

	collectionName := "data_collection"
	points := []map[string]interface{}{
		{
			"id":      1,
			"payload": map[string]interface{}{"title": "测试标题", "text": "测试内容"},
			"vector":  embeddingResponse.Data[0].Embedding,
		},
	}
	res, err := PutPoints(collectionName, points)
	if err != nil {
		t.Errorf("Error putting points: %v", err)
	}
	log.Println(res, err)
}
func TestSearchPoints2(t *testing.T) {
	gpt := NewChatGptTool("https://openai.api2d.net", "fk188528-JAPbwe87SKzXwGBroAIdcOLfSC1bAMVU")
	response, err := gpt.GetEmbedding("测试", "text-embedding-ada-002")
	var embeddingResponse EmbeddingResponse
	json.Unmarshal([]byte(response), &embeddingResponse)

	collectionName := "data_collection"
	params := map[string]interface{}{"exact": false, "hnsw_ef": 128}
	vector := embeddingResponse.Data[0].Embedding
	limit := 10
	points, err := SearchPoints(collectionName, params, vector, limit)
	if err != nil {
		t.Errorf("Error searching points: %v", err)
	}
	log.Println(string(points))
}
func TestGetPoints(t *testing.T) {

	collectionName := "data_collection"
	limit := 1
	offset := "16c802db-b52f-434a-b038-7777edd0b5c9"
	points, err := GetPoints(collectionName, uint(limit), offset)
	if err != nil {
		t.Errorf("Error searching points: %v", err)
	}
	log.Println(string(points))
}
