#!/bin/bash
echo "Starting the backup script..."
ROOTDIR="/home/user/backup/"
YEAR=`date +%Y`
MONTH=`date +%m`
DAY=`date +%d`
HOUR=`date +%H`
SERVER="localhost"
#BLACKLIST="information_schema performance_schema"
ADDITIONAL_MYSQLDUMP_PARAMS="--skip-lock-tables --skip-add-locks --quick --single-transaction"
MYSQL_USER="root"
MYSQL_PASSWORD="qwerty123"

MYSQLTARGET="$ROOTDIR/backup-mysql-${YEAR}-${MONTH}-${DAY}.sql"
TARTARGET="$ROOTDIR/backup-${YEAR}-${MONTH}-${DAY}.tar.gz"


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

echo "running mysqldump"
mysqldump ${ADDITIONAL_MYSQLDUMP_PARAMS} -u ${MYSQL_USER} -p${MYSQL_PASSWORD} -h $SERVER $db --all-databases > $MYSQLTARGET
tar -czvf $TARTARGET $ROOTDIR

echo
echo "dump completed"