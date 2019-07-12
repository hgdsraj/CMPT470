cd contact_app
docker-compose down -v --rmi all --remove-orphans
docker-compose build
docker-compose run app rake db:create RAILS_ENV=production
docker-compose run app rake db:migrate db:seed RAILS_ENV=production
docker-compose up web
