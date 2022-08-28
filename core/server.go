package core

import (
	"TelegramBot/config"
	"TelegramBot/core/handlers"
	"fmt"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/xelaj/mtproto/telegram"
	"gopkg.in/telebot.v3"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
)

const userkey = "user"

var secret = []byte("secret")

type Server struct {
	port   uint32
	client *telegram.Client
	api    *telebot.Bot
}

func New(port uint32, api *telebot.Bot) *Server {
	return &Server{
		port: port,
		api:  api,
	}
}

func (s *Server) Run(cfg *config.Config) {
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

		context.Redirect(http.StatusPermanentRedirect, "/")
	})

	private.POST("/login", func(context *gin.Context) {
		phone := context.PostForm("phone")

		phoneLen := len(phone)
		if phoneLen == 0 || phoneLen > 12 || phoneLen < 11 {
			context.String(http.StatusBadRequest, "Номер телефона некорректен")
			return
		}

		setCode, err := s.client.AuthSendCode(phone, int32(cfg.AppID), cfg.AppHash, &telegram.CodeSettings{})
		if err != nil {
			context.String(http.StatusInternalServerError, err.Error())
			return
		}

		context.Data(http.StatusOK, "text/html; charset=utf-8", []byte(fmt.Sprintf(tgAuthCode, phone, setCode.PhoneCodeHash)))
	})

	private.POST("/login/submit", func(context *gin.Context) {
		phone := context.PostForm("phone")
		phoneLen := len(phone)
		if phoneLen == 0 || phoneLen > 12 || phoneLen < 11 {
			context.String(http.StatusInternalServerError, "Произошла непредвиденная ошибка")
			return
		}

		hash := context.PostForm("hash")
		hashLen := len(hash)
		if hashLen == 0 {
			context.String(http.StatusInternalServerError, "Произошла непредвиденная ошибка")
			return
		}

		code := context.PostForm("code")
		codeLen := len(code)
		if codeLen == 0 || codeLen > 5 || codeLen < 5 {
			context.String(http.StatusBadRequest, "Некорректный код")
			return
		}

		_, err := s.client.AuthSignIn(
			phone,
			hash,
			code,
		)
		if err != nil {
			if strings.Contains(err.Error(), "Two-steps verification is enabled and a password is required") {
				context.Data(http.StatusOK, "text/html; charset=utf-8", []byte(tgPasswd))

				return
			}

			context.String(http.StatusInternalServerError, err.Error())
			return
		}
	})

	private.POST("/login/password", func(context *gin.Context) {
		password := context.PostForm("password")
		if len(password) == 0 {
			context.String(http.StatusBadRequest, "Некорректный пароль")
			return
		}

		accountPassword, err := s.client.AccountGetPassword()
		if err != nil {
			context.String(http.StatusInternalServerError, err.Error())
			return
		}

		// GetInputCheckPassword is fast response object generator
		inputCheck, err := telegram.GetInputCheckPassword(password, accountPassword)
		if err != nil {
			context.String(http.StatusInternalServerError, err.Error())
			return
		}

		_, err = s.client.AuthCheckPassword(inputCheck)
		if err != nil {
			context.String(http.StatusInternalServerError, err.Error())
			return
		}

		context.Redirect(http.StatusPermanentRedirect, "/")
	})

	private.GET("/volgograd", func(context *gin.Context) {
		_, err := s.getReports(cfg, handlers.VlgName)
		if err != nil {
			fmt.Println(err)
		}
	})

	private.GET("/krasnodar", func(context *gin.Context) {
		_, err := s.getReports(cfg, handlers.KrdName)
		if err != nil {
			fmt.Println(err)
		}
	})

	private.GET("/moscow", func(context *gin.Context) {
		_, err := s.getReports(cfg, handlers.MskName)
		if err != nil {
			fmt.Println(err)
		}
	})

	err := engine.Run(fmt.Sprintf(":%d", s.port))
	if err != nil {
		log.Println(err)
	}
}

func (s *Server) connectApp(cfg *config.Config) error {
	client, err := telegram.NewClient(telegram.ClientConfig{
		// where to store session configuration. must be set
		SessionFile: "storage/session.json",
		// host address of mtproto server. Actually, it can be any mtproxy, not only official
		ServerHost: cfg.Server,
		// public keys file is path to file with public keys, which you must get from https://my.telegram.org
		PublicKeysFile:  cfg.Key,
		AppID:           cfg.AppID,   // app id, could be find at https://my.telegram.org
		AppHash:         cfg.AppHash, // app hash, could be find at https://my.telegram.org
		InitWarnChannel: true,        // if we want to get errors, otherwise, client.Warnings will be set nil
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

		if s.client == nil {
			err := s.connectApp(cfg)
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
			}
		}

		_, err := s.client.AccountGetAccountTtl()
		if err != nil && !strings.Contains(c.Request.URL.String(), "login") {
			c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(tgAuth))
			c.Abort()
			return
		}

		// Continue down the chain to handler etc
		c.Next()
	}
}

func (s *Server) getReports(cfg *config.Config, city string) ([]byte, error) {
	channels := cfg.Channels[city]

	for _, channel := range channels {
		id, err := strconv.Atoi(channel.ID)

		chat, err := s.client.GetChatByID(id)
		if err != nil {
			return nil, fmt.Errorf("can't get chat: %w", err)
		}

		channelSimpleData, ok := chat.(*telegram.Channel)
		if !ok {
			return nil, fmt.Errorf("not a chan")
		}

		msgs, err := s.client.MessagesGetHistory(&telegram.MessagesGetHistoryParams{
			Peer:       &telegram.InputPeerChannel{ChannelID: channelSimpleData.ID, AccessHash: channelSimpleData.AccessHash},
			OffsetID:   0,
			OffsetDate: 0,
			AddOffset:  0,
			Limit:      math.MaxInt32,
			MaxID:      math.MaxInt32,
			MinID:      0,
			Hash:       0,
		})
		if err != nil {
			return nil, fmt.Errorf("can't get gistory %w", err)
		}

		s.getMessages(msgs)
	}

	return nil, nil
}

func (s *Server) getMessages(msgs telegram.MessagesMessages) {
	chanMsgs, ok := msgs.(*telegram.MessagesChannelMessages)
	fmt.Println(chanMsgs, ok)

	for _, untypedMsg := range chanMsgs.Messages {
		msg, ok := untypedMsg.(*telegram.MessageObj)
		if !ok {
			continue
		}

		fmt.Println(msg.Message)

		if msg.Media != nil {
			untypedPhoto, ok := msg.Media.(*telegram.MessageMediaPhoto)
			if ok {
				photo, ok := untypedPhoto.Photo.(*telegram.PhotoObj)
				if ok {

					file, err := s.client.UploadGetFile(&telegram.UploadGetFileParams{
						Location: &telegram.InputPhotoFileLocation{
							ID:            photo.ID,
							AccessHash:    photo.AccessHash,
							FileReference: photo.FileReference,
						},
						Limit: 1048576,
					})
					if err != nil {
						fmt.Println(err)
						continue
					}

					fmt.Println(file)
				}
			}
		}
	}
}
