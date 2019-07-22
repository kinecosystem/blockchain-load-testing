grep 'Advancing LCL: \[seq=[0-9]*' $1 | \
while read line; do
	ts=$(echo $line | echo `cut -b16-38`)
	node=$(echo $line | echo `cut -b40-45`)
	ledger=$(echo $line | echo `cut -f4 -d=` | cut -f1 -d,)
	echo "insert into ledger_db_end values ('$node', '${ts}Z', $ledger);"
done
