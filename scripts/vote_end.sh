grep 'setConfirmCommit i:[0-9]*' $1 | \
while read line; do
	ts=$(echo $line | echo `cut -b16-38`)
	node=$(echo $line | echo `cut -b40-45`)
	ledger=$(echo $line | echo `awk -F': ' '{print $2}' | cut -f1 -d" "`)
	echo "insert into ledger_vote_end values ('$node', '${ts}Z', $ledger);"
done
