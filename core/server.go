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
	"log"
	"net/http"
	"sort"
	"strconv"
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

func (s *Server) Run(cfg *config.Config) {
	engine := gin.Default()

	engine.Use(sessions.Sessions("mysession", sessions.NewCookieStore(secret)))

	private := engine.Group("/")
	private.Use(func(context *gin.Context) {
		session := sessions.Default(context)
		user := session.Get(userkey)
		if user == nil {
			context.Data(http.StatusOK, "text/html; charset=utf-8", []byte(login))
			context.Abort()
			return
		}

		context.Next()
	})

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
		reports, err := s.getReports(context, database.TableVLG, cfg.Token)
		if err != nil {
			context.String(http.StatusInternalServerError, err.Error())
			log.Println(err.Error())
			return
		}

		context.Data(http.StatusOK, "text/html; charset=utf-8", reports)
	})

	private.GET("/krasnodar", func(context *gin.Context) {
		reports, err := s.getReports(context, database.TableKRD, cfg.Token)
		if err != nil {
			context.String(http.StatusInternalServerError, err.Error())
			log.Println(err.Error())
			return
		}

		context.Data(http.StatusOK, "text/html; charset=utf-8", reports)
	})

	private.GET("/moscow", func(context *gin.Context) {
		reports, err := s.getReports(context, database.TableMSK, cfg.Token)
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

func (s *Server) getReports(context *gin.Context, table string, token string) ([]byte, error) {
	records, err := s.db.ReadAll(table)
	if err != nil {
		return nil, ErrNoReports
	}

	recordsReverse := sort.Reverse(sort.StringSlice(records))
	fmt.Println(recordsReverse)

	var recordsHTML string

	from := 0
	to := 100

	minString := context.Query("from")
	maxString := context.Query("to")

	if len(minString) != 0 {
		fromQuery, err := strconv.Atoi(minString)
		if err == nil {
			from = fromQuery
		}
	}

	if len(maxString) != 0 {
		toQuery, err := strconv.Atoi(maxString)
		if err == nil {
			to = toQuery
		}
	}

	countRecords := len(records)

	from, to = validateRange(from, to, countRecords)

	for _, recordJSON := range records[from:to] {
		var record *database.Message

		err = json.Unmarshal([]byte(recordJSON), &record)
		if err != nil {
			return nil, err
		}

		var media string
		if len(record.Media) != 0 {
			media = fmt.Sprintf(fmtMedia, fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", token, record.Media))
		}

		var text string
		if len(record.Text) != 0 {
			text = fmt.Sprintf(fmtText, record.Text)
		}

		recordsHTML += fmt.Sprintf(fmtReport, record.Place, record.From, record.Timestamp.String(), media+text)
	}

	switch {
	case from == 0 && countRecords <= to:
		return []byte(fmt.Sprintf(fmtReportBody, recordsHTML)), nil
	case from == 0 && countRecords > to:
		oldFrom, oldTo := getRange(from, countRecords, 100)
		return []byte(fmt.Sprintf(fmtReportBody, recordsHTML+fmt.Sprintf(fmtLoadReport, fmt.Sprintf("/%s", table), oldFrom, oldTo, "старее"))), nil
	case from > 0 && countRecords <= to:
		newFrom, newTo := getRange(from, countRecords, -100)
		return []byte(fmt.Sprintf(fmtReportBody, fmt.Sprintf(fmtLoadReport, fmt.Sprintf("/%s", table), newFrom, newTo, "новее")+recordsHTML)), nil
	case from > 0 && countRecords > to:
		newFrom, newTo := getRange(from, countRecords, -100)
		oldFrom, oldTo := getRange(from, countRecords, 100)

		return []byte(fmt.Sprintf(fmtReportBody, fmt.Sprintf(fmtLoadReport, fmt.Sprintf("/%s", table), newFrom, newTo, "новее")+recordsHTML+fmt.Sprintf(fmtLoadReport, fmt.Sprintf("/%s", table), oldFrom, oldTo, "старее"))), nil
	}

	return []byte(fmt.Sprintf(fmtReportBody, recordsHTML)), nil
}

func getRange(from, max, step int) (int, int) {
	newFrom := from + step
	if newFrom < 0 {
		newFrom = 0
	}

	newTo := newFrom + 100
	if newTo > max {
		newTo = max
	}

	return newFrom, newTo
}

func validateRange(from, to, max int) (int, int) {
	if from < 0 {
		from = 0
	}

	if to < 0 {
		to = from + 100
	}

	if to > max {
		to = max
	}

	if from > to {
		if to-100 > 0 {
			from = to - 100
		} else {
			from = 0
		}
	}

	return from, to
}
