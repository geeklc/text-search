### 说明
使用ollama + weaviate实现文本语义搜索

### 环境准备
✅ollama 0.3.11

✅weaviate:1.25.17

✅Embedding模型   znbang/bge:large-zh-v1.5-f32

✅golang 1.22.7


### 文件说明
/src/docker-weaviate下有weaviate的docker-compose.yaml文件

src/weaviate/collection.go ：创建带有ollama模型的class

src/weaviate/collection1.go: 创建普通的class模型，代码处理text-to-vector，保存时携带向量，查询时携带向量。

src/weaviate/collection2.go： 创建带有ollama模型的class，模型使用nomic-embed-text

src/weaviate/collection-beg.go： 创建带有ollama模型的class，模型使用znbang/bge:large-zh-v1.5-f32