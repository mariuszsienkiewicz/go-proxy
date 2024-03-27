#!/bin/bash

mysql -u root -p$MYSQL_ROOT_PASSWORD -e"CREATE USER '$GOPROXY_REPLICATION_USER'@'%' IDENTIFIED BY '$GOPROXY_REPLICATION_PASSWORD';"
mysql -u root -p$MYSQL_ROOT_PASSWORD -e"GRANT REPLICATION SLAVE ON *.* TO '$GOPROXY_REPLICATION_USER'@'%';"
mysql -u root -p$MYSQL_ROOT_PASSWORD -e"FLUSH PRIVILEGES;"
mysql -u root -p$MYSQL_ROOT_PASSWORD -e"SHOW MASTER STATUS\G;"