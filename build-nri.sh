go build -o 00-nri-core-binding nri_core_binding.go
sudo mkdir -p /opt/nri/plugins/
sudo mv 00-nri-core-binding /opt/nri/plugins/