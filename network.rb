require 'myst'
require 'nokogiri'

include Myst::Providers::VCloud

def process_network(data)
  case data[:_action]
  when 'create'
    create_network(data)
  when 'delete'
    delete_network(data)
  end
end

def create_network(data)
  provider = get_provider(data)

  datacenter      = provider.datacenter(data[:datacenter_name])
  router          = datacenter.router(data[:router_name])

  router.wait_for_tasks

  private_network_request = PrivateNetwork.new.instantiate(router,
                                                           data[:name],
                                                           { start_address: data[:start_address],
                                                             end_address:   data[:end_address] },
                                                           data[:netmask],
                                                           data[:gateway],
                                                           dns(data))
  datacenter.add_private_network(private_network_request)
  'network.create.vcloud.done'
rescue => e
  puts e
  puts e.backtrace
  'network.create.vcloud.error'
end

def delete_network(data)
  provider = get_provider(data)

  datacenter      = provider.datacenter(data[:datacenter_name])
  router          = datacenter.router(data[:router_name])
  network         = datacenter.private_network(data[:name])

  # wait for all prior tasks to complete
  router.wait_for_tasks

  # Filthy hack, because the vcloud sdk doesn't let us delete a network as a non-admin
  return 'network.delete.vcloud.done' if network.network.nil?
  network_href = network.network.getReference.getHref
  url = URI.parse(network_href.gsub(/network/, 'admin/network'))

  task = delete_network_request(url, provider.client.vcloud_token)

  # Wait for the delete to finish
  if task.nil?
    'network.delete.vcloud.error'
  else
    task.wait_for_task
    'network.delete.vcloud.done'
  end
rescue => e
  puts e
  puts e.backtrace
  'network.delete.vcloud.error'
end


def dns(data)
  return ['8.8.8.8', '8.8.4.4'] if data[:dns].nil?
  data[:dns]
end


class NetworkTask
  attr_reader :token, :name, :href, :owner, :status, :progress, :xml, :http

  def initialize(args)
    @token = args[:token]
    @http = args[:http]
    load args[:xml]
  end

  def load(xml)
    @xml = Nokogiri::XML(xml)
    @href = @xml.at_xpath('/xmlns:Task').get_attribute('href')
    @name = @xml.at_xpath('/xmlns:Task').get_attribute('operationName')
    @status = @xml.at_xpath('/xmlns:Task').get_attribute('status')

    if @owner == @xml.at_xpath('/xmlns:Task/xmlns:Owner')
      @owner = @xml.at_xpath('/xmlns:Task/xmlns:Owner').get_attribute('href')
    end

    @progress = if @xml.at_xpath('/xmlns:Task/xmlns:Progress')
                  @xml.at_xpath('/xmlns:Task/xmlns:Progress').content
                end
  end

  def update
    request = Net::HTTP::Get.new(@href)
    request['accept'] = 'application/*+xml;version=5.5'
    request['x-vcloud-authorization'] = token
    res = http.request(request)
    load res.body
  end

  def wait_for_task
    while status == 'running'
      sleep 2
      update
    end
    return false if status == 'error'
    true
  end
end

def delete_network_request(url, token)
  5.times do
    req = Net::HTTP::Delete.new(url.path)
    req['x-vcloud-authorization'] = token
    req['Accept'] = 'application/*+xml;version=5.1'

    http = Net::HTTP.new(url.host, url.port)
    http.use_ssl = true
    res = http.start { |h| h.request(req) }

    return NetworkTask.new(xml: res.body, http: http, token: token) if res.code == '202'

    sleep 10
  end

  nil
end
