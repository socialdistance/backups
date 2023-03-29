#!/bin/bash
echo "Starting the backup script..."
ROOTDIR="/home/user/backup/"
YEAR=`date +%Y`
MONTH=`date +%m`
DAY=`date +%d`
SERVER="localhost"
#BLACKLIST="information_schema performance_schema"
ADDITIONAL_MYSQLDUMP_PARAMS="--skip-lock-tables --skip-add-locks --quick --single-transaction"
MYSQL_USER="root"
MYSQL_PASSWORD="qwerty123"

DOCKERCONTAINERNAME="database_backup"

PG_USER="postgres"
PG_PASSWORD=""

IP=$(ip route get 8.8.8.8 | awk -F"src " 'NR==1{split($2,a," ");print a[1]}')

MYSQLTARGET="$ROOTDIR/backup-mysql-${YEAR}-${MONTH}-${DAY}.sql"
PSQLTARGET="$ROOTDIR/backup-postgresql-${YEAR}-${MONTH}-${DAY}.sql"
PSQLTARGETDOCKER="$ROOTDIR/backup-postgresql-docker-${YEAR}-${MONTH}-${DAY}.sql"


TARTARGET="$ROOTDIR/${IP}-backup-${YEAR}-${MONTH}-${DAY}.tar.gz"

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

echo "running psqldump"
#pg_dump -U $PG_USER $DATABASE > $PSQLTARGET
pg_dump > $PSQLTARGET

echo "running psqldump in docker"
docker exec -t $DOCKERCONTAINERNAME pg_dumpall -c -U postgres > $PSQLTARGETDOCKER

tar -czvf $TARTARGET -C $ROOTDIR .

echo
echo "dump completed"