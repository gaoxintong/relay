./build.sh
ssh root@123.232.115.10 -p 9022 "/root/services/gateway/kill.sh"
scp -P 9022 ./gateway root@123.232.115.10:/root/services/gateway
ssh root@123.232.115.10 -p 9022 "/root/services/gateway/run.sh"
