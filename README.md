Данный проект сделать как тестовое задание в рамках отборочного этапа на стажировку в AvitoTech.

Иерархия соответствует заданию, docker-compose находится в корневой папке, модули лежат в соотвтетствующих директориях. Каждый компонент содержит Dockerfile. Трафик к приложению проходит через nginx по порту 8089. Состояние базы Redis сохраняется вне контейнера, что позволяет сохранять данные после перезагрузки контейнера.

Требование "Конфигурация компонентов  «не запекается» внутри образов" было интерпретировано как одно из антипаттернов Docker "Hardcoding secrets and configuration into container images". Все конфигурационные файлы, секреты (пароли и логины для редиса), ssl сертификаты передаются в момент сборки через volumes или .env.


Реализованны необходимые endpoint'ы:
- /set_key - принимает в качестве значения ключ и значение в  формате json:{"\<key\>":"\<value\>"}. В случае успеха, возвращает "Accepted"
- /get_key?key=\<key\> - принимает в качестве параметра значение ключа и возвращает значение в виде html файла. В случае ошибки ключа возвращает 404, в случае неправильного написанного адреса (пример: /get_key/) возвращает 403
- /del_key - принимает в качестве значения ключ формата json:{"key":"\<value\>"}, и удаляет его из Redis'a. Как подтверждение, отправляет на /get_key с параметром, равным ключу удаленной записи.
- Во всех остальных роутах вернется 403 ошибка
  Все контейнеры в качестве базового образа используют Debian. Развертывание выполняется по средствам Docker Compose. Приложение реализовано на Golang. В Redis осуществляется работа только со строками.

В Redis включена поддержка аутентификации по имени пользователя и паролю. Осуществлен TLS обмен данными между клиентом и сервером, путем взаимодействия с заранее сгенерированными сертификатами.

Развертка объекта:
- git clone https://github.com/Azmekin/go_backend_DevOps - скопируйте проект
- cd go_backend_DevOps
- docker-compose up -d - запустите docker-compose

