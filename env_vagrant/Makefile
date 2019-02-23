setup:
	sudo apt-get update
	sudo apt-get -y upgrade
	cd /tmp
	wget https://dl.google.com/go/go1.11.linux-amd64.tar.gz
	sudo tar -xvf go1.11.linux-amd64.tar.gz
	sudo mv go /usr/local
	echo 'export GOROOT=/usr/local/go
	export GOPATH=$HOME/go
	export PATH=$GOPATH/bin:$GOROOT/bin:$PATH' >> ~/.profile
