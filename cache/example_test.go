package cache

import (
	"fmt"
	"time"
)

func ExampleCache_DecrWithLimit() {
	conf := &Config{
		ServerAddr:      "127.0.0.1:6379",
		MaxIdle:         50,
		MaxActive:       400,
		IdleTimeout:     time.Duration(300) * time.Second,
		MaxConnLifetime: time.Duration(300) * time.Second,
	}
	c, _ := New(Redis, conf)
	c.Set("nuladganteng", "0", time.Duration(300)*time.Second)
	res, err := c.DecrWithLimit("nuladganteng", 10, -5)
	fmt.Println(res)
	fmt.Println(err)
}

func ExampleCache_HGetSet() {
	conf := &Config{
		ServerAddr:      "127.0.0.1:6379",
		MaxIdle:         50,
		MaxActive:       400,
		IdleTimeout:     time.Duration(300) * time.Second,
		MaxConnLifetime: time.Duration(300) * time.Second,
	}
	c, _ := New(Redis, conf)
	c.HSet("nulad", "ganteng", "positif", time.Duration(300)*time.Second)
	err := c.HGetSet("nulad", "ganteng", "positif", "negatif", time.Duration(300)*time.Second)
	fmt.Println(err)
}

func ExampleCache_ZCard() {
	conf := &Config{
		ServerAddr:      "127.0.0.1:6379",
		MaxIdle:         50,
		MaxActive:       400,
		IdleTimeout:     time.Duration(300) * time.Second,
		MaxConnLifetime: time.Duration(300) * time.Second,
	}
	c, _ := New(Redis, conf)
	c.ZAddNX("key", "member_1", 100)
	c.ZAddNX("key", "member_2", 101)
	card, err := c.ZCard("key")
	fmt.Println(card, err)
}

func ExampleCache_TTL() {
	conf := &Config{
		ServerAddr:      "127.0.0.1:6379",
		MaxIdle:         50,
		MaxActive:       400,
		IdleTimeout:     time.Duration(300) * time.Second,
		MaxConnLifetime: time.Duration(300) * time.Second,
	}
	c, _ := New(Redis, conf)
	c.ZAddNX("key", "member_1", 100)
	ttl, err := c.TTL("key")
	fmt.Println(ttl, err)
}

func ExampleCacheCluster() {
	config := ConfigCacheCluster{
		Addrs: []string{"127.0.0.1:6001", "127.0.0.1:6002", "127.0.0.1:6003",
			"127.0.0.1:6005", "127.0.0.1:6007"},
		MaxRedirects:   8,
		ReadOnly:       false,
		RouteByLatency: true,
		RouteRandomly:  false,
	}
	currentClient, err := New(RedisCluster, config)

	if err != nil {
		panic(err)
	}

	err = currentClient.Set("testnumbertwo", "firstvalue", 60*time.Minute)

	result, err := currentClient.Get("testnumbertwo")

	if err != nil {
		panic(err)
	}

	if result != "firstvalue" {
		fmt.Print("Result not expected")
	}
}

func ExampleCacheScript() {
	config := ConfigCacheCluster{
		Addrs: []string{"127.0.0.1:6001", "127.0.0.1:6002", "127.0.0.1:6003",
			"127.0.0.1:6005", "127.0.0.1:6007"},
		MaxRedirects:   8,
		ReadOnly:       false,
		RouteByLatency: true,
		RouteRandomly:  false,
	}
	currentClient, err := New(RedisCluster, config)

	if err != nil {
		panic(err)
	}

	currentClient.HSet("xxxxx_counter", "firstcolumn", "last", 5)
	err = currentClient.HGetSet("xxxxx_counter", "firstcolumn", "thirdvalue", "last", 5)
	if err != nil {
		fmt.Print("result not expected with err : ", err)
	}
}

func ExampleRedigoCluster_SetNX() {
	config := ConfigCacheCluster{
		Addrs: []string{"127.0.0.1:6001", "127.0.0.1:6002", "127.0.0.1:6003",
			"127.0.0.1:6005", "127.0.0.1:6007"},
		MaxRedirects:   8,
		ReadOnly:       false,
		RouteByLatency: true,
		RouteRandomly:  false,
	}
	currentClient, err := New(RedisCluster, config)

	err = currentClient.SetNX("lambo", "firstcolumn", 5)

	if err != nil {
		fmt.Print("result not expected")
	}
}

func ExampleRedigoCluster_Set() {
	config := ConfigCacheCluster{
		Addrs: []string{"127.0.0.1:6001", "127.0.0.1:6002", "127.0.0.1:6003",
			"127.0.0.1:6005", "127.0.0.1:6007"},
		MaxRedirects:   8,
		ReadOnly:       false,
		RouteByLatency: true,
		RouteRandomly:  false,
	}

	currentClient, _ := New(RedisCluster, config)

	currentClient.Set("firstuser", "hisvalue", 5)

	result, _ := currentClient.Get("firstuser")

	if result != "hisvalue" {
		fmt.Print("result not expected : ", result)
	}
}

func ExampleRedigoCluster_HSet() {
	config := ConfigCacheCluster{
		Addrs: []string{"127.0.0.1:6001", "127.0.0.1:6002", "127.0.0.1:6003",
			"127.0.0.1:6005", "127.0.0.1:6007"},
		MaxRedirects:   8,
		ReadOnly:       false,
		RouteByLatency: true,
		RouteRandomly:  false,
	}

	currentClient, err := New(RedisCluster, config)

	err = currentClient.HSet("leftuser", "firstcolumn", "firstvalue", 5)

	result, err := currentClient.HGet("leftuser", "firstcolumn")

	if result != "firstvalue" {
		fmt.Print("Result not expected : ", err)
	}
}

func ExampleRedigoCluster_HSetNX() {
	config := ConfigCacheCluster{
		Addrs: []string{"127.0.0.1:6001", "127.0.0.1:6002", "127.0.0.1:6003",
			"127.0.0.1:6005", "127.0.0.1:6007"},
		MaxRedirects:   8,
		ReadOnly:       false,
		RouteByLatency: true,
		RouteRandomly:  false,
	}

	currentClient, err := New(RedisCluster, config)

	err = currentClient.HSetNX("seconduser", "secondcolumn", "secondvalue", 5)

	if err != nil {
		fmt.Print("Result not expected : ", err)
	}
}

func ExampleRedigoCluster_HMSet() {
	config := ConfigCacheCluster{
		Addrs: []string{"127.0.0.1:6001", "127.0.0.1:6002", "127.0.0.1:6003",
			"127.0.0.1:6005", "127.0.0.1:6007"},
		MaxRedirects:   8,
		ReadOnly:       false,
		RouteByLatency: true,
		RouteRandomly:  false,
	}

	currentClient, err := New(RedisCluster, config)

	err = currentClient.HMSet("checkhasher", map[string]string{"fifthvalue": "fourthvalue"}, 5)

	if err != nil {
		fmt.Print("Result not expected : ", err)
	}
}

func ExampleRedigoCluster_HMGet() {
	config := &ConfigCacheCluster{
		Addrs: []string{"127.0.0.1:6001", "127.0.0.1:6002", "127.0.0.1:6003",
			"127.0.0.1:6005", "127.0.0.1:6007"},
		MaxRedirects:   8,
		ReadOnly:       false,
		RouteByLatency: true,
		RouteRandomly:  false,
	}

	currentClient, err := New(RedisCluster, config)

	currentClient.HMSet("secondcheckhasher", map[string]string{"fifthvalue": "fourthvalue"}, 10)
	result, err := currentClient.HMGet("secondcheckhasher", "fifthvalue")

	if len(result) != 1 {
		fmt.Print("Result not expected : ", err)
	}
}

func ExampleRedigoCluster_SetNX2() {
	config := ConfigCacheCluster{
		Addrs: []string{"127.0.0.1:6001", "127.0.0.1:6002", "127.0.0.1:6003",
			"127.0.0.1:6005", "127.0.0.1:6007"},
		MaxRedirects:   8,
		ReadOnly:       false,
		RouteByLatency: true,
		RouteRandomly:  false,
	}

	currentClient, err := New(RedisCluster, config)

	currentClient.Set("vivaket", "thevalue", 5)

	err = currentClient.SetNX("vivaket", "check", 0)

	if err == nil {
		fmt.Print("result not expected : ", err)
	}
}

func ExampleRedigoCluster_ScanKeys() {
	config := ConfigCacheCluster{
		Addrs: []string{"127.0.0.1:6001", "127.0.0.1:6002", "127.0.0.1:6003",
			"127.0.0.1:6005", "127.0.0.1:6007"},
		MaxRedirects:   8,
		ReadOnly:       false,
		RouteByLatency: true,
		RouteRandomly:  false,
	}

	currentClient, _ := New(RedisCluster, config)

	currentClient.Set("mikasa", "warrior", 60)

	result, _ := currentClient.ScanKeys("key*")

	if len(result) != 1 {
		fmt.Print("Result not expected : ", len(result))
	}
}

func ExampleRedigoCluster_SDiff() {
	config := ConfigCacheCluster{
		Addrs: []string{"127.0.0.1:6001", "127.0.0.1:6002", "127.0.0.1:6003",
			"127.0.0.1:6005", "127.0.0.1:6007"},
		MaxRedirects:   8,
		ReadOnly:       false,
		RouteByLatency: true,
		RouteRandomly:  false,
	}

	currentClient, _ := New(RedisCluster, config)

	result, _ := currentClient.SDiff("{key}10", "{key}15")

	if result[0] != "b" {
		fmt.Print("Result not expected : ", result[0])
	}
}

func ExampleRedigoCluster_SMove() {
	config := ConfigCacheCluster{
		Addrs: []string{"127.0.0.1:6001", "127.0.0.1:6002", "127.0.0.1:6003",
			"127.0.0.1:6005", "127.0.0.1:6007"},
		MaxRedirects:   8,
		ReadOnly:       false,
		RouteByLatency: true,
		RouteRandomly:  false,
	}

	currentClient, _ := New(RedisCluster, config)

	currentClient.SAdd("{ka}1", "a")
	currentClient.SAdd("{ka}2", "a")

	result, _ := currentClient.SMove("a", "{ka}1", "{ka}2")

	if result != 1 {
		fmt.Print("Result not expected : ", result)
	}
}

func ExampleRedigoCluster_Expire() {
	config := ConfigCacheCluster{
		Addrs: []string{"127.0.0.1:6001", "127.0.0.1:6002", "127.0.0.1:6003",
			"127.0.0.1:6005", "127.0.0.1:6007"},
		MaxRedirects:   8,
		ReadOnly:       false,
		RouteByLatency: true,
		RouteRandomly:  false,
	}

	currentClient, _ := New(RedisCluster, config)

	result, _ := currentClient.Expire("{ka}2", 60)

	if result != 1 {
		fmt.Print("Result not expected : ", result)
	}
}

func ExampleRedigoCluster_IncrXX() {
	config := ConfigCacheCluster{
		Addrs: []string{"127.0.0.1:6001", "127.0.0.1:6002", "127.0.0.1:6003",
			"127.0.0.1:6005", "127.0.0.1:6007"},
		MaxRedirects:   8,
		ReadOnly:       false,
		RouteByLatency: true,
		RouteRandomly:  false,
	}

	currentClient, _ := New(RedisCluster, config)
	currentClient.Set("abc", "1", 60)
	currentClient.IncrBy("abc", 1)

	result, _ := currentClient.Get("abc")

	if result != "2" {
		fmt.Print("Result not expected : ", result)
	}
}

func ExampleRedigoCluster_DecrWithLimit() {
	config := ConfigCacheCluster{
		Addrs: []string{"127.0.0.1:6001", "127.0.0.1:6002", "127.0.0.1:6003",
			"127.0.0.1:6005", "127.0.0.1:6007"},
		MaxRedirects:   8,
		ReadOnly:       false,
		RouteByLatency: true,
		RouteRandomly:  false,
	}

	currentClient, _ := New(RedisCluster, config)
	currentClient.Set("justtemp", "10", 60)
	result, _ := currentClient.DecrWithLimit("justtemp", 1, 1)

	if result != 9 {
		fmt.Print("Result not expected")
	}

}

func ExampleCache_GeoAdd() {
	conf := &Config{
		ServerAddr:      "127.0.0.1:6379",
		MaxIdle:         50,
		MaxActive:       400,
		IdleTimeout:     time.Duration(300) * time.Second,
		MaxConnLifetime: time.Duration(300) * time.Second,
	}
	c, _ := New(Redis, conf)
	geoPoints := []*GeoPoint{
		{
			Member:    "jakarta selatan",
			Longitude: 106.802079,
			Latitude:  -6.284106,
		},
		{
			Member:    "jakarta pusat",
			Longitude: 106.836615,
			Latitude:  -6.182313,
		},
	}
	_, e := c.GeoAdd("geo-info", geoPoints...)
	if e != nil {
		fmt.Printf("error not expected")
	}
}

func ExampleCache_GeoHash() {
	conf := &Config{
		ServerAddr:      "127.0.0.1:6379",
		MaxIdle:         50,
		MaxActive:       400,
		IdleTimeout:     time.Duration(300) * time.Second,
		MaxConnLifetime: time.Duration(300) * time.Second,
	}
	c, _ := New(Redis, conf)

	_, e := c.GeoHash("geo-info", []string{"jakarta selatan", "jakarta pusat"}...)
	if e != nil {
		fmt.Printf("error not expected")
	}
}

func ExampleCache_GeoRadius() {
	conf := &Config{
		ServerAddr:      "127.0.0.1:6379",
		MaxIdle:         50,
		MaxActive:       400,
		IdleTimeout:     time.Duration(300) * time.Second,
		MaxConnLifetime: time.Duration(300) * time.Second,
	}
	c, _ := New(Redis, conf)

	georadiusQuery := &GeoRadiusQuery{
		Radius:      6000,
		Unit:        GeoRadiusKiloMeter,
		WithCoord:   false,
		WithDist:    true,
		WithGeoHash: true,
		Count:       3,
		Sort:        GeoRadiusDesc,
	}
	geoLoc, e := c.GeoRadius("geo-info", 96, 2.58, georadiusQuery)
	if e != nil {
		fmt.Printf("error not expected")
	}
	for _, g := range geoLoc {
		fmt.Println(g.Name, g.Latitude, g.Longitude, g.Distance, g.GeoHash)
	}
}
