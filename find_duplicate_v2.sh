#!/bin/bash

# Проверка наличия аргумента с путем к каталогу
if [ $# -eq 0 ]; then
    echo "Использование: $0 /путь/к/каталогу [количество_потоков]"
    exit 1
fi

DIRECTORY="$1"

# Проверка существования каталога
if [ ! -d "$DIRECTORY" ]; then
    echo "Ошибка: Каталог '$DIRECTORY' не существует!"
    exit 1
fi

# Определение количества потоков для параллельного выполнения
# По умолчанию используем количество ядер процессора
if [ $# -ge 2 ] && [[ "$2" =~ ^[0-9]+$ ]]; then
    NUM_THREADS=$2
else
    # Определяем количество процессоров для параллельной обработки на macOS
    if [ -x /usr/sbin/sysctl ]; then
        NUM_THREADS=$(sysctl -n hw.ncpu)
    else
        # По умолчанию, если не удалось определить
        NUM_THREADS=4
    fi
fi

echo "Сканирование каталога '$DIRECTORY' на наличие дубликатов файлов..."
echo "Используется $NUM_THREADS потоков для параллельной обработки"

# Временные файлы для хранения хешей и списка файлов
TEMP_DIR=$(mktemp -d)
TEMP_FILE="$TEMP_DIR/hashes.txt"
FILES_LIST="$TEMP_DIR/files.txt"
COUNTER_FILE="$TEMP_DIR/counter"

# Подсчет общего количества файлов для индикатора прогресса
echo "Подсчет количества файлов..."
find "$DIRECTORY" -type f > "$FILES_LIST"
TOTAL_FILES=$(wc -l < "$FILES_LIST")
echo "Всего файлов: $TOTAL_FILES"

if [ "$TOTAL_FILES" -eq 0 ]; then
    echo "В каталоге не найдено файлов для обработки."
    rm -rf "$TEMP_DIR"
    exit 0
fi

# Функция для обновления счетчика и отображения прогресса
update_progress() {
    local current=$1
    local percentage=$((current * 100 / TOTAL_FILES))
    local completed=$((percentage / 2))
    local remaining=$((50 - completed))
    
    printf "\r[%-${completed}s%-${remaining}s] %d%%  (%d из %d файлов)" \
        "$(printf '%0.s#' $(seq 1 $completed))" \
        "$(printf '%0.s-' $(seq 1 $remaining))" \
        "$percentage" "$current" "$TOTAL_FILES"
}

# Инициализация счетчика
echo 0 > "$COUNTER_FILE"

# Создаем отдельные списки файлов для каждого потока
echo "Распределение файлов между потоками..."
for ((i=1; i<=NUM_THREADS; i++)); do
    awk -v thread=$i -v threads=$NUM_THREADS 'NR % threads == (thread - 1) % threads' "$FILES_LIST" > "$TEMP_DIR/thread_$i.list"
done

# Функция для обработки списка файлов одним потоком
process_thread_list() {
    local thread_id=$1
    local thread_file="$TEMP_DIR/thread_$thread_id.list"
    local result_file="$TEMP_DIR/result_$thread_id"
    local thread_counter=0
    local global_counter=0
    
    while IFS= read -r file; do
        # Вычисляем SHA256 хеш файла
        if [ -f "$file" ]; then  # Проверяем, что файл все еще существует
            hash=$(shasum -a 256 "$file" | awk '{print $1}')
            echo "$hash|$file" >> "$result_file"
        fi
        
        thread_counter=$((thread_counter + 1))
        
        # Обновляем глобальный счетчик безопасным способом
        # Используем атомарную операцию для macOS без flock
        local current_count=0
        {
            # Читаем текущее значение
            read current_count < "$COUNTER_FILE"
            # Увеличиваем и записываем обратно атомарно
            echo $((current_count + 1)) > "$COUNTER_FILE"
        }
        
        # Читаем глобальный счетчик заново
        global_counter=$(<"$COUNTER_FILE")
        
        # Обновляем прогресс каждые несколько файлов
        if [ $((global_counter % 10)) -eq 0 ] || [ "$global_counter" -eq "$TOTAL_FILES" ]; then
            # Обновляем глобальный счетчик
            update_progress "$global_counter"
        fi
    done < "$thread_file"
}

# Запускаем обработку файлов параллельно
echo "Вычисление SHA256 хешей (параллельно в $NUM_THREADS потоках)..."

# Запускаем потоки обработки
for ((i=1; i<=NUM_THREADS; i++)); do
    process_thread_list $i &
done

# Ждем завершения всех процессов
wait

# Финальное обновление прогресса
update_progress "$(<"$COUNTER_FILE")"
echo -e "\n\nЗавершено вычисление хешей"

# Объединяем результаты всех потоков
cat "$TEMP_DIR"/result_* > "$TEMP_FILE" 2>/dev/null || true

# Проверяем, есть ли результаты
if [ ! -s "$TEMP_FILE" ]; then
    echo "Ошибка: не удалось создать хеши файлов или результаты обработки пусты."
    rm -rf "$TEMP_DIR"
    exit 1
fi

echo "Анализ хешей..."
echo "Поиск дубликатов..."
echo "Сортировка и группировка результатов..."

# Создаем временный файл для результатов
DUPLICATE_RESULTS="$TEMP_DIR/duplicates"

echo "Результаты поиска дубликатов:" > "$DUPLICATE_RESULTS"
echo "==============================" >> "$DUPLICATE_RESULTS"

# Находим дубликаты
sort "$TEMP_FILE" | awk -F'|' '
{
    hash = $1;
    file = $2;
    
    files_by_hash[hash] = files_by_hash[hash] ? files_by_hash[hash] "\n  - " file : "  - " file;
    count_by_hash[hash]++;
}
END {
    found = 0;
    for (hash in files_by_hash) {
        if (count_by_hash[hash] > 1) {
            print "Найдены дубликаты с хешем SHA256: " hash;
            print "Количество файлов: " count_by_hash[hash];
            print "Файлы:" files_by_hash[hash] "\n";
            found = 1;
        }
    }
    if (!found) {
        print "Дубликаты файлов не найдены.";
    }
}' >> "$DUPLICATE_RESULTS"

# Выводим результаты
cat "$DUPLICATE_RESULTS"

# Считаем статистику
duplicate_groups=$(grep -c "^Найдены дубликаты с хешем SHA256" "$DUPLICATE_RESULTS")
total_duplicates=0

if [ "$duplicate_groups" -gt 0 ]; then
    # Извлекаем и суммируем количество файлов в каждой группе дубликатов
    total_duplicates=$(grep "^Количество файлов:" "$DUPLICATE_RESULTS" | awk '{sum += $3} END {print sum}')
    
    echo "=============================="
    echo "Статистика:"
    echo "Найдено групп дубликатов: $duplicate_groups"
    echo "Общее количество файлов-дубликатов: $total_duplicates"
fi

# Удаление временных файлов
rm -rf "$TEMP_DIR"

echo "Сканирование завершено."
