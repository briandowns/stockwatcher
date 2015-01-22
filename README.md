## stockticker

A simple CLI application to monitor a given set of stocks.

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
$ stockticker -s GOOG,IBM,YHOO,CSCO,AAPL,FB,TWTR -i 5
```

Output from the above command shown below.

```sh

   AAPL  109.55       -    -
   CSCO   27.84       -    -
     FB   76.73       -    -
   GOOG  518.03       -    -
    IBM  152.08       -    -
   TWTR   37.83       -    -
   YHOO   48.18       -    -

```
