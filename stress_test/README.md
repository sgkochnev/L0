Запуск стресс теста на 2 мин. по 100 запросов в секунду
```bash
vegeta attack -duration=120s -rate=100 -targets=./stress_test/target.list -output=./stress_test/attack.bin
```

Формирование графика
```bash
vegeta plot -title=Attack ./stress_test/attack.bin > ./stress_test/results.html    
```

Формирование отчета в терминале
```bash
vegeta report ./stress_test/attack.bin 
```

Формирование отчета в формате json
```bash
vegeta report -type=json ./stress_test/attack.bin 
```