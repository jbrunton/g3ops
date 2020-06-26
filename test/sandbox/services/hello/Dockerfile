FROM ruby:2.7-alpine

WORKDIR "/app"

COPY ./Gemfile* ./
RUN gem install bundler -v 2.1.4
RUN bundle install
COPY . .

ENTRYPOINT [ "./entrypoint.sh" ]
