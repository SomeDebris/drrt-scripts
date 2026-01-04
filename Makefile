programs=drrt-scheduler drrt-namer drrt-datasheet-updater

all: $(programs)

$(programs): drrt-%: %
	go build -o $@ ./$^

clean:
	rm -f $(programs)
