#!/bin/zsh
# This script produces a string that can be used to transfer 1 foo token to the account with the
# name 'voter'. THe sender is the account with the name 'admin'. The admin account is an account
# of the acli application and the voter is an account of the vcli application.
# The password for the acli account needs to be '12345678'.

s1='echo "12345678" | acli tx bank send $(acli keys show admin -a) '
s2=$(vcli keys show voter -a --home ~/.vcli)
s3=' 1foo -y'
echo $s1$s2$s3

