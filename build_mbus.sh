
#!/bin/sh

BUILD_DIR=./build
if [ ! -d "$BUILD_DIR" ]; then
  REPO_URL="https://github.com/rscada/libmbus"
  git clone  $REPO_URL $BUILD_DIR
fi

pushd ${BUILD_DIR}
./build.sh
popd

go build .