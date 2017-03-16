require 'rubygems'
require 'bundler/setup'
require 'json'

require_relative 'router.rb'
require_relative 'network.rb'
require_relative 'instance.rb'

unless defined? @@test
  @data = { type: ARGV[0] }
  @data.merge! JSON.parse(ARGV[1], symbolize_names: true)

  original_stdout = $stdout
  $stdout = StringIO.new
  begin
      case @data[:_component]
      when "router"
        @data[:type] = process_router(@data)
      when "network"
        @data[:type] = process_network(@data)
      when "instance"
        @data[:type] = process_instance(@data)
      end

    if @data[:type].include? 'error'
      @data['error'] = { code: 0, message: $stdout.string.to_s }
      @data[:_state] = 'errored'
    else
      @data[:_state] = 'completed'
    end
  ensure
    $stdout = original_stdout
  end

  puts @data.to_json
end
