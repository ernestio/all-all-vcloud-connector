ruby '2.3.0', engine: 'jruby', engine_version: '9.1.2.0'

source 'https://rubygems.org'

gem 'nokogiri'
gem 'myst', path: '/opt/ernest-libraries/myst/', platform: :jruby

group :development, :test do
  gem 'pry'
end

group :test do
  gem 'rspec'
  gem 'rubocop',   require: false
  gem 'simplecov', require: false
end
