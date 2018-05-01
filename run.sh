docker stop explorer-mongo
docker run --rm -d -p 28000:27017 -v /tmp/.data:/data/db --name=explorer-mongo  mongo:3.6.4
sleep 5
export EXPLORER_MONGODB=mongodb://localhost:28000/explorer
go run main.go