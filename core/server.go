package core

import (
	"TelegramBot/config"
	"TelegramBot/core/database"
	"TelegramBot/tgerror"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	scribble "github.com/nanobox-io/golang-scribble"
	"github.com/xelaj/mtproto/telegram"
	"gopkg.in/telebot.v3"
	"log"
	"net/http"
)

const (
	userkey      = "user"
	ErrNoReports = tgerror.TelegramError("Отчёты не найдены")
)

var secret = []byte("secret")

type Server struct {
	port   uint32
	client *telegram.Client
	db     *scribble.Driver
}

func New(port uint32, db *scribble.Driver) *Server {
	return &Server{
		port: port,
		db:   db,
	}
}

func (s *Server) Run(cfg *config.Config, api *telebot.Bot) {
	engine := gin.Default()

	engine.Use(sessions.Sessions("mysession", sessions.NewCookieStore(secret)))

	private := engine.Group("/")
	private.Use(s.AuthRequired(cfg))

	private.GET("/", func(context *gin.Context) {
		context.Data(http.StatusOK, "text/html; charset=utf-8", []byte(cities))
	})

	engine.POST("/", func(context *gin.Context) {
		session := sessions.Default(context)

		login := context.PostForm("login")
		password := context.PostForm("password")

		if login != cfg.Login || password != cfg.Password {
			context.String(http.StatusUnauthorized, "Incorrect login or password")

			return
		}

		session.Set(userkey, login)
		if err := session.Save(); err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
			return
		}

		context.Data(http.StatusOK, "text/html; charset=utf-8", []byte(cities))
	})

	private.GET("/volgograd", func(context *gin.Context) {
		reports, err := s.getReports(database.TableVLG, cfg.Token, api)
		if err != nil {
			context.String(http.StatusInternalServerError, err.Error())
			log.Println(err.Error())
			return
		}

		context.Data(http.StatusOK, "text/html; charset=utf-8", reports)
	})

	private.GET("/krasnodar", func(context *gin.Context) {
		reports, err := s.getReports(database.TableKRD, cfg.Token, api)
		if err != nil {
			context.String(http.StatusInternalServerError, err.Error())
			log.Println(err.Error())
			return
		}

		context.Data(http.StatusOK, "text/html; charset=utf-8", reports)
	})

	private.GET("/moscow", func(context *gin.Context) {
		reports, err := s.getReports(database.TableMSK, cfg.Token, api)
		if err != nil {
			context.String(http.StatusInternalServerError, err.Error())
			log.Println(err.Error())
			return
		}

		context.Data(http.StatusOK, "text/html; charset=utf-8", reports)
	})

	private.GET("/logout", func(context *gin.Context) {
		session := sessions.Default(context)
		user := session.Get(userkey)
		if user == nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session token"})
			return
		}

		session.Delete(userkey)
		if err := session.Save(); err != nil {
			log.Println("failed to save user session")
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
			return
		}

		context.Redirect(http.StatusTemporaryRedirect, "/")
	})
	err := engine.Run(fmt.Sprintf(":%d", s.port))
	if err != nil {
		log.Println(err)
	}
}

func (s *Server) connectApp(cfg *config.Config) error {
	client, err := telegram.NewClient(telegram.ClientConfig{
		SessionFile:     "storage/session.json",
		ServerHost:      cfg.Server,
		PublicKeysFile:  cfg.Key,
		AppID:           cfg.AppID,
		AppHash:         cfg.AppHash,
		InitWarnChannel: true,
	})

	s.client = client

	return err
}

func (s *Server) AuthRequired(cfg *config.Config) func(c *gin.Context) {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get(userkey)
		if user == nil {
			// Abort the request with the appropriate error code
			c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(login))
			c.Abort()
			return
		}

		// Continue down the chain to handler etc
		c.Next()
	}
}

func (s *Server) getReports(table string, token string, api *telebot.Bot) ([]byte, error) {
	records, err := s.db.ReadAll(table)
	if err != nil {
		return nil, ErrNoReports
	}

	var recordsHTML string

	for _, recordJSON := range records {
		var record *database.Message

		err = json.Unmarshal([]byte(recordJSON), &record)
		if err != nil {
			return nil, err
		}

		var media string
		if len(record.Media) != 0 {
			file, err := api.FileByID(record.Media)
			if err != nil {
				return nil, err
			}

			media = fmt.Sprintf(fmtMedia, fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", token, file.FilePath))
		}

		var text string
		if len(record.Text) != 0 {
			text = fmt.Sprintf(fmtText, record.Text)
		}

		recordsHTML += fmt.Sprintf(fmtReport, record.Place, record.From, record.Timestamp.String(), media+text)
	}

	return []byte(fmt.Sprintf(fmtReportBody, recordsHTML)), nil
}
