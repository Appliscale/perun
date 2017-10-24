#!/usr/bin/env bash

SKIP_TESTS=false

for i in "$@"
do
case $i in
    --skip_tests)
    SKIP_TESTS=true
    shift
    ;;
    *)
    ;;
esac
done



if [ ! -f $GOPATH/bin/go-bindata ]; then
    echo "Installing go-bindata library"
    go get -u github.com/jteeuwen/go-bindata/...
    go install github.com/jteeuwen/go-bindata/...
fi

SPECIFICATION_DIRECTORY=$GOPATH/src/github.com/Appliscale/cftool/cfspecification/specification
$GOPATH/bin/go-bindata -pkg cfspecification -o $SPECIFICATION_DIRECTORY/CloudFormationResourceSpecification.go $SPECIFICATION_DIRECTORY/CloudFormationResourceSpecification.json
if [ $? -ne 0 ]
then
    exit 1
fi

go get -t -v ./...
if [ $? -ne 0 ]
then
    exit 1
fi

if [ "$SKIP_TESTS" == false ] ; then
    echo "Running tests..."
    go test github.com/Appliscale/cftool/...
    if [ $? -ne 0 ]
    then
        exit 1
    fi
fi

echo "Installing CFTool..."
go install github.com/Appliscale/cftool
if [ $? -ne 0 ]
then
    exit 1
else
    echo "Installation completed!"
fi