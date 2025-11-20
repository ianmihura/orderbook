# Orderbook

A dummy orderbook architecture. Simulates a queue of pending transactions at different prices. When a match is found between two prices, the matching engine will execute the trade.

Orberbooks are common in public market excahges.

### How to use

Once you launch the project (for example `go run .`), the orderbook gets filled with random orders. There will always be some number of automatic traders (Market Makers, institutionals, nosie traders) adding orders to the orderbook, keeping the asset price moving.

Type `display` or `d` in the console to see the orderbook move in real time (press `c` to close it).

You can also submit orders to the market, by typing `new` or `n`, and following the wizard. You can check your portfolio balance typing `portfolio` or `p`.

Type `help` to see how else to interact with the market.
