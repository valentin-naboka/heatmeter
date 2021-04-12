
#!/bin/bash
CONFIG_FILE="generated/config.go"

if test -f "$CONFIG_FILE"; then
    echo "$CONFIG_FILE has already been generated. Remove it first, if you want to regenerate."
    exit 0
fi

KEY=$(openssl rand -hex 32)
IV=$(openssl rand -hex 16)

decrypt(){
  echo -n $1 | openssl enc -aes-256-cfb -e -a -A -K $KEY -iv $IV
}

read -p "Enter host: " HOST
DEC_HOST=$(decrypt $HOST)
printf $DEC_HOST

read -p "Enter email: " EMAIL
DEC_EMAIL=$(decrypt $EMAIL)

read -p "Enter password: " PASS
DEC_PASS=$(decrypt $PASS)

read -p "Enter bot api token: " BOT_TOKEN
DEC_BOT_TOKEN=$(decrypt $BOT_TOKEN)

echo "package generated

const Data1 string = \`$KEY\`
const Data2 string = \`$IV\`

const Data3 string = \`$DEC_HOST\`
const Data4 string = \`$DEC_EMAIL\`
const Data5 string = \`$DEC_PASS\`
const Data6 string = \`$DEC_BOT_TOKEN\`" > $CONFIG_FILE
