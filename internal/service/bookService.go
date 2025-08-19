package service

import (
	"LibraryManagement/internal/api"
	"LibraryManagement/internal/repo/dao"
	"log"
	"strconv"
)

type BookService interface {
	Add(dto *api.BookInfoReq) error
	Delete(ids []string) error
	Update(dto *api.BookUpdateReq) error
	List(dto *api.BookSearchReq) (*api.BookSearchResp, error)
	GetByID(id uint) (*api.BookInfoResp, error)

	// ES搜索功能
	SearchBooks(req *api.BookSearchReq) (*api.BookSearchResp, error)
	SearchByTitle(title string, exact bool) (*[]api.BookInfoResp, error)
	SearchByContent(content string) (*[]api.BookInfoResp, error)

	// 索引管理
	InitializeESIndex() error
	ReindexAllBooks() error
}

type bookServiceImpl struct {
	esService BookESService
}

func NewBookService() BookService {
	return &bookServiceImpl{
		esService: NewBookESService(),
	}
}

func (b *bookServiceImpl) Add(dto *api.BookInfoReq) error {

	// 先保存到数据库
	book, err := dao.ApiDao.BookAddDAO(dto)
	if err != nil {
		return err
	}

	// 同步到ES
	if err := b.esService.IndexBook(book); err != nil {
		log.Printf("同步书籍到ES失败: %v", err)
		// TODO 回滚数据库操作，或者记录到消息队列中重试
	}

	return nil

}

func (b *bookServiceImpl) Delete(ids []string) error {
	// 先从数据库删除
	err := dao.ApiDao.BookDeleteDAO(ids)
	if err != nil {
		return err
	}

	// 从ES中删除
	for _, idStr := range ids {
		if id, parseErr := strconv.ParseUint(idStr, 10, 64); parseErr == nil {
			if esErr := b.esService.DeleteBook(uint(id)); esErr != nil {
				log.Printf("从ES删除书籍失败 (ID: %d): %v", id, esErr)
			}
		}
	}

	return nil

}

// 数据库乐观锁更新
func (b *bookServiceImpl) Update(dto *api.BookUpdateReq) error {

	// 先更新数据库
	book, err := dao.ApiDao.BookUpdateDAO(dto)
	if err != nil {
		return err
	}

	// 同步更新到ES
	if err := b.esService.UpdateBook(book); err != nil {
		log.Printf("同步更新书籍到ES失败: %v", err)
	}

	return nil
}

func (b *bookServiceImpl) List(dto *api.BookSearchReq) (*api.BookSearchResp, error) {

	books, err := dao.ApiDao.BookListDAO(dto)

	return books, err
}

func (b *bookServiceImpl) GetByID(id uint) (*api.BookInfoResp, error) {
	book, err := dao.ApiDao.BookGetByIDDAO(id)
	if err != nil {
		return nil, err
	}

	return &api.BookInfoResp{
		ID:      book.ID,
		Title:   book.Title,
		Count:   book.Count,
		ISBN:    book.ISBN,
		Author:  book.Author,
		Content: book.Content,
		Summary: book.Summary,
	}, nil
}

// SearchBooks ES综合搜索
func (b *bookServiceImpl) SearchBooks(req *api.BookSearchReq) (*api.BookSearchResp, error) {
	return b.esService.SearchBooks(req)
}

// SearchByTitle 标题搜索（精确或模糊）
func (b *bookServiceImpl) SearchByTitle(title string, exact bool) (*[]api.BookInfoResp, error) {
	docs, err := b.esService.SearchByTitle(title, exact)
	if err != nil {
		return nil, err
	}

	books := make([]api.BookInfoResp, 0, len(docs))
	for _, doc := range docs {
		books = append(books, api.BookInfoResp{
			ID:      doc.ID,
			Title:   doc.Title,
			ISBN:    doc.ISBN,
			Author:  doc.Author,
			Summary: doc.Summary,
			// Content在列表中通常不返回，以减少数据传输
		})
	}

	return &books, nil
}

// SearchByContent 内容模糊搜索
func (b *bookServiceImpl) SearchByContent(content string) (*[]api.BookInfoResp, error) {
	docs, err := b.esService.SearchByContent(content)
	if err != nil {
		return nil, err
	}

	books := make([]api.BookInfoResp, 0, len(docs))
	for _, doc := range docs {
		books = append(books, api.BookInfoResp{
			ID:      doc.ID,
			Title:   doc.Title,
			ISBN:    doc.ISBN,
			Author:  doc.Author,
			Summary: doc.Summary,
		})
	}

	return &books, nil
}

// InitializeESIndex 初始化ES索引
func (b *bookServiceImpl) InitializeESIndex() error {
	return b.esService.CreateIndex()
}

// ReindexAllBooks 重新索引所有书籍数据
func (b *bookServiceImpl) ReindexAllBooks() error {
	log.Println("开始重新索引所有书籍...")

	// 删除现有索引
	if err := b.esService.DeleteIndex(); err != nil {
		log.Printf("删除现有索引失败: %v", err)
	}

	// 创建新索引
	if err := b.esService.CreateIndex(); err != nil {
		return err
	}

	// 从数据库获取所有书籍并重新索引
	searchReq := &api.BookSearchReq{
		Page:     1,
		PageSize: 1000, // 批量处理
	}

	page := 1
	for {
		searchReq.Page = page
		listResp, err := dao.ApiDao.BookListDAO(searchReq)
		if err != nil {
			return err
		}

		books := listResp.Books
		if len(books) == 0 {
			break
		}

		// 批量索引到ES
		for _, bookResp := range books {
			// 获取完整的书籍信息
			book, err := dao.ApiDao.BookGetByIDDAO(bookResp.ID)
			if err != nil {
				log.Printf("获取书籍详情失败 (ID: %d): %v", bookResp.ID, err)
				continue
			}

			if err := b.esService.IndexBook(book); err != nil {
				log.Printf("索引书籍失败 (ID: %d): %v", book.ID, err)
			}
		}

		if len(books) < searchReq.PageSize {
			break
		}
		page++
	}

	log.Println("重新索引完成")
	return nil
}
