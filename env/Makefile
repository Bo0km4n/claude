setup:
	sudo apt-get update
	sudo apt-get -y upgrade
	cd /tmp
	wget https://dl.google.com/go/go1.11.linux-amd64.tar.gz -O ~/tmp/go1.11.linux-amd64.tar.gz
	sudo tar -xvf ~/tmp/go1.11.linux-amd64.tar.gz -C ~/tmp go
	sudo mv ~/tmp/go /usr/local
	echo "export GOROOT=/usr/local/go" >> ~/.bashrc
	echo "export GOPATH=$HOME/go" >> ~/.bashrc
	echo "export PATH=$GOPATH/bin:$GOROOT/bin:$PATH" >> ~/.bashrc
