



# How to expose your bosh-lites to the internet ?

- Go to VPC network > Firewall rules and look for something like 'apollo-external'
- Make sure the Source IP ranges is '0.0.0.0/0'
- Edit the rules
- Add a target tag like 'bosh-lite'
- Add 'icmp' to allowed protocols
- Save
