#!/bin/zsh

echo "Warning"
echo "This script must be run from inside the scripts/codecov_upload directory!"

# Constants
CODECOV_TOKEN0=$(cat token.txt)

# Run coverage tool and testing tool + put reports in this directory
cd ../..
go test -v -cover ./... -coverprofile scripts/codecov_upload/coverage.out -coverpkg ./...
go tool cover -func scripts/codecov_upload/coverage.out -o scripts/codecov_upload/coverage2.out

echo " "
echo "Coverage report generated!"
echo " "

# Upload reports
codecovcli create-commit -t $CODECOV_TOKEN0
codecovcli create-report -t $CODECOV_TOKEN0
codecovcli do-upload -t $CODECOV_TOKEN0