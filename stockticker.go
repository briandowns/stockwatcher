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
	"os/exec"
	"os/signal"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

const TIMEOUT = time.Duration(time.Second * 10)
const URL = "http://finance.yahoo.com/webservice/v1/symbols/%s/quote?format=json"

var UP = 8593   // '↑'
var DOWN = 8595 //'↓'

var re = regexp.MustCompile(`^\d.+\.\d{2}`) // this is to have only 2 decimal places
var signalChan = make(chan os.Signal, 1)    // channel to catch ctrl-c

var (
	symbolFlag   = flag.String("s", "", "Symbols for ticker, comma seperate (no spaces)")
	intervalFlag = flag.Int("i", 0, "Interval for stock data to be updated in seconds")
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
	quotes   map[string]map[string]float64
	interval time.Duration
	m        *sync.Mutex
}

// clearScreen runs a shell clear command
func clearScreen() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
}

func NewStockTicker(i time.Duration) *stockticker {
	return &stockticker{
		quotes:   make(map[string]map[string]float64),
		interval: i,
		m:        &sync.Mutex{},
	}
}

func (t *stockticker) add(symbol string) {
	t.m.Lock()
	defer t.m.Unlock()
	if _, ok := t.quotes[symbol]; !ok {
		t.quotes[symbol] = map[string]float64{}
	}

}

func (t *stockticker) updateStock(symbol string, price float64) {
	t.m.Lock()
	defer t.m.Unlock()
	t.quotes[symbol] = map[string]float64{
		"previous": t.quotes[symbol]["current"],
		"current":  price,
	}
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

func (t *stockticker) runner() {
	var wg sync.WaitGroup
	for k, _ := range t.quotes {
		wg.Add(1)
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
}

func (t *stockticker) printData() {
	var keys []string
	for k := range t.quotes {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	pos := 1
	for _, k := range keys {
		if t.quotes[k]["previous"] == 0.00 || t.quotes[k]["previous"] == t.quotes[k]["current"] {
			printTb(1, pos, fmt.Sprintf("%6s %7v %7s %4s\n", k, t.quotes[k]["current"], "-", "-"))
			pos++
			//fmt.Printf("%6s %7v %%%5s %4s\n", k, v["current"], "-", "-")
		} else if t.quotes[k]["current"] > t.quotes[k]["previous"] {
			printTb(1, pos, fmt.Sprintf("%6s %7v %7v %4s\n", k, t.quotes[k]["current"], t.quotes[k]["previous"], UP))
			//fmt.Printf("%6s %7v %%%5v %4s\n", k, v["current"], 100*(v["previous"]/v["current"]), green(UP))
			pos++
		} else {
			printTb(1, pos, fmt.Sprintf("%6s %7v %7v %4s\n", k, t.quotes[k]["current"], t.quotes[k]["previous"], DOWN))
			//fmt.Printf("%6s %7v %%%5v %4s\n", k, v["current"], 100*(v["previous"]/v["current"]), red(DOWN))
			pos++
		}
	}
}

func printTb(x, y int, msg string) {
	for _, c := range []rune(msg) {
		termbox.SetCell(x, y, c, termbox.ColorWhite, termbox.ColorDefault)
		x += runewidth.RuneWidth(c)
	}
	termbox.Flush()
}

func main() {
	flag.Parse()

	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for range signalChan {
			os.Exit(1)
		}
	}()

	t := NewStockTicker(time.Duration(*intervalFlag) * time.Second)

	switch {
	case strings.Contains(*symbolFlag, ","):
		for _, a := range strings.Split(*symbolFlag, ",") {
			t.add(a)
		}
	default:
		t.add(*symbolFlag)
	}

	err := termbox.Init()
	if err != nil {
		log.Fatal(err)
	}
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	event := make(chan termbox.Event)
	go func() {
		for {
			// Post events to channel
			event <- termbox.PollEvent()
		}
	}()

loop:
	for {
		t.runner()
		t.printData()

		// Poll key event or timeout (maybe)
		select {
		case <-event:
			break loop
			return
		case <-time.After(t.interval):
			continue loop
		}
	}
	close(event)
	time.Sleep(1 * time.Second)
	termbox.Close()
	os.Exit(0)
}
