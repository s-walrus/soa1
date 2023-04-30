# SOA1: Serializer test

*Домашнее задание для курса по микросервисам*

В этом репозитории приложение, которое тестирует разные форматы сериализации.


### Как устроено

Каждый формат тестируется в отдельном docker-контейнере.
Результаты тестирования доступны через прокси,
который по-запросу сам запрашивает их у тестировщиков через tcp.
Сам прокси размёщен в отдельном docker-контейнере.


### Как запустить

```bash
# clone the repository
git clone https://github.com/s-walrus/soa1 && cd soa1

# compile and run servers
cd docker
docker-compose up -d

# make a request
echo 'get_result yaml' | nc localhost 2000
# > YAML – 106 – 17.985µs – 19.055µ

# supported formats: gob xml json protobuf avro yaml message-pack
echo 'get_result gob' | nc localhost 2000
echo 'get_result xml' | nc localhost 2000

# shut down servers
docker-compose down
```

Тестируемые форматы:
+ `gob` (native format in Go)
+ `xml`
+ `json`
+ `protobuf`
+ `avro` (apache avro)
+ `yaml`
+ `message-pack`
