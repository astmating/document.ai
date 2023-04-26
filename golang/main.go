package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
	"github.com/tidwall/gjson"
	"io"
	"io/ioutil"
	"knowledge/utils"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

func main() {
	err := godotenv.Load(".env")
	// 读取环境变量
	utils.OPENAI_API_KEY = os.Getenv("OPENAI_KEY")
	utils.OPENAI_API_BASE = os.Getenv("OPENAI_API_BASE")
	utils.QdrantBase = os.Getenv("QDRANT_BASE")
	utils.QdrantPort = os.Getenv("QDRANT_PORT")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// 创建 Gin 引擎
	router := gin.Default()
	//启用跨域中间件
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS,DELETE,PUT")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Origin")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})
	router.LoadHTMLGlob("templates/*")
	router.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "text/html;charset=utf-8")
		// 关闭输出缓冲，使得每次写入的数据能够立即发送给客户端
		f, ok := c.Writer.(http.Flusher)
		if !ok {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		// 向客户端输出数据
		html := []rune(`Go-Knowledge AI，自建私有数据知识库 · 与知识库AI聊天`)
		for _, str := range html {
			fmt.Fprintf(c.Writer, string(str))
			f.Flush()
			time.Sleep(100 * time.Millisecond)
		}
	})
	//首页
	router.GET("/:collectName", func(c *gin.Context) {
		collectName := c.Param("collectName")
		c.HTML(http.StatusOK, "index.html", gin.H{
			"collectName": collectName,
		})
	})
	//后台界面
	router.GET("/:collectName/qdrant", func(c *gin.Context) {
		collectName := c.Param("collectName")
		c.HTML(http.StatusOK, "qdrant.html", gin.H{
			"collectName": collectName,
		})
	})
	//集合列表
	router.GET("/collects", func(c *gin.Context) {
		list, err := utils.GetCollections()
		if err != nil {
			c.Writer.Write([]byte(err.Error()))
			return
		}
		c.Writer.Write(list)
	})
	//向量列表
	router.GET("/:collectName/points", func(c *gin.Context) {
		collectName := c.Param("collectName")
		list, err := utils.GetPoints(collectName, 1000, nil)
		if err != nil {
			c.Writer.Write([]byte(err.Error()))
			return
		}
		c.Writer.Write(list)
	})
	//删除向量
	router.GET("/:collectName/delPoints", func(c *gin.Context) {
		collectName := c.Param("collectName")
		id := c.Query("id")
		i, err := strconv.Atoi(id)
		var points interface{}
		if err != nil {
			points = []string{id}
		} else {
			points = []int{i}
		}
		list, err := utils.DeletePoints(collectName, points)
		if err != nil {
			c.Writer.Write([]byte(err.Error()))
			return
		}
		c.Writer.Write([]byte(list))
	})
	//搜索
	router.GET("/:collectName/search", func(c *gin.Context) {
		keywords, _ := url.QueryUnescape(c.Query("keywords"))
		collectName := c.Param("collectName")
		gpt := utils.NewChatGptTool()
		message, err := MakePrompt(collectName, keywords)
		if err != nil {
			return
		}
		res, err := gpt.ChatGPT3Dot5Turbo(message)
		log.Println(message, res, err)
		c.Writer.Write([]byte(res))
	})
	//搜索流式响应
	router.GET("/:collectName/searchStream", func(c *gin.Context) {
		keywords, _ := url.QueryUnescape(c.Query("keywords"))
		collectName := c.Param("collectName")
		gpt := utils.NewChatGptTool()
		message, err := MakePrompt(collectName, keywords)
		if err != nil {
			c.Writer.Write([]byte(err.Error()))
		}
		log.Printf("请求openai：%+v\n", message)
		stream, _ := gpt.ChatGPT3Dot5TurboStream(message)
		// 将响应头中的Content-Type设置为text/plain，表示响应内容为文本
		c.Header("Content-Type", "text/html;charset=utf-8;")

		// 关闭输出缓冲，使得每次写入的数据能够立即发送给客户端
		f, ok := c.Writer.(http.Flusher)
		if !ok {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		for {
			response, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				log.Println("\nStream finished")
				break
			} else if err != nil {
				log.Printf("\nStream error: %v\n", err)
				break
			} else {
				data := response.Choices[0].Delta.Content
				log.Println(data)
				c.Writer.Write([]byte(data))
				f.Flush()
			}
		}
		stream.Close()
	})
	//训练
	router.POST("/:collectName/training", func(c *gin.Context) {
		collectName := c.Param("collectName")
		//判断集合是否存在
		collectInfo, err := utils.GetCollection(collectName)
		collectInfoStatus := gjson.Get(string(collectInfo), "status").String()
		if collectInfoStatus != "ok" {
			utils.PutCollection(collectName)
		}
		//判断ID是否传递，以及是否为数值型或uuid
		id := c.PostForm("id")
		var pointId interface{}
		if id == "" {
			pointId = uuid.NewV4().String()
		} else {
			i, err := strconv.Atoi(id)
			if err != nil {
				pointId = id
			} else {
				pointId = i
			}
		}
		//向量化数据
		content := c.PostForm("content")
		res, err := Train(pointId, collectName, content)
		log.Println(err)
		c.Writer.Write([]byte(res))
	})
	//上传doc
	router.POST("/:collectName/uploadDoc", func(c *gin.Context) {
		collectName := c.Param("collectName")
		f, err := c.FormFile("file")
		if err != nil {
			c.JSON(200, gin.H{
				"code": 400,
				"msg":  "上传失败!" + err.Error(),
			})
			return
		} else {

			fileExt := strings.ToLower(path.Ext(f.Filename))
			if fileExt != ".docx" && fileExt != ".txt" && fileExt != ".xlsx" {
				c.JSON(200, gin.H{
					"code": 400,
					"msg":  "上传失败!只允许txt或docx或xlsx文件",
				})
				return
			}
			fileName := collectName + f.Filename
			c.SaveUploadedFile(f, fileName)
			text := ""
			if fileExt == ".txt" {
				// 打开txt文件
				file, _ := os.Open(fileName)
				// 一次性读取整个txt文件的内容
				txt, _ := ioutil.ReadAll(file)
				text = string(txt)
				file.Close()
			} else if fileExt == ".docx" {
				text, err = utils.ReadDocxAll(fileName)
			} else {
				text, err = utils.ReadExcelAll(fileName)
			}
			removeErr := os.Remove(fileName)
			if removeErr != nil {
				log.Println("Remove error:", fileName, removeErr)
			}
			if err != nil {
				c.JSON(200, gin.H{
					"code": 400,
					"msg":  err.Error(),
				})
				return
			}
			chunks := SplitTextByLength(text, 300)
			for _, chunk := range chunks {
				pointId := uuid.NewV4().String()
				Train(pointId, collectName, chunk)
			}
			os.Remove(fileName)
			c.JSON(200, gin.H{
				"code": 200,
			})
		}
	})
	// 启动服务器
	if err := router.Run(":8083"); err != nil {
		panic(err)
	}
}

//搜索功能
func Search(collectionName, keywords, searchType string) (string, error) {
	gpt := utils.NewChatGptTool()
	message, err := MakePrompt(collectionName, keywords)
	if err != nil {
		return "", err
	}
	res, err := gpt.ChatGPT3Dot5Turbo(message)
	log.Println(message, res, err)
	return res, err
}
func MakePrompt(collectionName, keywords string) ([]utils.Gpt3Dot5Message, error) {
	gpt := utils.NewChatGptTool()
	response, err := gpt.GetEmbedding(keywords, "text-embedding-ada-002")
	if err != nil {
		return nil, err
	}
	var embeddingResponse utils.EmbeddingResponse
	json.Unmarshal([]byte(response), &embeddingResponse)

	params := map[string]interface{}{"exact": false, "hnsw_ef": 128}
	vector := embeddingResponse.Data[0].Embedding
	limit := 8
	points, _ := utils.SearchPoints(collectionName, params, vector, limit)
	result := gjson.Get(string(points), "result").Array()
	message := make([]utils.Gpt3Dot5Message, 0)
	content := ""
	line := ""
	for key, row := range result {
		key++
		line = fmt.Sprintf("%d. %s\n", key, row.Get("payload.text").String())
		arr := []rune(line)
		if len(arr) > 300 {
			line = string(arr[:300])
		}
		//message = append(message, utils.Gpt3Dot5Message{
		//	Role:    "assistant",
		//	Content: line,
		//})
		content += line
	}
	message = append(message, utils.Gpt3Dot5Message{
		Role:    "system",
		Content: "你现在扮演知识库AI机器人。请严格根据提供的参考信息总结归纳后回答问题。对于与参考信息无关的问题或者不理解的问题等，你应拒绝并告知用户“未查询到相关信息，请提供详细的问题信息。”\n避免引用任何当前或过去的政治人物或事件，以及可能引起争议或分裂的历史人物或事件。",
	})
	message = append(message, utils.Gpt3Dot5Message{
		Role:    "user",
		Content: fmt.Sprintf("问题是：\"%s\"\n参考信息是:\"%s\"\n", keywords, content),
	})
	return message, nil
}

//训练功能
func Train(pointId interface{}, collectName, content string) (string, error) {
	gpt := utils.NewChatGptTool()

	response, err := gpt.GetEmbedding(content, "text-embedding-ada-002")
	if err != nil {
		return "", err
	}

	var embeddingResponse utils.EmbeddingResponse
	json.Unmarshal([]byte(response), &embeddingResponse)

	points := []map[string]interface{}{
		{
			"id":      pointId,
			"payload": map[string]interface{}{"text": content},
			"vector":  embeddingResponse.Data[0].Embedding,
		},
	}
	res, err := utils.PutPoints(collectName, points)
	return res, err
}

//长文本分块
func SplitTextByLength(text string, length int) []string {
	var blocks []string
	runes := []rune(text)
	for i := 0; i < len(runes); i += length {
		j := i + length
		if j > len(runes) {
			j = len(runes)
		}
		blocks = append(blocks, string(runes[i:j]))
	}
	return blocks
}
