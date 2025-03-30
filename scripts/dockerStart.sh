#!/usr/bin/sh
systemctl start docker.service
cd ../db/
docker-compose up
exit 0