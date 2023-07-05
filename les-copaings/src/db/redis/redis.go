package redis

import (
	"context"
	"fmt"
	"github.com/anhgelus/discord-bots/les-copaings/src/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

type RedisCredentials struct {
	Address  string
	Password string
	DB       int
}

type ConnectedUser struct {
	UserID        string
	GuildID       string
	IsConnected   bool
	TimeConnected uint
}

var Credentials RedisCredentials

var Ctx = context.Background()

func (rc *RedisCredentials) GetClient() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     rc.Address,
		Password: rc.Password,
		DB:       rc.DB,
	})
	err := client.Ping(Ctx).Err()
	return client, err
}

func GenerateConnectedUser(member *discordgo.Member) ConnectedUser {
	client, _ := Credentials.GetClient()
	defer client.Close()
	guildID := member.GuildID
	userID := member.User.ID
	connect := client.Get(Ctx, genKey(guildID, userID, "connected"))
	user := ConnectedUser{
		UserID:      userID,
		GuildID:     guildID,
		IsConnected: connect.Val() == "true",
	}
	user.GenerateTimeConnected()
	return user
}

func (user *ConnectedUser) Connect() {
	client, _ := Credentials.GetClient()
	defer client.Close()
	user.IsConnected = true
	user.TimeConnected = 0
	client.Set(Ctx, user.genKey("connected"), "true", 0)
	client.Set(Ctx, user.genKey("connect_at"), time.Now().Unix(), 0)
}

func (user *ConnectedUser) Disconnect() {
	client, _ := Credentials.GetClient()
	defer client.Close()

	user.GenerateTimeConnected()
	user.IsConnected = false

	client.Set(Ctx, user.genKey("connected"), "true", 0)
	client.Set(Ctx, user.genKey("connect_at"), 0, 0)
}

func genKey(guildID string, userID string, ext string) string {
	return fmt.Sprintf("%s:%s:%s", guildID, userID, ext)
}

func (user *ConnectedUser) genKey(ext string) string {
	return fmt.Sprintf("%s:%s:%s", user.GuildID, user.UserID, ext)
}

func (user *ConnectedUser) GenerateTimeConnected() {
	if !user.IsConnected {
		user.TimeConnected = 0
		return
	}
	client, _ := Credentials.GetClient()
	defer client.Close()
	connectAtStr := client.Get(Ctx, genKey(user.GuildID, user.UserID, "connect_at"))
	connectAt, err := strconv.Atoi(connectAtStr.Val())
	if err != nil {
		utils.SendAlert("redis.go - Str to Int Conversion", err.Error())
		return
	}
	user.TimeConnected = CalcTime(uint(connectAt))
}

func CalcTime(connectAt uint) uint {
	timeConnected := uint(time.Now().Unix() - int64(connectAt))
	// Limit the time connect to 6 hours
	if 21600 < timeConnected {
		return 21600
	}
	return timeConnected
}
