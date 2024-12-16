export SEALIGHTS_TOKEN="<redacted>"
export OS_ARCH="linux-amd64"
export SERVICE_NAME="go-calc-demo"
export BRANCH_NAME="main"
export BUILD_NAME="konflux-$(head /dev/urandom | tr -dc 'a-z0-9' | head -c 6)"

wget -qO- https://agents.sealights.co/slgoagent/latest/slgoagent-${OS_ARCH}.tar.gz | tar -xzv
wget -qO- https://agents.sealights.co/slcli/latest/slcli-${OS_ARCH}.tar.gz | tar -xzv

./slcli config init --lang go --token "${SEALIGHTS_TOKEN}"
./slcli config create-bsid --app ${SERVICE_NAME} --branch ${BRANCH_NAME} --build ${BUILD_NAME}

# Run Sealights scan
./slcli scan --bsid buildSessionId.txt --path-to-scanner ./slgoagent \
    --workspacepath ./ --scm git --scmProvider github

podman build -t quay.io/flacatus/konflux-sealights:0.2 -f Dockerfile

podman rm -f sealights-test || true

podman run --name sealights-test -d -it -p 8080:8080 --rm quay.io/flacatus/konflux-sealights:0.2

export BUILD_SESSION_ID=$(cat buildSessionId.txt)
export CONTAINER_ROUTE_URL="http://127.0.0.1:8080"
