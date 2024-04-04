run:
	@rm -rf .sasy;\
	go build cmd/sasy/main.go;\
	./main init;\
	./main commit;\
	rm main;

del:
	@rm -rf .sasy;
init:
	@rm -rf .sasy;\
	go run cmd/sasy/main.go init;

commit: 
	@go run cmd/sasy/main.go commit;

rm:
	@rm main;

build:
	@rm -rf .sast;\
	go build cmd/sasy/main.go;\

