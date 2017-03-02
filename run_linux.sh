#! /bin/sh
./golang-ots-linux -inside=true -insert=true -num=100 -cols=1
./golang-ots-linux -inside=true -insert=true -num=100 -cols=100
./golang-ots-linux -inside=true -insert=true -num=100 -cols=1000
./golang-ots-linux -inside=true -insert=true -num=100 -cols=5000
./golang-ots-linux -inside=true -insert=true -num=100 -cols=10000
./golang-ots-linux -inside=true -insert=true -num=100 -cols=30000
./golang-ots-linux -inside=true -insert=true -num=100 -cols=50000

./golang-ots-linux -inside=true -get=true -num=100 -cols=1
./golang-ots-linux -inside=true -get=true -num=100 -cols=100
./golang-ots-linux -inside=true -get=true -num=100 -cols=1000
./golang-ots-linux -inside=true -get=true -num=100 -cols=5000
./golang-ots-linux -inside=true -get=true -num=100 -cols=10000
./golang-ots-linux -inside=true -get=true -num=100 -cols=30000
./golang-ots-linux -inside=true -get=true -num=100 -cols=50000
exit