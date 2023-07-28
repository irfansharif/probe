
    go test . -v -run 'TestReadBandwidth'   -p 1 -count 1
    go test . -v -run 'TestReadIOPS'        -p 1 -count 1
    go test . -v -run 'TestWriteBandwidth'  -p 1 -count 1
    go test . -v -run 'TestWriteIOPS'       -p 1 -count 1


    === RUN   TestSupported
        probe_test.go:30: supported = true
    --- PASS: TestSupported (0.00s)

    === RUN   TestWriteIOPS/with-max=true
        probe_test.go:60: write iops (max = 1000) = 999

    === RUN   TestWriteIOPS/with-max=false
        probe_test.go:60: write iops = 21484

    === RUN   TestWriteBandwidth/with-max=true
        probe_test.go:82: write bandwidth (max = 80 MiB/s) = 81 MiB/s

    === RUN   TestWriteBandwidth/with-max=false
        probe_test.go:82: write bandwidth = 1.1 GiB/s

    === RUN   TestReadIOPS/with-max=true
        probe_test.go:104: read iops (max = 1000) = 999

    === RUN   TestReadIOPS/with-max=false
        probe_test.go:104: read iops = 71493

    === RUN   TestReadBandwidth/with-max=true
        probe_test.go:126: read bandwidth (max = 80 MiB/s) = 81 MiB/s

    === RUN   TestReadBandwidth/with-max=false
        probe_test.go:126: read bandwidth = 3.4 GiB/s
