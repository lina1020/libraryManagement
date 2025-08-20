package service

import (
	"LibraryManagement/internal/api"
	"LibraryManagement/internal/es"
	"LibraryManagement/internal/model"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8/esapi"
)

const (
	BooksIndex        = "books"
	DefaultSearchSize = 100
	RequestTimeout    = 5 * time.Second
)

type BookESService interface {
	// 索引管理
	CreateIndex() error
	DeleteIndex() error

	// 文档操作
	IndexBook(book *model.Book) error
	UpdateBook(book *model.Book) error
	DeleteBook(id uint) error
	GetBook(id uint) (*model.ESBookDocument, error)

	// 搜索功能
	SearchBooks(req *api.BookSearchReq) (*api.BookSearchResp, error)
	SearchByTitle(title string, exact bool) ([]model.ESBookDocument, error)
	SearchByContent(content string) ([]model.ESBookDocument, error)
}

type bookESServiceImpl struct{}

func (s *bookESServiceImpl) CreateIndex() error {
	if es.Client == nil {
		log.Println("ES客户端未初始化，跳过索引创建")
		return nil
	}

	indexMapping := `{
		"mappings": {
			"properties": {
				"id": {"type": "long"},
				"title": {
					"type": "text",
					"analyzer": "ik_max_word",
					"search_analyzer": "ik_smart",
					"fields": {
						"keyword": {"type": "keyword"}
					}
				},
				"author": {
					"type": "text",
					"analyzer": "ik_max_word",
					"search_analyzer": "ik_smart",
					"fields": {
						"keyword": {"type": "keyword"}
					}
				},
				"count": {"type": "long"},
				"isbn": {"type": "keyword"},
				"content": {
					"type": "text",
					"analyzer": "ik_max_word",
					"search_analyzer": "ik_smart"
				},
				"summary": {
					"type": "text",
					"analyzer": "ik_max_word",
					"search_analyzer": "ik_smart"
				}
			}
		},
		"settings": {
			"number_of_shards": 1,
			"number_of_replicas": 0
		}
	}`

	req := esapi.IndicesCreateRequest{
		Index: BooksIndex,
		Body:  strings.NewReader(indexMapping),
	}

	// req.Do(...) 需要 context 是为了“控制请求生命周期”
	res, err := req.Do(context.Background(), es.Client)
	if err != nil {
		return fmt.Errorf("创建索引失败: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return fmt.Errorf("解析错误响应失败: %w", err)
		}
		// 如果索引已存在，不视为错误
		if errorType, ok := e["error"].(map[string]interface{})["type"]; ok && errorType == "resource_already_exists_exception" {
			log.Printf("索引 %s 已存在", BooksIndex)
			return nil
		}
		return fmt.Errorf("创建索引失败: %v", e["error"])
	}

	log.Printf("成功创建索引: %s", BooksIndex)
	return nil
}

func (s *bookESServiceImpl) DeleteIndex() error {
	if es.Client == nil {
		log.Println("ES客户端未初始化，跳过删除索引")
		return nil
	}

	req := esapi.IndicesDeleteRequest{Index: []string{BooksIndex}}
	res, err := req.Do(context.Background(), es.Client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("删除索引失败: %s", res.Status())
	}

	return nil

}

func (s *bookESServiceImpl) IndexBook(book *model.Book) error {
	// 1. 参数校验
	if es.Client == nil {
		log.Println("ES客户端未初始化，跳过索引操作")
		return nil
	}

	// 2. 构建文档
	doc := model.ESBookDocument{
		ID:      book.ID,
		Title:   book.Title,
		Count:   book.Count,
		Author:  book.Author,
		ISBN:    book.ISBN,
		Content: book.Content,
		Summary: book.Summary,
	}

	// 3. 序列化文档
	data, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("序列化文档失败: %w", err)
	}

	// 4. 创建请求
	req := esapi.IndexRequest{
		Index:      BooksIndex,
		DocumentID: strconv.FormatUint(uint64(book.ID), 10),
		Body:       bytes.NewReader(data),
		// 索引完这个文档后，立即刷新（refresh）索引，让这个文档马上可以被搜索到。
		Refresh: "true",
	}

	// 5. 执行请求
	res, err := req.Do(context.Background(), es.Client)
	if err != nil {
		return fmt.Errorf("索引文档失败: %w", err)
	}
	defer res.Body.Close()

	// 6. 检查 HTTP 状态
	if res.IsError() {
		return fmt.Errorf("索引文档失败: %s", res.Status())
	}

	log.Printf("成功索引书籍: %s (ID: %d)", book.Title, book.ID)
	return nil
}

func (s *bookESServiceImpl) UpdateBook(book *model.Book) error {
	// 直接重新索引
	return s.IndexBook(book)
}

func (s *bookESServiceImpl) DeleteBook(id uint) error {
	if es.Client == nil {
		return nil
	}

	req := esapi.DeleteRequest{
		Index:      BooksIndex,
		DocumentID: strconv.FormatUint(uint64(id), 10),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), es.Client)
	if err != nil {
		return fmt.Errorf("删除文档失败: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() && res.StatusCode != 404 {
		return fmt.Errorf("删除文档失败: %s", res.Status())
	}

	return nil
}

func (s *bookESServiceImpl) GetBook(id uint) (*model.ESBookDocument, error) {
	if es.Client == nil {
		return nil, fmt.Errorf("ES客户端未初始化")
	}

	req := esapi.GetRequest{
		Index:      BooksIndex,
		DocumentID: strconv.FormatUint(uint64(id), 10),
	}

	res, err := req.Do(context.Background(), es.Client)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		if res.StatusCode == 404 {
			return nil, fmt.Errorf("文档不存在")
		}
		return nil, fmt.Errorf("获取文档失败: %s", res.Status())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	source := result["_source"].(map[string]interface{})
	doc := decodeESDoc(source)

	return &doc, nil
}

// SearchBooks 综合搜索书籍
func (s *bookESServiceImpl) SearchBooks(req *api.BookSearchReq) (*api.BookSearchResp, error) {
	if es.Client == nil {
		return nil, fmt.Errorf("ES客户端未初始化")
	}

	// 设置默认分页参数
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	from := (page - 1) * pageSize

	// 构建查询
	var query map[string]interface{}

	if req.Keyword != "" {
		// 全文搜索
		query = map[string]interface{}{
			"bool": map[string]interface{}{
				"should": []map[string]interface{}{
					{
						"multi_match": map[string]interface{}{
							"query":  req.Keyword,
							"fields": []string{"title^3", "author^2", "content", "summary^1.5"},
							"type":   "best_fields",
						},
					},
				},
				"minimum_should_match": 1,
			},
		}
	} else {
		// 构建布尔查询
		var must []map[string]interface{}

		if req.Title != "" {
			must = append(must, map[string]interface{}{
				"match": map[string]interface{}{
					"title": req.Title,
				},
			})
		}

		if req.Author != "" {
			must = append(must, map[string]interface{}{
				"match": map[string]interface{}{
					"author": req.Author,
				},
			})
		}

		if req.ISBN != "" {
			must = append(must, map[string]interface{}{
				"term": map[string]interface{}{
					"isbn": req.ISBN,
				},
			})
		}

		if req.Content != "" {
			must = append(must, map[string]interface{}{
				"match": map[string]interface{}{
					"content": req.Content,
				},
			})
		}

		if len(must) == 0 {
			query = map[string]interface{}{"match_all": map[string]interface{}{}}
		} else {
			query = map[string]interface{}{
				"bool": map[string]interface{}{
					"must": must,
				},
			}
		}
	}

	// 构建搜索请求
	searchBody := map[string]interface{}{
		"query": query,
		"from":  from,
		"size":  pageSize,
		"sort": []map[string]interface{}{
			{"_score": map[string]string{"order": "desc"}},
			{"id": map[string]string{"order": "desc"}},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(searchBody); err != nil {
		return nil, fmt.Errorf("编码搜索请求失败: %w", err)
	}

	// 执行搜索（带超时）
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := es.Client.Search(
		es.Client.Search.WithContext(ctx),
		es.Client.Search.WithIndex(BooksIndex),
		es.Client.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, fmt.Errorf("搜索失败: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("搜索失败: %s", res.Status())
	}

	// 解析响应
	var result struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source struct {
					ID      float64 `json:"id"`
					Title   string  `json:"title"`
					Count   float64 `json:"count"`
					Author  string  `json:"author"`
					ISBN    string  `json:"isbn"`
					Summary string  `json:"summary"`
				} `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("es响应解码失败: %w", err)
	}

	// 转换为响应结构
	books := make([]api.BookInfoResp, 0, len(result.Hits.Hits))
	for _, hit := range result.Hits.Hits {
		books = append(books, api.BookInfoResp{
			ID:      uint(hit.Source.ID),
			Title:   hit.Source.Title,
			Count:   uint(hit.Source.Count),
			Author:  hit.Source.Author,
			ISBN:    hit.Source.ISBN,
			Summary: hit.Source.Summary,
		})
	}

	totalPages := int((result.Hits.Total.Value + int64(pageSize) - 1) / int64(pageSize))

	return &api.BookSearchResp{
		Books:      books,
		Total:      result.Hits.Total.Value,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// SearchByTitle 根据标题搜索（支持精确和模糊）
func (s *bookESServiceImpl) SearchByTitle(title string, exact bool) ([]model.ESBookDocument, error) {
	if es.Client == nil {
		return []model.ESBookDocument{}, nil
	}

	var query map[string]interface{}
	if exact {
		query = map[string]interface{}{
			"term": map[string]interface{}{
				"title.keyword": title,
			},
		}
	} else {
		query = map[string]interface{}{
			"match": map[string]interface{}{
				"title": title,
			},
		}
	}

	searchBody := map[string]interface{}{
		"query": query,
		"size":  100, // 限制返回数量
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(searchBody); err != nil {
		return nil, err
	}

	res, err := es.Client.Search(
		es.Client.Search.WithContext(context.Background()),
		es.Client.Search.WithIndex(BooksIndex),
		es.Client.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("搜索失败: %s", res.Status())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	hits := result["hits"].(map[string]interface{})["hits"].([]interface{})
	documents := make([]model.ESBookDocument, 0, len(hits))

	for _, hit := range hits {
		source := hit.(map[string]interface{})["_source"].(map[string]interface{})
		documents = append(documents, decodeESDoc(source))
	}

	return documents, nil
}

// SearchByContent 根据内容模糊搜索
func (s *bookESServiceImpl) SearchByContent(content string) ([]model.ESBookDocument, error) {
	if es.Client == nil {
		return []model.ESBookDocument{}, nil
	}

	query := map[string]interface{}{
		"match": map[string]interface{}{
			"content": content,
		},
	}

	searchBody := map[string]interface{}{
		"query": query,
		"size":  100,
		"highlight": map[string]interface{}{
			"fields": map[string]interface{}{
				"content": map[string]interface{}{},
			},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(searchBody); err != nil {
		return nil, err
	}

	res, err := es.Client.Search(
		es.Client.Search.WithContext(context.Background()),
		es.Client.Search.WithIndex(BooksIndex),
		es.Client.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("搜索失败: %s", res.Status())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	hits := result["hits"].(map[string]interface{})["hits"].([]interface{})
	documents := make([]model.ESBookDocument, 0, len(hits))

	for _, hit := range hits {
		source := hit.(map[string]interface{})["_source"].(map[string]interface{})
		doc := decodeESDoc(source)

		documents = append(documents, doc)
	}

	return documents, nil
}

func NewBookESService() BookESService {
	return &bookESServiceImpl{}
}

// ---------- 工具函数 ----------

func decodeESDoc(source map[string]interface{}) model.ESBookDocument {
	doc := model.ESBookDocument{}
	if id, ok := source["id"].(float64); ok {
		doc.ID = uint(id)
	}
	if title, ok := source["title"].(string); ok {
		doc.Title = title
	}
	if author, ok := source["author"].(string); ok {
		doc.Author = author
	}
	if isbn, ok := source["isbn"].(string); ok {
		doc.ISBN = isbn
	}
	if content, ok := source["content"].(string); ok {
		doc.Content = content
	}
	if summary, ok := source["summary"].(string); ok {
		doc.Summary = summary
	}

	if count, ok := source["count"].(float64); ok {
		doc.Count = uint(count)
	}
	return doc
}
