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
	XpLostSaved   uint
}

var Credentials RedisCredentials

var Ctx = context.Background()

const (
	xpLostSavedKey = "xp_lost_saved"
	connectedKey   = "connected"
	connectAtKey   = "connect_at"
	lastEventKey   = "last_event"
)

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
	connect := client.Get(Ctx, genKey(guildID, userID, connectedKey))
	raw := client.Get(Ctx, genKey(guildID, userID, xpLostSavedKey))
	var xpLostSaved uint
	if raw.Err() == redis.Nil {
		xpLostSaved = 0
	} else if raw.Err() == nil {
		xpLostSavedStr := raw.Val()
		t, err := strconv.Atoi(xpLostSavedStr)
		if err != nil {
			utils.SendAlert("redis.go - Str to Int Conversion for GenerateConnectedUser", err.Error())
			return ConnectedUser{}
		}
		xpLostSaved = uint(t)
	} else {
		utils.SendAlert("redis.go - Error while fetching xp lost saved", raw.Err().Error())
	}
	user := ConnectedUser{
		UserID:      userID,
		GuildID:     guildID,
		IsConnected: connect.Val() == "true",
		XpLostSaved: xpLostSaved,
	}
	last := client.Get(Ctx, user.genKey(lastEventKey))
	if last.Err() == redis.Nil {
		user.UpdateLastEvent()
	}
	user.GenerateTimeConnected()
	return user
}

func (user *ConnectedUser) Connect() {
	client, _ := Credentials.GetClient()
	defer client.Close()
	user.IsConnected = true
	user.TimeConnected = 0
	client.Set(Ctx, user.genKey(connectedKey), "true", 0)
	client.Set(Ctx, user.genKey(connectAtKey), time.Now().Unix(), 0)
}

func (user *ConnectedUser) Disconnect() {
	client, _ := Credentials.GetClient()
	defer client.Close()

	user.GenerateTimeConnected()
	user.IsConnected = false

	client.Set(Ctx, user.genKey(connectedKey), "false", 0)
	client.Set(Ctx, user.genKey(connectAtKey), 0, 0)
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
	connectAtStr := client.Get(Ctx, genKey(user.GuildID, user.UserID, connectAtKey))
	connectAt, err := strconv.Atoi(connectAtStr.Val())
	if err != nil {
		utils.SendAlert("redis.go - Str to Int Conversion for Time Connected", err.Error())
		return
	}
	user.TimeConnected = CalcTime(uint(connectAt))
}

func (user *ConnectedUser) UpdateLastEvent() {
	client, _ := Credentials.GetClient()
	defer client.Close()

	client.Set(Ctx, user.genKey(lastEventKey), time.Now().Unix(), 0)
	client.Del(Ctx, user.genKey(xpLostSavedKey))
}

func (user *ConnectedUser) TimeSinceLastEvent() int64 {
	client, _ := Credentials.GetClient()
	defer client.Close()

	lastStr := client.Get(Ctx, user.genKey(lastEventKey))
	if lastStr.Err() == redis.Nil {
		return 0
	}
	last, err := strconv.Atoi(lastStr.Val())
	if err != nil {
		utils.SendAlert("redis.go - Str to Int Conversion", err.Error())
		return 0
	}
	return time.Now().Unix() - int64(last)
}

func (user *ConnectedUser) UpdateLostXp(xp uint) {
	client, _ := Credentials.GetClient()
	defer client.Close()

	user.XpLostSaved += xp

	client.Set(Ctx, user.genKey(xpLostSavedKey), fmt.Sprintf("%d", user.XpLostSaved), 0)
}

func (user *ConnectedUser) LeaveGuild() {

}

func CalcTime(connectAt uint) uint {
	timeConnected := uint(time.Now().Unix() - int64(connectAt))
	// Limit the time connected to 6 hours
	if 21600 < timeConnected {
		return 21600
	}
	return timeConnected
}
