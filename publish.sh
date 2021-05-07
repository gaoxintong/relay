./build.sh
ssh root@123.232.115.10 -p 9022 "/root/services/prg/kill.sh"
scp -P 9022 ./prg root@123.232.115.10:/root/services/prg
ssh root@123.232.115.10 -p 9022 "/root/services/prg.sh"
