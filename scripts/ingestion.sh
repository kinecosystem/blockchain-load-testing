while read line; do
	time=$(awk '{print $1}' <<<$line | awk -F= '{print $2}') 
	first=$(awk '{print $7}' <<<$line | awk -F= '{print $2}')
	last=$(awk '{print $8}' <<<$line | awk -F= '{print $2}')

	time="'$(cut -b2-25 <<<$time)'"

	for (( ledger=$first; ledger<=$last; ledger++)); do
		echo "insert into ingestion values ($ledger, $time);"
	done
 done
