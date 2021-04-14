
#!/bin/bash

if [ $# -gt 1 ]; then
    echo -e " Usage: $0 [OPTION]\n\n" \
         "Options: \n" \
          "--arm    build for Raspberry Pi"
    exit 1
fi

BUILD_DIR=./build
if [ ! -d "$BUILD_DIR" ]; then
  REPO_URL="https://github.com/rscada/libmbus"
  git clone  $REPO_URL $BUILD_DIR
fi

pushd ${BUILD_DIR}

if [ $(uname) == "Darwin" ]; then
  sed -i '' 's/^.*\&\& \.\/configure$/& --enable-shared=no/' build.sh
  
  if [ $# -eq 1 ] && [ $1 == "--arm" ]; then
    echo "Warning: --arm option has no effect on $(uname)"
  fi

else
  sed -i 's/^.*\&\& \.\/configure$/& --enable-shared=no/' build.sh
  
  if [ $# -eq 1 ] && [ $1 == "--arm" ]; then
    sed -i 's/^.*\&\& \.\/configure --enable-shared=no$/& --build=x86_64-ubuntu-linux --host=arm-linux-gnueabihf /' build.sh
  fi  
fi

./build.sh
popd

if [ $# -eq 1 ] && [ $1 == "--arm" ] && [ $(uname) != "Darwin" ]; then
  env CC="arm-linux-gnueabihf-gcc" LD="arm-linux-gnueabihf-ld"  GOOS=linux GOARCH=arm GOARM=5 CGO_ENABLED=1 go build -v  .
else
  go build -v  .
fi
