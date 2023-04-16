# CS 598 FTS Warmup Project

## Authors

{Haorong Lu, Lan Zhang, Jiaai Xu} @ UIUC MCS

## Design and Implementation

Using the Multiple-Write-Multiple-Read (MWMR) shared register protocol as a basis, we developed two main programs in Golang: the client and the replica. The client program includes four important functions, namely GetPhase and SetPhase for the Writer and Reader, respectively, which send the corresponding RPC requests to the replica. We use gRPC as the communication protocol between clients and replicas. Since the GetPhase/SetPhase functions for both the Writer and Reader are almost identical in terms of RPC message format, we define only two gRPC methods: GetPhase and SetPhase. The replica's main responsibility is to act as a backend for these two gRPC methods and update its own local storage, which is represented by a global hash table (map in Go).

A challenging aspect of the clients was determining how to concurrently send gRPC requests to all replicas and wait for a majority of responses before proceeding. To accomplish this, we used a combination of goroutines, sync.WaitGroup, and channels. Below is a simplified version of our approach:

```go
var wg sync.WaitGroup
wg.Add(n)
ch := make(chan Pair, n)

for i := 0; i < n; i++ {
	go func(rid int) {
		// Send a gRPC request and put the response in `temp`
		// ...... 
		ch <- temp
		wg.Done()
	}(i)
}

go func() {
	wg.Wait()
}()

done := make([]Pair, 0, f+1)
for pair := range ch {
	done = append(done, pair)
	if len(done) >= f+1 {
		break
	}
}
```

Another important consideration is that we use a global map as the replica storage, which is not thread-safe in Go. Since multiple gRPC streams may access the same key simultaneously, we need to add a RWLock to the map. To minimize the locking overhead, we added a fine-grained lock to each key-value pair instead of using a global lock. This design improves the replica's concurrency performance, especially when handling a large number of clients. The remaining parts were relatively straightforward, involving the translation of MWMR provisions into corresponding code. If you have any further questions about the design and implementation, you can refer directly to our source code.

## Testing Correctness

Before any evaluation on the performance of our implementation, correctness needs to be promised first, which is tested as follows.

1. Start 5 servers on individual Bare Metal PCs with a data store of size 1000 (1000 key-value pairs). One replica is set to fail in the middle of the running. Bash Command: `./scripts/run_server.sh -s 1000`
2. Start 32 multiple clients concurrently on another cloud Bare Metal PC where each client is going to send 10000 reads and 10000 writes rpc requests to the distributed storage system. Then the test script will track every read and write event with its key, timestamp and value, then log it into a log file. Bash Command: `./scripts/client_test.sh -i 1000 -n 32 -r 10000 -w 10000 -c true` (`-c` means correctness mode, the program will skip useless event tracking operations. Here we choose 1000 as the data store size and 32 as the number of clients because they maximize the concurrency and thus maximizing the possibility of exposing the linearizability problem. )
3. Run a Python script called verify_correctness.py to verify the linearizability of each client’s log file. If the script prints “Linearizability Verified!”, then linearizability of all clients are verified and correctness is verified. Bash Command: `python3 data_processing/verify_correctness.py`

### Logistics of the Verification Script

Under testing correctness mode, each client opens its own log file with its id as the unique identifier under the folder “logs_for_correctness”. The log files format is as follows:

```go
INFO: 2023/02/26 17:46:03 main.go:159: ithwrite: 300 key: 474 timestamp: 6 cid: 1 value: 791
INFO: 2023/02/26 17:46:03 main.go:146: ithread: 301 key: 529 timestamp: 7 cid: 14 value: 166
INFO: 2023/02/26 17:46:03 main.go:146: ithread: 302 key: 327 timestamp: 4 cid: 19 value: 320
INFO: 2023/02/26 17:46:03 main.go:146: ithread: 303 key: 857 timestamp: 8 cid: 29 value: 176
INFO: 2023/02/26 17:46:03 main.go:146: ithread: 304 key: 929 timestamp: 9 cid: 9 value: 994
INFO: 2023/02/26 17:46:03 main.go:159: ithwrite: 304 key: 876 timestamp: 6 cid: 1 value: 106
```

where ‘key’, ‘timestamp’, ‘cid’ and the relative order of each log data are needed for verification. As for one single client, it will send a new rpc request only after the previous one is finished so it is promised that in one log file, the real-time order is the same as the order in the log file. Therefore, the only next step to verify the linearizability is to check that the timestamp is not decreasing for the same key. The Python script is shown below.

```python
import os

def isT1BeforeT2(t1: list, t2: list):
   if t1[0] < t2[0]: return True
   elif t1[0] > t2[0]: return False
   else:
       if t1[1] <= t2[1]: return True
       else: return False

directory = "logs_for_correctness/logs_with_32_clients"
for filename in os.listdir(directory):
   data = {}
   filepath = os.path.join(directory, filename)
   if not os.path.isfile(filepath):
       print("not a file {}".format(filepath))
       continue
   f = open(filepath, 'r')
   for line in f.readlines():
       line = line.split()
       if line[7] not in data:
           data[line[7]] = [int(line[9]), int(line[11])]
       else:
           if isT1BeforeT2(data[line[7]], [int(line[9]), int(line[11])]): data[line[7]] = [int(line[9]), int(line[11])]
           else:  
               print("NOT linear!")
               print(filename,  "line: ", line)
               f.close()
               exit(1)
   f.close()

print("Linearizability Verified!")
```

