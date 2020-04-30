#!/bin/zsh
# This script sets up 8 vcli environments. In the end there will be directories ~/.vcli1 to ~/
# .vcli8 each with an account called 'voter' which uses the password '12345678'.
# The script prints commands for transferring 1 foo token to the voter accounts. They require
# that one acli application has been setup with an 'admin' account using the password '12345678'.

if [ -e ~/.vcli1 ]; then
  rm -r ~/.vcli1
fi
if [ -e ~/.vcli2 ]; then
  rm -r ~/.vcli2
fi
if [ -e ~/.vcli3 ]; then
  rm -r ~/.vcli3
fi
if [ -e ~/.vcli4 ]; then
  rm -r ~/.vcli4
fi
if [ -e ~/.vcli5 ]; then
  rm -r ~/.vcli5
fi
if [ -e ~/.vcli6 ]; then
  rm -r ~/.vcli6
fi
if [ -e ~/.vcli7 ]; then
  rm -r ~/.vcli7
fi
if [ -e ~/.vcli8 ]; then
  rm -r ~/.vcli8
fi

printf '12345678\n12345678\n' | vcli keys add voter --home ~/.vcli1
printf '12345678\n12345678\n' | vcli keys add voter --home ~/.vcli2
printf '12345678\n12345678\n' | vcli keys add voter --home ~/.vcli3
printf '12345678\n12345678\n' | vcli keys add voter --home ~/.vcli4
printf '12345678\n12345678\n' | vcli keys add voter --home ~/.vcli5
printf '12345678\n12345678\n' | vcli keys add voter --home ~/.vcli6
printf '12345678\n12345678\n' | vcli keys add voter --home ~/.vcli7
printf '12345678\n12345678\n' | vcli keys add voter --home ~/.vcli8

vcli config chain-id pbb --home ~/.vcli1
vcli config output json --home ~/.vcli1
vcli config indent true --home ~/.vcli1
vcli config trust-node true --home ~/.vcli1

vcli config chain-id pbb --home ~/.vcli2
vcli config output json --home ~/.vcli2
vcli config indent true --home ~/.vcli2
vcli config trust-node true --home ~/.vcli2

vcli config chain-id pbb --home ~/.vcli3
vcli config output json --home ~/.vcli3
vcli config indent true --home ~/.vcli3
vcli config trust-node true --home ~/.vcli3

vcli config chain-id pbb --home ~/.vcli4
vcli config output json --home ~/.vcli4
vcli config indent true --home ~/.vcli4
vcli config trust-node true --home ~/.vcli4

vcli config chain-id pbb --home ~/.vcli5
vcli config output json --home ~/.vcli5
vcli config indent true --home ~/.vcli5
vcli config trust-node true --home ~/.vcli5

vcli config chain-id pbb --home ~/.vcli6
vcli config output json --home ~/.vcli6
vcli config indent true --home ~/.vcli6
vcli config trust-node true --home ~/.vcli6

vcli config chain-id pbb --home ~/.vcli7
vcli config output json --home ~/.vcli7
vcli config indent true --home ~/.vcli7
vcli config trust-node true --home ~/.vcli7

vcli config chain-id pbb --home ~/.vcli8
vcli config output json --home ~/.vcli8
vcli config indent true --home ~/.vcli8
vcli config trust-node true --home ~/.vcli8

## Create voter account by sending assets to it ##
s1='echo "12345678" | acli tx bank send $(acli keys show admin -a) '
s2=`vcli keys show voter -a --home ~/.vcli1` 
s3=' 1foo -y'
echo $s1$s2$s3

s1='echo "12345678" | acli tx bank send $(acli keys show admin -a) '
s2=`vcli keys show voter -a --home ~/.vcli2` 
s3=' 1foo -y'
echo $s1$s2$s3

s1='echo "12345678" | acli tx bank send $(acli keys show admin -a) '
s2=`vcli keys show voter -a --home ~/.vcli3` 
s3=' 1foo -y'
echo $s1$s2$s3

s1='echo "12345678" | acli tx bank send $(acli keys show admin -a) '
s2=`vcli keys show voter -a --home ~/.vcli4` 
s3=' 1foo -y'
echo $s1$s2$s3

s1='echo "12345678" | acli tx bank send $(acli keys show admin -a) '
s2=`vcli keys show voter -a --home ~/.vcli5` 
s3=' 1foo -y'
echo $s1$s2$s3

s1='echo "12345678" | acli tx bank send $(acli keys show admin -a) '
s2=`vcli keys show voter -a --home ~/.vcli6` 
s3=' 1foo -y'
echo $s1$s2$s3

s1='echo "12345678" | acli tx bank send $(acli keys show admin -a) '
s2=`vcli keys show voter -a --home ~/.vcli7` 
s3=' 1foo -y'
echo $s1$s2$s3

s1='echo "12345678" | acli tx bank send $(acli keys show admin -a) '
s2=`vcli keys show voter -a --home ~/.vcli8` 
s3=' 1foo -y'
echo $s1$s2$s3
