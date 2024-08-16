FROM victoriametrics/vmbackup:v1.96.0 AS vmbackup

RUN mkdir /job
RUN apk add --no-cache bash
ADD backup-now.sh /job/backup-now.sh
ADD entry.sh /entry.sh
RUN chmod +x /job/backup-now.sh /entry.sh

RUN crontab -l | { cat; echo "0 * * * * bash /job/backup-now.sh 2>&1"; } | crontab -

ENTRYPOINT ["/entry.sh"]