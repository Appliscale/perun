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

if [ ! -f /etc/.Appliscale/cftool/config.yaml ]; then
    echo "Extracting configuration file..."
    sudo mkdir -p /etc/.Appliscale/cftool/
    sudo cp config.yaml /etc/.Appliscale/cftool/config.yaml
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