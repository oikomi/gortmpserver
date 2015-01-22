#!/bin/sh

echo "# Hack runtime"

cat >> $GOROOT/src/runtime/runtime.c << EOF

void
runtime·GetGoId(int32 ret)
{
	ret = g->goid;
	USED(&ret);
}

EOF

cat >> $GOROOT/src/runtime/extern.go << EOF

func GetGoId() int32

EOF

cd $GOROOT/src
./make.bash

cat > $$$$.go << EOF
package main

import (
	"fmt"
	"runtime"
)

func main() {
	runtime.GetGoId()
	fmt.Print("done")
}
EOF
x=`go run $$$$.go`
rm $$$$.go

echo ""
echo ""
echo "---"

if [ $x = "done" ]; then
	echo "Done"
else
	echo "Failed"
fi
