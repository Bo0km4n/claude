up:
	make -C env
down:
	make -C env down
proxy:
	go build ./pkg/proxy
tablet:
	go build ./pkg/tablet
proto:
	make -C ./pkg/common/proto
