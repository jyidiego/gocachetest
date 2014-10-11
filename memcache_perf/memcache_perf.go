package memcache_perf

import (
	"errors"
	"flag"
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"log"
	"strings"
	"time"
)

func GetFromMemcache(mc *memcache.Client, key string) error {
	defer un(trace("Starting GetFromMemcache"))
	it, err := mc.Get(key)
	if err != nil {
		fmt.Printf("Error: %s", err)
	} else {
		fmt.Printf("Retrieving Key %s -> %s\n", it.Key, it.Value)
	}
	return nil
}

func SetToMemcache(mc *memcache.Client, key string, value []byte) error {
	mc.Set(&memcache.Item{Key: key, Value: value})
	return nil
}

func trace(s string) (string, time.Time) {
	log.Println("START:", s)
	return s, time.Now()
}

func un(s string, startTime time.Time) {
	endTime := time.Now()
	log.Println("  END:", s, "ElapsedTime in seconds:", endTime.Sub(startTime))
}

type ServerList []string

func (s *ServerList) String() string {
	return fmt.Sprint(*s)
}

func (s *ServerList) Set(value string) error {
	if len(*s) > 0 {
		return errors.New("memcache flag already set")
	}
	for _, server := range strings.Split(value, ",") {
		*s = append(*s, server)
	}
	return nil
}

var serverList ServerList
var key string
var value string

func init() {
	flag.Var(&serverList, "memcache", "comma-seperated list of servers")
	flag.StringVar(&key, "key", "TestKey", "A string for the key")
	flag.StringVar(&value, "value", "The Cowboys are 4-1 and loving it", "A string for the value")
}

func main() {
	flag.Parse()
	if serverList == nil {
		serverList = append(serverList, "localhost:11211")
	}
	fmt.Printf("Server list: %s", serverList)
	mc := memcache.New(strings.Join(serverList, ","))
	SetToMemcache(mc, "TestKey", []byte("Cowboys are 4-1 and we are loving it!"))
	GetFromMemcache(mc, "TestKey")
}
