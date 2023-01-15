#!/bin/sh

migrate create -ext sql -dir migrations -seq $1
