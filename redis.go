package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/go-redis/redis"
)

const (
	ShortlinkKey = "shortlink:%s"
	LineIDKey    = "lineId:%s"
)

type RedisCli struct {
	Cli    *redis.Client
	Status bool
}

// ShortlinkInfo short_link info
type ShortlinkInfo struct {
	Status    bool   `json:"Status"`
	Count     int    `json:"Count"`
	LineId    string `json:"LineId"`
	ShortLink string `json:"ShortLink"`
	Prize     *Prize `json:"Prize"`
	LuckDate  string `json:"LuckDate"`
	CreatedAt string `json:"CreatedAt"`
}

var rdbConnCount int

// InitializeRedis 初始化Redis连接
func (r *RedisCli) InitializeRedis() (err error) {
	if r.Status {
		return nil
	}
	r.Cli = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	rdbConnCount++
	_, err = r.Cli.Ping().Result()
	if err != nil {
		fmt.Printf("Redis conn failed %d times\n", rdbConnCount)
		if rdbConnCount >= 3 {
			panic(err)
			// return err
		}
		<-time.After(time.Second * 1)
		return r.InitializeRedis()
	}
	r.Status = true
	return nil
}

// Exists方法封装
func (r *RedisCli) Exists(key string) (bool, error) {
	v, err := r.Cli.Exists(key).Result()
	if err != nil {
		return false, err
	}

	return v > 0, nil
}

func (r *RedisCli) Shorten(lineId string) (*ShortlinkInfo, *StatusError) {
	if len(lineId) < 4 {
		return nil, LogicError(errors.New("lineID長度太小！！"))
	}
	if len(lineId) > 8 {
		return nil, LogicError(errors.New("lineID長度太大！！"))
	}
	b, err := r.Exists(fmt.Sprintf(LineIDKey, lineId))
	if err != nil {
		return nil, ServerError(err)
	}
	if b {
		return nil, LogicError(errors.New("lineID已經存在"))
	}

	eid := r.genstr(lineId, 8)

	shortLinkInfo := &ShortlinkInfo{
		Status:    false,
		Count:     0,
		LineId:    lineId,
		ShortLink: eid,
		Prize:     nil,
		LuckDate:  "",
		CreatedAt: time.Now().String(),
	}

	jsonStr, err := json.Marshal(shortLinkInfo)
	if err != nil {
		return nil, ServerError(err)
	}

	err = r.Cli.Set(fmt.Sprintf(LineIDKey, lineId), eid, 0).Err()
	if err != nil {
		return nil, ServerError(err)
	}

	err = r.Cli.Set(fmt.Sprintf(ShortlinkKey, eid), jsonStr, 0).Err()
	if err != nil {
		return nil, ServerError(err)
	}

	return shortLinkInfo, nil
}

func (r *RedisCli) GetShortlinkInfo(eid string) (*ShortlinkInfo, *StatusError) {
	jsonStr, err := r.Cli.Get(fmt.Sprintf(ShortlinkKey, eid)).Result()
	if err == redis.Nil {
		return nil, &StatusError{Code: 404, Err: errors.New("unknown short url")}
	} else if err != nil {
		return nil, &StatusError{Code: 500, Err: fmt.Errorf("redis error: %s", err)}
	}
	res := &ShortlinkInfo{}
	if e := json.Unmarshal([]byte(jsonStr), res); e != nil {
		return nil, &StatusError{Code: 500, Err: fmt.Errorf("json.Unmarshal failed:%s", e)}
	}
	return res, nil
}

func (r *RedisCli) getUrls() ([]*ShortlinkInfo, error) {
	urls, err := r.Cli.Keys("shortlink:*").Result()
	if err != nil {
		return nil, err
	}

	res := []*ShortlinkInfo{}
	for _, url := range urls {
		infoStr, err := r.Cli.Get(url).Result()
		if err != nil {
			return nil, err
		}
		info := &ShortlinkInfo{}
		if err = json.Unmarshal([]byte(infoStr), info); err != nil {
			return nil, StatusError{Code: 500, Err: fmt.Errorf("json.Unmarshal failed:%s", err)}
		}
		res = append(res, info)
	}
	return res, nil
}

func (r *RedisCli) genstr(str string, length int) string {
	m := len(str)
	if m > 8 || m < 4 {
		return ""
	}
	Base62 := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte{}
	rd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		bytes = append(bytes, Base62[rd.Intn(62)])
	}
	eid := string(bytes)
	b, _ := r.Exists(fmt.Sprintf(ShortlinkKey, eid))
	if b {
		eid = r.genstr(str, length)
	}
	return eid
}
