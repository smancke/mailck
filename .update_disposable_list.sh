#!/bin/bash

file=disposable_list.go

cat > $file <<EOF
package mailck

var DisposableDomains = []string{
EOF

curl -s  https://rawcdn.githack.com/andreis/disposable/master/domains.txt \
     | sed  's/\(.*\)/\t"\1",/' >> $file


echo >> $file
echo "}" >> $file
