to build & run:
OPTION 1:
```
./run.sh
```

OPTION 2:
```
cd contact_app
sudo service apache2 stop
sudo service nginx stop
docker-compose down -v --rmi all --remove-orphans
docker-compose build
docker-compose run app rake db:create RAILS_ENV=production
docker-compose run app rake db:migrate db:seed RAILS_ENV=production
docker-compose up -d web   
```
Runs on localhost:3000

From https://itnext.io/docker-rails-puma-nginx-postgres-999cd8866b18