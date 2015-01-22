## stockticker

A simple CLI application to watch a given set of stocks.

### Build & Install

```sh
$ make install
```

### Usage

```sh
$ stockwatcher --help
```

### Examples

```sh
$ stockwatcher -s GOOG,IBM,YHOO,CSCO,AAPL,FB,TWTR -i 5
```

Output from the above command shown below.

Displayed as Symbol, Current Price, Previous Price, and up/down indicator.

```sh

   AAPL  111.27  111.26    ↑
   CSCO   27.83   27.86    ↓
     FB   77.37       -    -
   GOOG  529.03     529    ↑
    IBM  153.52   153.5    ↑
   TWTR   39.04   39.02    ↑
   YHOO   48.41   48.44    ↓

```
