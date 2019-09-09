/|/{gsub(/ /, "", $0); print "insert into tx_ledger values ('"$1"', "$2", "$3");"}
