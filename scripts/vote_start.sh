grep 'triggerNextLedger.* slot:[0-9]*' $1 | \
while read line; do
	ts=$(echo $line | echo `cut -b16-38`)
	node=$(echo $line | echo `cut -b40-45`)
	ledger=$(echo $line | echo `awk -F': ' '{print $5}'`)
	echo "insert into ledger_vote_start values ('$node', '${ts}Z', $ledger);"
done
