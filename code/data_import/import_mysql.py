import pymysql
import pandas as pd
from qdrant_client import QdrantClient
from qdrant_client.http.models import Distance, VectorParams
from qdrant_client.http.models import PointStruct
import os
import tqdm
import openai
from dotenv import load_dotenv

def to_embeddings(items):
    sentence_embeddings = openai.Embedding.create(
        model="text-embedding-ada-002",
        input=items[1]
    )
    return [items[0], items[1], sentence_embeddings["data"][0]["embedding"]]


if __name__ == '__main__':
    client = QdrantClient("127.0.0.1", port=6333)
    collection_name = "data_collection"
    load_dotenv()
    openai.api_key = os.getenv("OPENAI_API_KEY")
    # 创建collection
    client.recreate_collection(
        collection_name=collection_name,
        vectors_config=VectorParams(size=1536, distance=Distance.COSINE),
    )
    # Connect to the database
    connection = pymysql.connect(host='127.0.0.1',
                                 port=3306,
                                 user='root',
                                 password='1111',
                                 db='bt_db_name',
                                 charset='utf8mb4')
    
    # Read data from MySQL table
    data = pd.read_sql_query("SELECT * FROM article", connection)
    count = 0
    for row in tqdm.tqdm(data.iterrows()):
        print(row[1]['title'], row[1]['content'])
        item = to_embeddings([row[1]['title'], row[1]['content']])
        client.upsert(
            collection_name=collection_name,
            wait=True,
            points=[
                PointStruct(id=count, vector=item[2], payload={"title": item[0], "text": item[1]}),
            ],
        )
        count += 1
    
    # Close connections
    connection.close()