#!/bin/bash
cat /var/fixtures/data.sql | /usr/bin/mysql -u user --password=password db
