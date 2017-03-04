#!/bin/bash

file=disposable_list.go

cat > $file <<EOF
package mailck

// DisposableDomains is a list of fake mail providers.
// The list was taken from https://github.com/andreis/disposable
// License: MIT 
// Last updated: `date`
var DisposableDomains = map[string]bool{
EOF

curl -s  https://rawcdn.githack.com/andreis/disposable/master/domains.txt \
     | sed  's/\(.*\)/\t"\1": true,/' >> $file


echo >> $file
echo "}" >> $file
