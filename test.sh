export TOKEN=`cat .token`;
export GUILD=`cat .guild`;
go run bot.go -token="$TOKEN" -guild="$GUILD"