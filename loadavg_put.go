package main

import (
	"flag"
	"github.com/go-redis/redis"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/load"
	"os"
	"time"
)

func main() {
	redisPtr := flag.String("redis", "", "Redis server")
	flag.Parse()
	if *redisPtr == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	var loadAverage int
	var secondi time.Duration = 300
	rdb := redis.NewClient(&redis.Options{
		Addr:     *redisPtr + ":" + "6379", // use default Addr
		Password: "",                       // no password set
		DB:       1,                        // use default DB
	})
	err := rdb.Ping().Err()
	if err != nil {
		panic(err)
	}

	// We use this second connection to dynamically populate list of servers

	rdb2 := redis.NewClient(&redis.Options{
		Addr:     *redisPtr + ":" + "6379", // use default Addr
		Password: "",                       // no password set
		DB:       2,                        // use default DB
	})
	err = rdb2.Ping().Err()
	if err != nil {
		panic(err)
	}

	for true {
		location, _ := time.LoadLocation("Europe/Rome")
		t := time.Now().In(location)
		d := time.Duration(5 * time.Minute)
		hostname, err := os.Hostname()
		if err != nil {
			panic(err)
		}
		hostname_time := t.Round(d).Format("15:04:05.999999999") + "_" + hostname

		load, _ := load.Avg()
		loadAverage = int(load.Load5)
		cpuNumber, err := cpu.Counts(true)
		loadAverage = loadAverage / cpuNumber
		err = rdb.Set(hostname_time, loadAverage, 86100*time.Second).Err()
		err = rdb2.Set(hostname, "1", 600*time.Second).Err()
		if err != nil {
			panic(err)
		}

		time.Sleep(secondi * time.Second)
	}
	rdb.Close()
	rdb2.Close()
}
