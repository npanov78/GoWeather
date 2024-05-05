# GoWeather

Консольная утилита, показывающая прогноз погоды в удобном табличном виде. 
Основана на работе с API http://api.weatherapi.com/v1/

## Запуск проекта 
1. Установить переменную окружения ```export API_KEY=<KEY>```
2. Собрать бинарь ```go build main.go -O forecast```
3. Запустить проект ```./forecast -d <days> -c <city>```

## Установка утилиты в систему Linux
1. Собрать бинать 
2. Копировтаь бинарь в директорию $PATH: ```cp forecast /usr/bin/forecast```
3. Получить API_KEY и задать его в ```.bashrc```: ```echo "export API_KEY=<your_key> >> ~/.bashrc"```
4. Задать alias для своего города и количества дней: ```echo "alias forecast='forecast -d 7 -c Saint-Petersburg'" >> ~/.bashrc```
