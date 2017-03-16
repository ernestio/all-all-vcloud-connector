require 'myst'

include Myst::Providers::VCloud

def process_router(data)
  case data[:_action]
  when 'create', 'update'
    configure_router(data)
  end
end

def get_provider(data)
  usr = ENV['DT_USR'] || data[:datacenter_username]
  pwd = ENV['DT_PWD'] || data[:datacenter_password]
  credentials = usr.split('@')

  Provider.new(endpoint:     data[:vcloud_url],
               organisation: credentials.last,
               username:     credentials.first,
               password:     pwd)
end

def configure_router(data)
  provider = get_provider(data)

  datacenter      = provider.datacenter(data[:datacenter_name])
  router          = datacenter.router(data[:name])

  firewall = router.firewall
  firewall.purge_rules

  data[:firewall_rules].each do |rule|
    firewall.add_rule({ ip: rule[:source_ip], port_range: rule[:source_port] },
                      { ip: rule[:destination_ip], port_range: rule[:destination_port] }, rule[:protocol].to_sym)
  end

  router.update_service(firewall)

  nat = router.nat
  nat.purge_rules

  data[:nat_rules].each do |rule|
    interface_ref = router.interface_network_reference(rule[:network])
    nat.add_rule(rule[:type],
                 { ip: rule[:origin_ip], port: rule[:origin_port] },
                 { ip: rule[:translation_ip], port: rule[:translation_port] },
                 rule[:protocol],
                 interface_ref)
  end

  router.update_service(nat)

  data[:_component] + '.' + data[:_action] + '.' + data[:_provider] + '.done'
rescue => e
  puts e
  puts e.backtrace
  data[:_component] + '.' + data[:_action] + '.' + data[:_provider] + '.error'
end
