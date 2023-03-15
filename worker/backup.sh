#!/bin/bash
echo "Starting the backup script..."
ROOTDIR="/home/user/backup/mysql/"
YEAR=`date +%Y`
MONTH=`date +%m`
DAY=`date +%d`
HOUR=`date +%H`
SERVER="localhost"
#BLACKLIST="information_schema performance_schema"
ADDITIONAL_MYSQLDUMP_PARAMS="--skip-lock-tables --skip-add-locks --quick --single-transaction"
MYSQL_USER="root"
MYSQL_PASSWORD="qwerty123"

# Read MySQL password from stdin if empty
if [ -z "${MYSQL_PASSWORD}" ]; then
 echo -n "Enter MySQL ${MYSQL_USER} password: "
 read -s MYSQL_PASSWORD
 echo
fi

# Check MySQL credentials
echo exit | mysql --user=${MYSQL_USER} --password=${MYSQL_PASSWORD} --host=${SERVER} -B 2>/dev/null
if [ "$?" -gt 0 ]; then
 echo "MySQL ${MYSQL_USER} - wrong credentials"
 exit 1
else
 echo "MySQL ${MYSQL_USER} - was able to connect."
fi

#creating backup path
if [ ! -d "$ROOTDIR/$YEAR/$MONTH/$DAY" ]; then
    mkdir -p "$ROOTDIR/$YEAR/$MONTH/$DAY"
    chmod -R 700 $ROOTDIR
fi

echo "running mysqldump"
dblist=`mysql -u ${MYSQL_USER} -p${MYSQL_PASSWORD} -h $SERVER -e "show databases" | sed -n '2,$ p'`
for db in $dblist; do
    echo "Backuping $db"
    mysqldump ${ADDITIONAL_MYSQLDUMP_PARAMS} -u ${MYSQL_USER} -p${MYSQL_PASSWORD} -h $SERVER $db | gzip --best > "$ROOTDIR/$YEAR/$MONTH/$DAY/`echo $db | sed 's/\//_/g'`.sql.gz"
            echo "Backup of $db ends with $? exit code"
    # isBl=`echo $BLACKLIST |grep $db`
#    if [ $? == 1 ]; then
#        mysqldump ${ADDITIONAL_MYSQLDUMP_PARAMS} -u ${MYSQL_USER} -p${MYSQL_PASSWORD} -h $SERVER $db | gzip --best > "$ROOTDIR/$YEAR/$MONTH/$DAY/`echo $db | sed 's/\//_/g'`.sql.gz"
#        echo "Backup of $db ends with $? exit code"
#    else
#        echo "Database $db is blacklisted, skipped"
    fi
done
echo
echo "dump completed"