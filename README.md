## stockticker

A simple CLI application to watch a given set of stocks.

### Build & Install

```sh
$ make install
```

### Usage

```sh
$ stockticker --help
```

### Examples

```sh
$ stockwatcher -s GOOG,IBM,YHOO,CSCO,AAPL,FB,TWTR -i 5
```

Output from the above command shown below.

```sh

   AAPL  111.09  111.25    ↓
   CSCO   27.69   27.68    ↑
     FB   77.37   77.41    ↓
   GOOG  528.57       -    -
    IBM  154.57  154.46    ↑
   TWTR   38.77    38.8    ↓
   YHOO   48.66       -    -

```
