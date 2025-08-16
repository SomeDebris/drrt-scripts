programs=drrt-scheduler

all: $(programs)

$(programs): drrt-%: %
	go build -o $@ ./$^

clean:
	rm -f $(programs)
