## stockticker

A simple CLI application to watch a given set of stocks.

### Build & Install

```sh
$ make install
```

### Usage

```sh
$ stockwatcher --help
Usage of stockwatcher:
  -i=1: Interval for stock data to be updated in seconds
  -s="": Symbols for ticker, comma seperate (no spaces)
```

### Example

Displayed as Symbol, Current Price, Previous Price, and up/down indicator.

```sh
$ stockwatcher -s GOOG,IBM,YHOO,CSCO,AAPL,FB,TWTR -i 1
```
```sh

   AAPL  111.27  111.26    ↑
   CSCO   27.83   27.86    ↓
     FB   77.37       -    -
   GOOG  529.03     529    ↑
    IBM  153.52   153.5    ↑
   TWTR   39.04   39.02    ↑
   YHOO   48.41   48.44    ↓

```

![http://i.picasion.com/pic79/0131736efe91ae428c0fe8f60fc92f3c.gif](http://i.picasion.com/pic79/0131736efe91ae428c0fe8f60fc92f3c.gif)
