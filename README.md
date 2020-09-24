# go-json
基于simdjson-go的json处理库，针对decode方法做了优化，实现了对象复用。在高并发场景下，能够有效减少内存分配，以降低gc带来的压力。
