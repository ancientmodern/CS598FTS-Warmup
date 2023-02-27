#!/bin/bash

# ./scripts/client_test.sh -i 1000000 -n 1 -r 100000 -w 0
# ./scripts/client_test.sh -i 1000000 -n 1 -r 50000 -w 50000
# ./scripts/client_test.sh -i 1000000 -n 1 -r 0 -w 100000

# ./scripts/client_test.sh -i 1000000 -n 2 -r 50000 -w 0
# ./scripts/client_test.sh -i 1000000 -n 2 -r 25000 -w 25000
# ./scripts/client_test.sh -i 1000000 -n 2 -r 0 -w 50000

# ./scripts/client_test.sh -i 1000000 -n 4 -r 25000 -w 0
# ./scripts/client_test.sh -i 1000000 -n 4 -r 12500 -w 12500
# ./scripts/client_test.sh -i 1000000 -n 4 -r 0 -w 25000

# ./scripts/client_test.sh -i 1000000 -n 8 -r 12500 -w 0
# ./scripts/client_test.sh -i 1000000 -n 8 -r 6250 -w 6250
# ./scripts/client_test.sh -i 1000000 -n 8 -r 0 -w 12500

# ./scripts/client_test.sh -i 1000000 -n 16 -r 6250 -w 0
# ./scripts/client_test.sh -i 1000000 -n 16 -r 3125 -w 3125
# ./scripts/client_test.sh -i 1000000 -n 16 -r 0 -w 6250

# ./scripts/client_test.sh -i 1000000 -n 32 -r 3125 -w 0
# ./scripts/client_test.sh -i 1000000 -n 32 -r 1562 -w 1563
# ./scripts/client_test.sh -i 1000000 -n 32 -r 0 -w 3125


# test for correctness
./scripts/client_test.sh -i 1000 -n 32 -r 10000 -w 10000 -c true

