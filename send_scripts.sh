scp -r output/ run_server.sh client_test.sh Lanz@ms1040.utah.cloudlab.us:~ &
scp -r output/ run_server.sh client_test.sh Lanz@ms1016.utah.cloudlab.us:~ &
scp -r output/ run_server.sh client_test.sh Lanz@ms1005.utah.cloudlab.us:~ &
scp -r output/ run_server.sh client_test.sh Lanz@ms1022.utah.cloudlab.us:~ &
scp -r output/ run_server.sh client_test.sh Lanz@ms1045.utah.cloudlab.us:~ &
scp -r output/ run_server.sh client_test.sh Lanz@ms1006.utah.cloudlab.us:~ &


wait
echo "finish"