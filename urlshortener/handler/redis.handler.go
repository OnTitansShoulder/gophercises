package handler

import (
	"fmt"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

func RedisHandler(redisServerURL string, fallback http.Handler) (http.HandlerFunc, error) {
	opt, err := redis.ParseURL(redisServerURL)
	if err != nil {
		return nil, errors.Wrap(err, "parse redis URL")
	}

	rdb := redis.NewClient(opt)
	handler := func(w http.ResponseWriter, r *http.Request) {
		val, err := rdb.Get(r.URL.Path).Result()
		if err != nil {
			fmt.Println(err.Error())
		} else {
			http.Redirect(w, r, val, 302)
		}
		fallback.ServeHTTP(w, r)
	}
	return handler, nil
}
