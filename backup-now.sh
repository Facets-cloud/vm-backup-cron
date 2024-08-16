#!/bin/bash
# Fetch environment variables

sleeptime=1m # Sleep for 10 minutes after a failed try.
maxtries=5    # 5 * 10 minutes = about 50 minutes total of waiting,
              # not counting running and failing.

echo "Starting Victoria metrics backup process..."
while ! /vmbackup-prod -storageDataPath="${STORAGE_DATA_PATH}" -dst="${DESTINATION}" -snapshot.createURL="${SNAPSHOT_CREATE_URL}/snapshot/create"; do
  echo "Backup attempt failed. Remaining tries: $maxtries"
  maxtries=$(( maxtries - 1 ))
  if [ "$maxtries" -eq 0 ]; then
    echo "Victoria metrics backup didn't succeed after multiple attempts! Exiting." >&2
    exit 1
  fi

  echo "Sleeping for $sleeptime before next attempt..."
  sleep "$sleeptime" || break
done
echo "Victoria metrics backup process completed."
