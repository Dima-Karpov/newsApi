package repository

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"newsApi/internal/domain"
	"testing"
	"time"
)

func TestNewsRepository_Save(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		news *domain.RSSItem
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "successfully save news",
			fields: fields{
				db: createMockDB(t, "SELECT VERSION()"), // создаем mock DB
			},
			args: args{
				news: &domain.RSSItem{
					Title:       "Test News",
					Description: "Test Description",
					Link:        "http://example.com",
					PublishedAt: time.Now(),
				},
			},
			wantErr: false,
		},
		{
			name: "news already exists",
			fields: fields{
				db: createMockDB(t, "SELECT VERSION()"),
			},
			args: args{
				news: &domain.RSSItem{
					Title:       "Test News",
					Description: "Test Description",
					Link:        "http://example.com", // с таким же линком
					PublishedAt: time.Now(),
				},
			},
			wantErr: true, // ошибка, так как новость уже существует
		},
		{
			name: "error when saving news",
			fields: fields{
				db: createMockDBWithError(t, "SELECT VERSION()"),
			},
			args: args{
				news: &domain.RSSItem{
					Title:       "Test News",
					Description: "Test Description",
					Link:        "http://example.com",
					PublishedAt: time.Now(),
				},
			},
			wantErr: true, // ошибка при сохранении новости
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &NewsRepository{
				db: tt.fields.db,
			}
			if err := r.Save(tt.args.news); (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// createMockDB создает mock базу данных для успешного выполнения запросов.
func createMockDB(t *testing.T, expectedQuery string) *gorm.DB {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock db: %v", err)
	}

	// Ожидаем запрос версии базы данных
	mock.ExpectQuery(expectedQuery).WillReturnRows(sqlmock.NewRows([]string{"version"}).AddRow("5.7.29"))

	// Возвращаем mock DB, преобразованный в *gorm.DB
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to initialize GORM DB: %v", err)
	}

	return gormDB
}

// createMockDBWithError создает mock базу данных с ошибкой при выполнении запроса.
func createMockDBWithError(t *testing.T, expectedQuery string) *gorm.DB {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock db: %v", err)
	}

	// Ожидаем запрос версии базы данных
	mock.ExpectQuery(expectedQuery).WillReturnRows(sqlmock.NewRows([]string{"version"}).AddRow("5.7.29"))

	// Ожидаем ошибку при сохранении
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `news_lists`").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnError(errors.New("error saving news"))
	mock.ExpectCommit()

	// Возвращаем mock DB, преобразованный в *gorm.DB
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to initialize GORM DB: %v", err)
	}

	return gormDB
}

func TestNewsRepository_GetNews(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		page     int
		pageSize int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []domain.NewsList
		total   int
		wantErr bool
	}{
		{
			name: "successfully get news with pagination",
			fields: fields{
				db: createMockDBForGetNews(t),
			},
			args: args{
				page:     1,
				pageSize: 10,
			},
			want: []domain.NewsList{
				{Title: "Test News 1", Link: "http://example.com/1", PublishedAt: time.Now()},
				{Title: "Test News 2", Link: "http://example.com/2", PublishedAt: time.Now()},
			},
			total:   20, // общее количество новостей
			wantErr: false,
		},
		{
			name: "error when getting news",
			fields: fields{
				db: createMockDBWithErrorForGetNews(t),
			},
			args: args{
				page:     1,
				pageSize: 10,
			},
			wantErr: true, // ошибка при извлечении новостей
		},
		{
			name: "error when counting news",
			fields: fields{
				db: createMockDBForGetNewsWithCountError(t),
			},
			args: args{
				page:     1,
				pageSize: 10,
			},
			wantErr: true, // ошибка при подсчете общего количества новостей
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &NewsRepository{
				db: tt.fields.db,
			}
			got, total, err := r.GetNews(tt.args.page, tt.args.pageSize)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNews() error = %v, wantErr %v", err, tt.wantErr)
			}
			if len(got) != len(tt.want) {
				t.Errorf("GetNews() got = %v, want %v", got, tt.want)
			}
			if total != tt.total {
				t.Errorf("GetNews() total = %v, want %v", total, tt.total)
			}
		})
	}
}

// createMockDBForGetNews создает mock DB с успешным возвратом новостей.
func createMockDBForGetNews(t *testing.T) *gorm.DB {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock db: %v", err)
	}

	// Ожидаем запрос на получение новостей
	mock.ExpectQuery("SELECT VERSION()").WillReturnRows(sqlmock.NewRows([]string{"version"}).AddRow("5.7.29"))

	// Ожидаем запрос на получение новостей с пагинацией
	mock.ExpectQuery("SELECT * FROM `news_lists` LIMIT 10 OFFSET 0").WillReturnRows(
		sqlmock.NewRows([]string{"id", "title", "description", "link", "published_at"}).
			AddRow(uuid.New(), "Test News 1", "Test Description", "http://example.com/1", time.Now()).
			AddRow(uuid.New(), "Test News 2", "Test Description", "http://example.com/2", time.Now()),
	)

	// Ожидаем запрос на подсчет общего числа новостей
	mock.ExpectQuery("SELECT count(*) FROM `news_lists`").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(20))

	// Возвращаем mock DB, преобразованный в *gorm.DB
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to initialize GORM DB: %v", err)
	}

	return gormDB
}

// createMockDBWithErrorForGetNews создает mock DB, который вернет ошибку при извлечении новостей.
func createMockDBWithErrorForGetNews(t *testing.T) *gorm.DB {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock db: %v", err)
	}

	// Ожидаем запрос на получение новостей с ошибкой
	mock.ExpectQuery("SELECT * FROM `news_lists` LIMIT 10 OFFSET 0").WillReturnError(errors.New("error getting news"))

	// Возвращаем mock DB, преобразованный в *gorm.DB
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to initialize GORM DB: %v", err)
	}

	return gormDB
}

// createMockDBForGetNewsWithCountError создает mock DB, который вернет ошибку при подсчете общего количества новостей.
func createMockDBForGetNewsWithCountError(t *testing.T) *gorm.DB {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock db: %v", err)
	}

	// Ожидаем запрос на получение новостей
	mock.ExpectQuery("SELECT VERSION()").WillReturnRows(sqlmock.NewRows([]string{"version"}).AddRow("5.7.29"))

	// Ожидаем успешный запрос на получение новостей
	mock.ExpectQuery("SELECT * FROM `news_lists` LIMIT 10 OFFSET 0").WillReturnRows(
		sqlmock.NewRows([]string{"id", "title", "description", "link", "published_at"}).
			AddRow(uuid.New(), "Test News 1", "Test Description", "http://example.com/1", time.Now()).
			AddRow(uuid.New(), "Test News 2", "Test Description", "http://example.com/2", time.Now()),
	)

	// Ожидаем ошибку при подсчете общего числа новостей
	mock.ExpectQuery("SELECT count(*) FROM `news_lists`").WillReturnError(errors.New("error counting news"))

	// Возвращаем mock DB, преобразованный в *gorm.DB
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to initialize GORM DB: %v", err)
	}

	return gormDB
}
