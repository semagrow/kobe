FROM python:3.6

COPY . /srv/flask_app
WORKDIR /srv/flask_app

RUN apt-get -y update
RUN apt-get install -y --no-install-recommends apt-utils
RUN apt-get -y install nginx \
    python3-dev \
    build-essential

COPY nginx.conf /etc/nginx

RUN pip install --upgrade pip
RUN pip install Flask \
    uWSGI \
    flask-oidc \
    requests

CMD service nginx start && \
    uwsgi --ini uwsgi.ini

