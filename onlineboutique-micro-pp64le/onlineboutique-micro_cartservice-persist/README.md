# CARTSERVICE-PERSIST

## INFO

The original _cartservice_ (https://github.com/go-micro/demo/tree/main/src/cartservice) requires an underlying Redis pod, however **it still saves data only on memory**.

I refactored the code, so that it actually stores data on the database (also using Valkey instead of Redis)

## NOTES

- **THE REFACTORED CODE / CODE ARCHITECTURE ARE UGLY**, but still enough for a functioning demo

- The code / Dockerfile here are for an x86 architecture. To "convert" them to ppc64le, follow the same instrucitons as per the other microservices of this repo

## TEST

To check that data is actually saved on the db:

- `oc port-forward svc/redis-cart 6379:6379`
- `valkey-cli keys \* &&  valkey-cli keys \* | xargs -I {} redis-cli hscan {} 0`