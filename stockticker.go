// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

const TIMEOUT = time.Duration(time.Second * 10)
const URL = "http://finance.yahoo.com/webservice/v1/symbols/%s/quote?format=json"

var re = regexp.MustCompile(`^\d.+\.\d{2}`)
var signalChan = make(chan os.Signal, 1) // channel to catch ctrl-c

var (
	symbolFlag   = flag.String("s", "", "Symbols for ticker, comma seperate (no spaces)")
	rateFlag     = flag.Int("r", 1, "Speed of stock data")
	intervalFlag = flag.Int("i", 5, "Interval for stock data to be updated")
)

type Stock struct {
	List List `json:"list"`
}

type List struct {
	Meta      Meta        `json:"meta"`
	Resources []Resources `json:"resources"`
}

type Meta struct {
	Type  string `json:"type"`
	Start uint   `json:"start"`
	Count uint   `json:"count"`
}

type Resources struct {
	Resource Resource `json:"resource"`
}

type Resource struct {
	Classname string `json:"classname"`
	Fields    Fields `json:"fields"`
}

type Fields struct {
	Name    string `json:"name"`
	Price   string `json:"price"`
	Symbol  string `json:"symbol"`
	TS      string `json:"ts"`
	Type    string `json:"type"`
	UTCTime string `json:"utctime"`
	Volume  string `json:"volume"`
}

type stockticker struct {
	symbolData map[string]float64
	speed      int
	interval   time.Duration
	m          *sync.Mutex
}

func NewStockTicker(s int, i time.Duration) *stockticker {
	return &stockticker{
		symbolData: make(map[string]float64),
		speed:      s,
		interval:   i,
		m:          &sync.Mutex{},
	}
}

func (t *stockticker) add(symbol string) {
	t.m.Lock()
	defer t.m.Unlock()
	if _, ok := t.symbolData[symbol]; !ok {
		t.symbolData[symbol] = 0.0
	}
}

func (t *stockticker) updateStock(symbol string, price float64) {
	t.m.Lock()
	defer t.m.Unlock()
	fmt.Println("updating for ", symbol)
	t.symbolData[symbol] = price
}

func query(symbol string) (*Stock, error) {
	data := &Stock{}
	client := http.Client{
		Timeout: TIMEOUT,
	}
	resp, err := client.Get(fmt.Sprintf(URL, symbol))
	if err != nil {
		return nil, errors.New("unable to retrive symbol data")
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
	return data, nil
}

func convertPrice(p string) float64 {
	price, err := strconv.ParseFloat(p, 64)
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
	return price
}

func main() {
	flag.Parse()
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for range signalChan {
			os.Exit(1)
		}
	}()
	t := NewStockTicker(*rateFlag, time.Duration(*intervalFlag)*time.Minute)
	switch {
	case strings.Contains(*symbolFlag, ","):
		for _, a := range strings.Split(*symbolFlag, ",") {
			fmt.Println("Adding " + a)
			t.add(a)
		}
	default:
		t.add(*symbolFlag)
	}
	var wg sync.WaitGroup
	for k, _ := range t.symbolData {
		wg.Add(1)
		fmt.Println("Getting for " + k)
		go func(k string) {
			defer wg.Done()
			stock, err := query(k)
			if err != nil {
				log.Fatalln(err)
				os.Exit(1)
			}
			t.updateStock(stock.List.Resources[0].Resource.Fields.Symbol,
				convertPrice(re.FindString(stock.List.Resources[0].Resource.Fields.Price)),
			)
		}(k)
	}
	wg.Wait()
	fmt.Println(t)
	os.Exit(0)
}
