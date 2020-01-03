for DEP in `cat dependencies`; do
	echo "getting $DEP"
	go get -u $DEP
done
