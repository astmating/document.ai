import requests
import json

openai_api_key = "fk188528-JAPbwe87SKzXwGBroAIdcOLfSC1bAMVU"
url = "https://openai.api2d.net/v1/embeddings"
headers = {
    "Content-Type": "application/json",
    "Authorization": "Bearer "+openai_api_key
}
data = {
    "input": "测试",
    "model": "text-embedding-ada-002"
}
print(url,headers,data)
# 发送请求
response = requests.post(url, headers=headers, data=json.dumps(data))

# 提取嵌入向量
embeddings = json.loads(response.content)
print(embeddings)