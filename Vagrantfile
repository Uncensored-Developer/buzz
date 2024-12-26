Vagrant.configure("2") do |config|
  config.vm.box = "hashicorp/bionic64"

  config.vm.synced_folder ".", "/home/vagrant/app"

  config.vm.provision :shell, path: "vagrant_setup_dep.sh"

  config.vm.provision :shell, inline: <<-SHELL
    sudo su - vagrant -c "cd /home/vagrant/app && make docker-up"
  SHELL

  config.vm.network "forwarded_port", guest: 8007, host: 8004
end
