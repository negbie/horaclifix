#!/bin/bash

echo "Updating Horaclifix... "
rm -rf /go/src/github.com/negbie/horaclifix
go get github.com/negbie/horaclifix
go install github.com/negbie/horaclifix

echo "HORACLIFIX GO Starting..."

COMMAND="horaclifix "

if [[ ${DEBUG:+1} ]] ; then
  echo "DEBUG Enabled..."
  COMMAND+="-d "
fi

if [[ ${VERBOSE:+1} ]] ; then
  echo "VERBOSE Enabled..."
  COMMAND+="-v "
fi

if [[ $HEP_HOST =~ ^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  echo "Adding HEP target: $HEP_HOST:$HEP_PORT"
  COMMAND+="-H $HEP_HOST:$HEP_PORT -HQ "
else
  HEP_IP=$(ping -q -c 1 $HEP_HOST | grep PING | sed -e "s/).*//" | sed -e "s/.*(//")
  if [[ ${HEP_IP:+1} ]] ; then
      echo "Adding HEP target: $HEP_IP:$HEP_PORT"
      COMMAND+="-H $HEP_IP:$HEP_PORT -HQ "
  else
      echo "Failed resolving $HEP_HOST ! Exiting..."
      exit 1;
  fi
fi


if [[ ${GRAYLOG_URL:+1} ]] ; then
  echo "Adding Graylog target: $GRAYLOG_URL"
  COMMAND+="-g $GRAYLOG_URL "
fi


if [[ ${STATSD_URL:+1} ]] ; then
  echo "Adding Statsd target: $STATSD_URL"
  COMMAND+="-s $STATSD_URL "
fi

echo "Executing: $COMMAND "
/bin/bash -c "$COMMAND"
