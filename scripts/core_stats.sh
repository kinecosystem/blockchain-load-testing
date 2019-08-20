zgrep 'Advancing LCL: \[seq=[0-9]*' $1 | \
while read line; do
	ts=$(echo $line | echo `cut -b16-38`)
	node=$(echo $line | echo `cut -b40-45`)
	ledger=$(echo $line | echo `cut -f4 -d=` | cut -f1 -d,)
	echo "insert into ledger_db_end values ('$node', '${ts}Z', $ledger);"
done

zgrep 'starting closeLedger() on ledgerSeq=[0-9]*' $1 | \
while read line; do
	ts=$(echo $line | echo `cut -b16-38`)
	node=$(echo $line | echo `cut -b40-45`)
	ledger=$(echo $line | echo `cut -f2 -d=`)
	echo "insert into ledger_db_start values ('$node', '${ts}Z', $ledger);"
done

zgrep 'setConfirmCommit i:[0-9]*' $1 | \
while read line; do
	ts=$(echo $line | echo `cut -b16-38`)
	node=$(echo $line | echo `cut -b40-45`)
	ledger=$(echo $line | echo `awk -F': ' '{print $2}' | cut -f1 -d" "`)
	echo "insert into ledger_vote_end values ('$node', '${ts}Z', $ledger);"
done

zgrep 'triggerNextLedger.* slot:[0-9]*' $1 | \
while read line; do
	ts=$(echo $line | echo `cut -b16-38`)
	node=$(echo $line | echo `cut -b40-45`)
	ledger=$(echo $line | echo `awk -F': ' '{print $5}'`)
	echo "insert into ledger_vote_start values ('$node', '${ts}Z', $ledger);"
done

zgrep 'Got consensus:' $1 | \
while read line; do
	node=$(echo $line | echo `cut -b40-45`)
	ledger=$(echo $line | echo `awk -F'=' '{print $2}' | cut -f1 -d","`)
	count=$(echo $line | echo `awk -F'=' '{print $4}' | cut -f1 -d","`)
	echo "insert into ledger_tx_count values ('$node', $count, $ledger);"
done
