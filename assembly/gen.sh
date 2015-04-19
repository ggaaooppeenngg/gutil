FILE=main.go
PACKAGE=main

if [ -n "$2" ];then
	FILE=$2
fi

if [ -n "$1" ];then
	PACKAGE=$1
fi

if [ ! -f "$FILE" ];then
	touch $FILE
fi

cat >$FILE <<EOF

package $PACKAGE

func foo() {

}
EOF
echo "initial file $FILE is generated!"
