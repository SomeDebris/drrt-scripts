programs=drrt-scheduler drrt-namer

all: $(programs)

$(programs): drrt-%: %
	go build -o $@ ./$^

clean:
	rm -f $(programs)
