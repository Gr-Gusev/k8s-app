# ДЗ №2 по курсу "Облачные вычисления и виртуализация информационных ресурсов"

## Описание

Web-приложение представляет собой сервис заметок, с возможностью их добавления и удаления.

## Клонирование репозитория и запуск приложения через minikube

```bash
# клонирование репозитория:
git clone https://github.com/Gr-Gusev/k8s-app.git
cd k8s-app

# запуск minikube:
minikube start

# сборка образа приложения:
docker build -t app .

# сохранение образа в архив:
docker image save -o app_image.tar app

# загрузка образа в minikube:
minikube image load app_image.tar

# запуск приложения:
minikube kubectl -- apply -f manifest.yaml
```

Приложение будет доступно после того, как оба пода (app и mysql) будут в статусе Running:

```bash
minikube kubectl -- get pods
```

Получение адреса, по которому доступно приложение, осуществляется через команду:

```bash
minikube service app-service --url
```

## Проверка PersistentVolume

Для проверки работы постоянного хранилища mysql, можно создать через web-интерфейс заметку,
затем удалить под mysql командой:

```bash
minikube kubectl -- delete pod mysql
```

и заново его создать командой:

```bash
minikube kubectl -- apply -f manifest.yaml
```
