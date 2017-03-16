require 'myst'

include Myst::Providers::VCloud

def process_instance(data)
  case data[:_action]
  when 'create'
    create_instance(data)
  when 'update'
    update_instance(data)
  when 'delete'
    delete_instance(data)
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

def create_instance(data)
  provider = get_provider(data)

  instance = ComputeInstance.new(client: provider.client)
  datacenter = provider.datacenter(data[:datacenter_name])
  network = datacenter.private_network(data[:network])
  image = provider.image(data[:reference_image],
                         data[:reference_catalog])

  datacenter.add_compute_instance(instance, data[:name], network, image)


  unless data[:shell_commands].nil? || data[:shell_commands].empty?
    custom_section = instance.vm.getGuestCustomizationSection
    custom_section.setEnabled(true)
    custom_section.setCustomizationScript(data[:shell_commands].join("\n"))
    instance.vm.updateSection(custom_section).waitForTask(0, 1000)
  end

  update_instance(data)

  'instance.create.vcloud.done'
rescue => e
  puts e
  puts e.backtrace
  'instance.create.vcloud.error'
end

def update_instance(data)
  provider = get_provider(data)

  datacenter  = provider.datacenter(data[:datacenter_name])
  instance    = datacenter.compute_instance(data[:name])
  instance.tasks.each { |task| task.waitForTask(0, 1000) }

  instance.power_off if instance.vapp.isDeployed

  instance.hostname = data[:hostname]
  instance.cpus = data[:cpus]
  instance.memory = data[:ram]

  if instance.nics.count == 0
    network_interace = NetworkInterface.new(client:     client,
                                            id:         0,
                                            network:    data[:network],
                                            ipaddress:  data[:ip],
                                            primary:    true)
    instance.add_nic(network_interace)
  end
  instance.tasks.each { |task| task.waitForTask(0, 1000) }

  if data[:disks]
    existing_disks = instance.disks.select(&:isHardDisk)

    data[:disks].each do |disk|
      if existing_disks[disk[:id]].nil?
        dsk = AttachedStorage.new(client: client, id: disk[:id], size: disk[:size])
        instance.add_disk(dsk)
      elsif existing_disks[disk[:id]].getHardDiskSize < disk[:size]
        existing_disks[disk[:id]].updateHardDiskSize(disk[:size])
        instance.vm.updateDisks(instance.vm.getDisks).waitForTask(0, 1000)
      end
    end
  end

  instance.tasks.each { |task| task.waitForTask(0, 1000) }
  instance.power_on

  'instance.update.vcloud.done'
rescue => e
  puts e
  puts e.backtrace
  'instance.update.vcloud.error'
end

def delete_instance(data)
  provider = get_provider(data)

  datacenter  = provider.datacenter(data[:datacenter_name])
  instance    = datacenter.compute_instance(data[:name])
  instance.tasks.each { |task| task.waitForTask(0, 1000) }

  instance.delete unless instance.vapp.nil?

  'instance.delete.vcloud.done'
rescue => e
  puts e
  puts e.backtrace
  'instance.delete.vcloud.error'
end
