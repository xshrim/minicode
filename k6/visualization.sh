# build
git clone 'https://github.com/loadimpact/k6'
cd k6
git submodule update --init
docker-compose up -d influxdb grafana
docker-compose run -v $PWD/samples:/scripts k6 run /scripts/es6sample.js

# run
# docker-compose up

# url
# http://localhost:3000

# dashboard
# import ID 2587, by Dave Cadwallader -
# https://grafana.com/dashboards/2587

# document
# https://docs.k6.io/docs/welcome
