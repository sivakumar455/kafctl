# Step 1: Get all packages
packages=$(go list ./...)

# Step 2: Format the packages for -coverpkg
formatted_packages=$(echo $packages | tr ' ' ',')

# Step 3: Run go test with coverage
go test -coverpkg=$formatted_packages -coverprofile=coverage.out ./...

# Step 4: View the coverage report (optional)
go tool cover -html=coverage.out