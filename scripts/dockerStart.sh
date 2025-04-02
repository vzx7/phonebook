#!/usr/bin/sh
sudo rm -rf /home/$USER/REPS_LEARN/phonebook/db/postgres/*
systemctl start docker.service
cd ../db/
docker-compose up
exit 0