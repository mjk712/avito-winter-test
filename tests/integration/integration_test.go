package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"

	"avito-winter-test/internal/config"
	httpServer "avito-winter-test/internal/http-server"
	"avito-winter-test/internal/models/dto"
	"avito-winter-test/internal/service"
	"avito-winter-test/internal/storage"
	testhelpers "avito-winter-test/tests/test-helpers"
)

type BaseTestSuite struct {
	suite.Suite
	pgContainer *testhelpers.PostgresContainer
	db          *sqlx.DB
	server      *http.Server
	baseURL     string
}

func (s *BaseTestSuite) setupTestEnvironment(ctx context.Context) {
	var err error

	// Поднимаем контейнер с PostgreSQL
	s.pgContainer, err = testhelpers.CreatePostgresContainer(ctx)
	s.Require().NoError(err)

	// Подключаемся к базе данных
	s.db, err = sqlx.Connect("postgres", s.pgContainer.ConnectionString)
	s.Require().NoError(err)

	// Инициализируем репозиторий и сервис
	repo := &storage.Repository{DB: s.db}
	s.runMigrations()
	merchShopService := service.NewMerchShopService(repo)

	// Создаем конфиг и сервер
	cfg := &config.Config{
		HTTPServer: config.HTTPServer{Address: ":8080"},
		DBConfig:   config.DBConfig{ConnectionString: s.pgContainer.ConnectionString},
		Env:        "local",
	}
	s.baseURL = "http://localhost:8080"

	s.server = httpServer.NewServer(ctx, slog.Default(), cfg, merchShopService)

	// Запускаем сервер
	go func() {
		s.T().Log("Starting test server")
		if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
			s.T().Fatalf("Server failed: %v", err)
		}
	}()

	// Ждем, пока сервер поднимется
	s.Require().Eventually(func() bool {
		resp, err := http.Get(fmt.Sprintf("%s/api/info", s.baseURL))
		return err == nil && resp.StatusCode == http.StatusUnauthorized
	}, 5*time.Second, 100*time.Millisecond, "Server failed to start")
}

func (s *BaseTestSuite) TearDownSuite() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	s.Require().NoError(s.server.Shutdown(ctx))
	s.Require().NoError(s.pgContainer.Terminate(ctx))
}

func (s *BaseTestSuite) runMigrations() {
	m, err := migrate.New("file://../../internal/storage/migrations", s.pgContainer.ConnectionString)
	if err != nil {
		log.Fatalf("Could not connect to migrations: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("failed to run migrations: %v", err)
	}
}

func (s *BaseTestSuite) createTestUser(username, password string) (string, int) {
	reqBody := fmt.Sprintf(`{"username": "%s", "password": "%s"}`, username, password)
	resp, err := http.Post(fmt.Sprintf("%s/api/auth", s.baseURL), "application/json", strings.NewReader(reqBody))
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	var authResp dto.AuthResponse
	s.Require().NoError(json.NewDecoder(resp.Body).Decode(&authResp))

	// Извлекаем userID из токена
	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(authResp.Token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	s.Require().NoError(err)

	return authResp.Token, int(claims["userID"].(float64))
}

func (s *BaseTestSuite) makeRequest(method, url, token string, body any) *http.Response {
	var reqBody *bytes.Reader
	if body != nil {
		data, _ := json.Marshal(body)
		reqBody = bytes.NewReader(data)
	} else {
		reqBody = bytes.NewReader(nil)
	}

	req, _ := http.NewRequest(method, url, reqBody)
	if token != "" {
		req.Header.Add("Authorization", token)
	}
	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	resp, err := http.DefaultClient.Do(req)
	s.Require().NoError(err)
	return resp
}

func (s *BaseTestSuite) getUserInfo(token string) dto.InfoResponse {
	resp := s.makeRequest("GET", fmt.Sprintf("%s/api/info", s.baseURL), token, nil)
	var info dto.InfoResponse
	s.Require().NoError(json.NewDecoder(resp.Body).Decode(&info))
	return info
}

// ---------- Тест покупки ----------

type PurchaseTestSuite struct {
	BaseTestSuite
	authToken    string
	testUserID   int
	testItemName string
}

func (s *PurchaseTestSuite) SetupSuite() {
	s.setupTestEnvironment(context.Background())

	// Создаем тестового пользователя
	s.authToken, s.testUserID = s.createTestUser("testuser", "testpass")
	s.testItemName = "t-shirt"
}

func (s *PurchaseTestSuite) TestFullPurchaseFlow() {
	// Покупка предмета
	resp := s.makeRequest("GET", fmt.Sprintf("%s/api/buy/%s", s.baseURL, s.testItemName), s.authToken, nil)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	// Проверяем баланс
	info := s.getUserInfo(s.authToken)
	s.Equal(1000-80, info.Coins)
	s.Len(info.Inventory, 1)
	s.Equal("t-shirt", info.Inventory[0].Type)
	s.Equal(1, info.Inventory[0].Quantity)
}

func TestPurchaseSuite(t *testing.T) {
	suite.Run(t, new(PurchaseTestSuite))
}

// ---------- Тест перевода монет ----------

type SendCoinTestSuite struct {
	BaseTestSuite
	authToken      string
	testUserID     int
	receiverToken  string
	receiverUserID int
}

func (s *SendCoinTestSuite) SetupSuite() {
	s.setupTestEnvironment(context.Background())

	// Создаем отправителя и получателя
	s.authToken, s.testUserID = s.createTestUser("sender", "senderpass")
	s.receiverToken, s.receiverUserID = s.createTestUser("receiver", "receiverpass")
}

func (s *SendCoinTestSuite) TestFullSendCoinFlow() {
	// Отправляем монеты
	sendCoinReq := dto.SendCoinRequest{ToUser: "receiver", Amount: 100}
	resp := s.makeRequest("POST", fmt.Sprintf("%s/api/send-coin", s.baseURL), s.authToken, sendCoinReq)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	// Проверяем баланс отправителя
	senderInfo := s.getUserInfo(s.authToken)
	s.Equal(1000-100, senderInfo.Coins)

	// Проверяем баланс получателя
	receiverInfo := s.getUserInfo(s.receiverToken)
	s.Equal(1000+100, receiverInfo.Coins)
}

func TestSendCoinSuite(t *testing.T) {
	suite.Run(t, new(SendCoinTestSuite))
}
