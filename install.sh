go build -o db-backup-cli ./cmd/main.go
sudo mv db-backup-cli /usr/local/bin/

which db-backup-cli
db-backup-cli --help
