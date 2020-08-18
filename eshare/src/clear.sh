#/bin/sh
rm -rf static/files/*.svg
mongo 127.0.0.1:27017 <<EOF
use eshare
db.document.drop()
db.page.drop()
EOF
