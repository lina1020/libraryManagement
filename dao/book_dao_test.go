package dao

import (
	"LibraryManagement/dto"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

/*
DAO 层测试 (dao_test.go)

数据库操作测试：测试增删改查操作
约束测试：测试数据库约束（唯一性、长度等）
乐观锁测试：测试版本号控制的乐观锁机制
并发测试：测试数据库操作的并发安全性
事务测试：测试事务的提交和回滚
性能测试：测试数据库操作的性能
*/

func TestBookAddDAO(t *testing.T) {
	gormDB, mock, err := getDBMock()
	assert.NoError(t, err)

	// 准备预期的SQL操作
	mock.ExpectBegin() // 预期事务的开始
	mock.ExpectExec("INSERT INTO `books` (`created_at`,`updated_at`,`deleted_at`,`title`,`count`,`isbn`,`version`) VALUES (?,?,?,?,?,?,?)").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "1", 2, "1548", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// 调用被测试的函数
	bookInfoDTO := &dto.BookInfoDTO{
		Title: "1",
		Count: 2,
		ISBN:  "1548",
	}
	err = BookAddDAO(gormDB, bookInfoDTO)
	assert.NoError(t, err)

	// 验证所有预期的操作都已执行
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestBookDeleteDAO(t *testing.T) {
	gormDB, mock, err := getDBMock()
	assert.NoError(t, err)

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `books` SET `deleted_at`=? WHERE `books`.`id` IN (?,?) AND `books`.`deleted_at` IS NULL").
		WithArgs(sqlmock.AnyArg(), 1, 2).
		WillReturnResult(sqlmock.NewResult(0, 2))
	mock.ExpectCommit()

	err = BookDeleteDAO(gormDB, []string{"1", "2"})
	assert.NoError(t, err)

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `books` SET `deleted_at`=? WHERE `books`.`id` IN (?,?) AND `books`.`deleted_at` IS NULL").
		WithArgs(sqlmock.AnyArg(), 1, 2).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err = BookDeleteDAO(gormDB, []string{"1", "2"})
	assert.NoError(t, err)

	// 验证所有预期的操作都已执行
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestBookListDAO(t *testing.T) {
	gormDB, mock, err := getDBMock()
	assert.NoError(t, err)

	// 准备预期的SQL操作
	mock.ExpectQuery("SELECT `books`.`id`,`books`.`title`,`books`.`count`,`books`.`isbn` FROM `books` WHERE `books`.`deleted_at` IS NULL").
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "count", "isbn"}).
			AddRow(1, "1", 2, "1548").
			AddRow(2, "12345", 2, "测试"))

	_, err = BookListDAO(gormDB, &dto.BookSearchDTO{
		Title: "",
		ISBN:  "",
	})
	assert.NoError(t, err)

	// 准备预期的SQL操作
	mock.ExpectQuery("SELECT `books`.`id`,`books`.`title`,`books`.`count`,`books`.`isbn` FROM `books` WHERE title LIKE ? AND `books`.`deleted_at` IS NULL").
		WithArgs("%2%").
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "count", "isbn"}))

	_, err = BookListDAO(gormDB, &dto.BookSearchDTO{
		Title: "2",
		ISBN:  "",
	})
	assert.NoError(t, err)

	// 准备预期的SQL操作
	mock.ExpectQuery("SELECT `books`.`id`,`books`.`title`,`books`.`count`,`books`.`isbn` FROM `books` WHERE title LIKE ? AND isbn = ? AND `books`.`deleted_at` IS NULL").
		WithArgs("%2%", "测试isbn").
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "count", "isbn"}).
			AddRow(2, "12345", 2, "测试isbn"))

	_, err = BookListDAO(gormDB, &dto.BookSearchDTO{
		Title: "2",
		ISBN:  "测试isbn",
	})
	assert.NoError(t, err)

	// 验证所有预期的操作都已执行
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

/*
重点测试以下场景：

✅ 成功更新（乐观锁通过）
❌ 乐观锁失败（数据被并发修改）
❌ 记录不存在
*/
func TestBookUpdateDAO(t *testing.T) {
	t.Run("更新成功：乐观锁匹配，更新成功", func(t *testing.T) {
		gormDB, mock, err := getDBMock()
		assert.NoError(t, err)

		dto := &dto.BookUpdateDTO{
			ID: 1,
			BookInfoDTO: dto.BookInfoDTO{
				Title: "新标题",
				Count: 100,
				ISBN:  "1234567890",
			},
		}

		// 模拟查询：先 SELECT 获取当前 book（含 version=5）
		rows := sqlmock.NewRows([]string{"id", "title", "count", "isbn", "version"}).
			AddRow(1, "旧标题", 50, "0987654321", 5)

		// 不需要写完整的 SQL，sqlmock 只要匹配开头就行。比如你写 "SELECT"，也能匹配所有 SELECT 语句。
		mock.ExpectQuery("SELECT * FROM `books` WHERE id = ? AND `books`.`deleted_at` IS NULL ORDER BY `books`.`id` LIMIT ?").
			WithArgs(1, 1).
			// 表示当这条查询执行时，返回我提前准备好的数据 rows。
			WillReturnRows(rows)

		mock.ExpectBegin()

		// 模拟 UPDATE：使用 version=5 做条件，version 更新为 6
		mock.ExpectExec("UPDATE `books` SET `count`=?,`isbn`=?,`title`=?,`version`=?,`updated_at`=? WHERE (id = ? AND version = ?) AND `books`.`deleted_at` IS NULL").
			WithArgs(100, "1234567890", "新标题", 6, sqlmock.AnyArg(), 1, 5).
			WillReturnResult(sqlmock.NewResult(1, 1)) // 影响 1 行

		mock.ExpectCommit()
		err = BookUpdateDAO(gormDB, dto)
		assert.NoError(t, err)

		// 验证所有期望都被满足
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("更新失败：乐观锁冲突，version 不匹配", func(t *testing.T) {
		gormDB, mock, err := getDBMock()
		assert.NoError(t, err)

		dto := &dto.BookUpdateDTO{
			ID: 1,
			BookInfoDTO: dto.BookInfoDTO{
				Title: "新标题",
				Count: 100,
				ISBN:  "1234567890",
			},
		}

		// 模拟查询：读取 version=5
		rows := sqlmock.NewRows([]string{"id", "title", "count", "isbn", "version"}).
			AddRow(1, "旧标题", 50, "0987654321", 5)

		mock.ExpectQuery("SELECT * FROM `books` WHERE id = ? AND `books`.`deleted_at` IS NULL ORDER BY `books`.`id` LIMIT ?").
			WithArgs(1, 1).
			WillReturnRows(rows)

		mock.ExpectBegin()
		// 模拟 UPDATE：version=5 不匹配（比如已被别人改到 6），影响行数为 0
		mock.ExpectExec("UPDATE `books` SET `count`=?,`isbn`=?,`title`=?,`version`=?,`updated_at`=? WHERE (id = ? AND version = ?) AND `books`.`deleted_at` IS NULL").
			WithArgs(100, "1234567890", "新标题", 6, sqlmock.AnyArg(), 1, 5).
			WillReturnResult(sqlmock.NewResult(0, 0)) // 影响 0 行 → 乐观锁失败

		mock.ExpectCommit()
		err = BookUpdateDAO(gormDB, dto)
		assert.Error(t, err)
		assert.Equal(t, "更新失败：数据已被其他用户修改，请刷新后重试", err.Error())

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("查询失败：记录不存在", func(t *testing.T) {
		gormDB, mock, err := getDBMock()
		assert.NoError(t, err)

		dto := &dto.BookUpdateDTO{ID: 999}

		// 模拟查询：没有找到记录
		mock.ExpectQuery("SELECT * FROM `books` WHERE id = ? AND `books`.`deleted_at` IS NULL ORDER BY `books`.`id` LIMIT ?").
			WithArgs(999, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		err = BookUpdateDAO(gormDB, dto)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
