#!/usr/bin/env bash

SKIP_ANALYZE=false
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
case $i in
    --skip_analyze)
    SKIP_ANALYZE=true
    shift
    ;;
    *)
    ;;
esac
done

echo "Extracting configuration file..."
if [ ! -f "$HOME/.config/cftool" ]; then
    mkdir -p "$HOME/.config/cftool"
fi
cp config.yaml "$HOME/.config/cftool/config.yaml"

go get -t -v ./...
if [ $? -ne 0 ]
then
    exit 1
fi

if [ "$SKIP_TESTS" == false ] ; then
    echo "Running tests..."
    go test github.com/Appliscale/cftool/... --v -cover
    if [ $? -ne 0 ]
    then
        exit 1
    fi
fi

if [ "$SKIP_ANALYZE" == false ] ; then
    echo "Analyzing code..."
    go tool vet -v ./
fi

echo "Installing CFTool..."
go install github.com/Appliscale/cftool
if [ $? -ne 0 ]
then
    exit 1
else
    echo "Installation completed!"
fi
