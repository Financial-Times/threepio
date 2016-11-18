GOARGS=
GOENV=CGO_ENABLED=0
EXECUTABLE=threepio

all: main

main:
	go build $(GOARGS)

clean:
	-rm -v $(EXECUTABLE)

dist: clean main 
