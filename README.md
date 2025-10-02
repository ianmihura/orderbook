# Orderbook

A dummy orderbook architecture, that simulates a queue of pending transactions of an asset at different prices. Orberbooks are common in public market excahges.

### How to use

Once you launch the project (for example `go run .`), the orderbook gets filled with random orders. There will always be some number of Market Makers adding random orders to the orderbook, keeping the asset price moving. Type `display` or `d` in the console to see the orderbook move in real time (press `c` to close it)

You can also submit transactions to the market, by typing `new` or `n`. You can check your portfolio balance typing `portfolio` or `p`.

Type `help` to see how else to interact with the market.
