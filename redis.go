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
	URLIDKEY     = "shortlinkkey"
	ShortlinkKey = "shortlink:%s"
)

type RedisCli struct {
	Cli    *redis.Client
	Status bool
}

// ShortlinkInfo short_link info
type ShortlinkInfo struct {
	LineId    string `json:"LineId"`
	ShortLink string `json:"ShortLink"`
	Prize     string `json:"Prize"`
	LuckDate  string `json:"LuckDate"`
	CreatedAt string `json:"created_at"`
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

// TODO 未做lineId唯一限制
func (r *RedisCli) Shorten(lineId string) (string, error) {
	// Incr global counter
	if err := r.Cli.Incr(URLIDKEY).Err(); err != nil {
		return "", err
	}

	// get global counter
	urlId, err := r.Cli.Get(URLIDKEY).Int()
	if err != nil {
		return "", err
	}

	// convert int to short link
	eid := genstr(lineId, 8, urlId)

	// set detail for short link
	shortLinkInfo := &ShortlinkInfo{
		LineId:    lineId,
		ShortLink: eid,
		Prize:     "",
		LuckDate:  "",
		CreatedAt: time.Now().String(),
	}

	// serialize short link info
	jsonStr, err := json.Marshal(shortLinkInfo)
	if err != nil {
		return "", err
	}

	// set key for short link to detail
	err = r.Cli.Set(fmt.Sprintf(ShortlinkKey, eid), jsonStr, 0).Err()
	if err != nil {
		return "", err
	}

	return eid, err
}

func (r *RedisCli) GetShortlinkInfo(eid string) (*ShortlinkInfo, error) {
	jsonStr, err := r.Cli.Get(fmt.Sprintf(ShortlinkKey, eid)).Result()
	if err == redis.Nil {
		return nil, StatusError{Code: 404, Err: errors.New("unknown short url")}
	} else if err != nil {
		return nil, err
	}
	res := &ShortlinkInfo{}
	if e := json.Unmarshal([]byte(jsonStr), res); e != nil {
		return nil, StatusError{Code: 500, Err: fmt.Errorf("json.Unmarshal failed:%s", e)}
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

func genstr(str string, length int, id int) string {
	Base62 := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const key = "licat233"
	n := len(key) - 1
	bytes := []byte{}
	for k, s := range []byte(str) {
		keys := key[k%n]
		ens := s + keys + byte(id)
		ens = ens % 62
		bytes = append(bytes, Base62[ens])
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := len(bytes); i < length; i++ {
		b := r.Intn(62)
		bytes = append(bytes, Base62[b])
	}
	res := string(bytes)
	if len(res) > 8 {
		res = res[:8]
	}
	return res
}
