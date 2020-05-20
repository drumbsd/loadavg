package main

import (
	"flag"
	"fmt"
	"github.com/go-redis/redis"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	var i int = 288
	var future = time.Date(2000, 2, 1, 23, 55, 0, 0, time.UTC)
	var s []time.Time
	var vettore []string
	redisPtr := flag.String("redis", "", "Redis server")
	flag.Parse()
	if *redisPtr == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	f, err := os.Create("loadavg.html")
	if err != nil {
		fmt.Println(err)
		return
	}

	for k := 0; k < i; k++ {
		future = future.Add(5 * time.Minute)
		s = append(s, future)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     *redisPtr + ":" + "6379", // use default Addr
		Password: "",                       // no password set
		DB:       1,                        // use default DB
	})
	err = rdb.Ping().Err()
	if err != nil {
		panic(err)
	}

	rdb2 := redis.NewClient(&redis.Options{
		Addr:     *redisPtr + ":" + "6379", // use default Addr
		Password: "",                       // no password set
		DB:       2,                        // use default DB
	})
	err = rdb2.Ping().Err()
	if err != nil {
		panic(err)
	}

	vettore, err = rdb2.Keys("*").Result()
	// Creazione prima riga tabella con la lista dei server
	_, _ = f.WriteString(`
	   <style type="text/css">
.tg  {border-collapse:collapse;border-spacing:0;}
.tg td{font-family:Arial, sans-serif;font-size:14px;padding:10px 5px;border-style:solid;border-width:1px;overflow:hidden;word-break:normal;border-color:black;}
.tg th{background:#eee;font-family:Arial, sans-serif;font-size:14px;font-weight:normal;padding:10px 5px;border-style:solid;border-width:1px;overflow:hidden;word-break:normal;border-color:black;}
.tg .tg-c3ow{border-color:inherit;text-align:center;vertical-align:top}
.tg .tg-0pky{border-color:inherit;text-align:left;vertical-align:top}

.tableFixHead          { overflow-y: auto; height: 100%; }
.tableFixHead thead th { position: sticky; top: 0; }

/* Just common table stuff. Really. */
table  { border-collapse: collapse; width: 100%;}
th, td { padding: 8px 16px; }
</style>
<div class="tableFixHead">
<table align=center border=1 class="tg"><thead><tr><th></th>

	   `)

	for _, v := range vettore {
		v2 := strings.Split(v,".")
		linea := fmt.Sprintln("<th class=tg-0pky>", v2[0], "</th>")
		_, _ = f.WriteString(linea)
	}
	//fmt.Println("</tr><tr>")
	_, _ = f.WriteString("</thead><tbody></tr><tr>")
	for _, v := range s {
		linea := fmt.Sprintln("<td class=tg-0pky>", v.Format("15:04:00"), "</td>")
		_, _ = f.WriteString(linea)
		for _, a := range vettore {
			chiave := v.Format("15:04:00") + "_" + a
			val, err := rdb.Get(chiave).Result()
			if err != nil {
				//panic(err)
				val = "-1"
			}
			//Cambiamo il colore dello sfondo in base al valore.

			i, err := strconv.Atoi(val)
			if i == 0 {
				linea = fmt.Sprintln("<td class=tg-0pky align=center>", i, "</td>")
			} else if i >= 1 {
				linea = fmt.Sprintln("<td class=tg-0pky align=center bgcolor=red>", i, "</td>")
			} else {
				linea = fmt.Sprintln("<td class=tg-0pky align=center bgcolor=green>", i, "</td>")
			}

			_, _ = f.WriteString(linea)
		}
		// fmt.Println("</tr>")
		_, _ = f.WriteString("</tr>")
	}
	//fmt.Println("</tr></table>
	_, _ = f.WriteString("</tr></tbody></table></div>")
	rdb.Close()
	rdb2.Close()
	f.Close()
}
